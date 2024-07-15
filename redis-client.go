package redisclient

import (
	"bufio"
	"errors"
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"
)

type RedisClient struct {
	conn   net.Conn
	reader bufio.Reader
	Opts
}

type RedisOptions struct {
	Opts
}

type Opts struct {
	host                   string
	port                   int
	db                     int
	username               string
	password               string
	ssl                    bool
	socketTimeout          time.Duration
	socketConnectTimeout   time.Duration
	socketKeepAlive        bool
	socketKeepaliveOptions *KeepAliveOptions
	connectionPool         *ConnectionPool
	unixSocketPath         string
	encoding               string
	encodingErrors         string
	charset                string
	errors                 []error
	decodeResponses        bool
	retryOnTimeout         bool
	retryOnError           []error
	sslKeyFile             string
	sslCertFile            string
	sslCertReqs            string
	sslCaCerts             string
	sslCaPath              string
	sslCaData              []byte
	sslCheckhostname       bool
	sslPassword            string
	sslValidateOcsp        bool
	sslValidateOcspstapled bool
	sslOcspContext         interface{} // this might need a more specific type
	sslocspexpectedcert    []byte
	maxconnections         int
	singleConnectionClient bool
	healthcheckInterval    time.Duration
	clientname             string
	libname                string
	libversion             string
	retry                  *RetryOptions
	redisConnectFunc       RedisConnectFunc
	credentialProvider     CredentialProvider
	protocol               int
}

// TODO: Fill in
type KeepAliveOptions struct{}
type RetryOptions struct{}
type CredentialProvider struct{}
type ConnectionPool struct{}
type RedisConnectFunc func() (net.Conn, error)
type SetOpts struct {
	NX      bool
	XX      bool
	Get     bool
	EX      int64 //seconds item exists for
	PX      int64 //milliseconds item exists for
	EXAT    *time.Time
	PXAT    *time.Time
	KeepTTL bool
}

type OptsFunc func(*Opts)

type SetOptsFunc func(*SetOpts)

func NewClient(host string, port int, opts ...OptsFunc) (*RedisClient, error) {

	defaultOpts := defaultOptions()

	for _, fn := range opts {
		fn(&defaultOpts)
	}

	var conn net.Conn

	var err error

	if defaultOpts.redisConnectFunc != nil {
		conn, err = defaultOpts.redisConnectFunc()
	} else {
		conn, err = net.Dial("tcp", fmt.Sprintf("%s:%d", host, port))
	}

	if err != nil {
		return nil, fmt.Errorf("failed to connect to redis %w", err)
	}

	client := &RedisClient{
		conn:   conn,
		reader: *bufio.NewReader(conn),
		Opts:   defaultOpts,
	}

	auth_err := client.authenticate()
	if auth_err != nil {
		client.Close()
		return nil, fmt.Errorf("error while authenticating %w", auth_err)
	}

	return client, nil
}

func (o *SetOpts) WithNX() *SetOpts {
	o.NX = true
	o.XX = false // NX and XX are mutually exclusive
	return o
}

func (o *SetOpts) WithXX() *SetOpts {
	o.XX = true
	o.NX = false // NX and XX are mutually exclusive
	return o
}

func (o *SetOpts) WithGet() *SetOpts {
	o.Get = true
	return o
}

func (o *SetOpts) WithEX(seconds int64) *SetOpts {
	o.EX = seconds
	o.PX = 0
	o.EXAT = nil
	o.PXAT = nil
	o.KeepTTL = false
	return o
}

func (o *SetOpts) WithPX(milliseconds int64) *SetOpts {
	o.PX = milliseconds
	o.EX = 0
	o.EXAT = nil
	o.PXAT = nil
	o.KeepTTL = false
	return o
}

func (o *SetOpts) WithEXAT(timestamp time.Time) *SetOpts {
	o.EXAT = &timestamp
	o.EX = 0
	o.PX = 0
	o.PXAT = nil
	o.KeepTTL = false
	return o
}

func (o *SetOpts) WithPXAT(timestamp time.Time) *SetOpts {
	o.PXAT = &timestamp
	o.EX = 0
	o.PX = 0
	o.EXAT = nil
	o.KeepTTL = false
	return o
}

func (o *SetOpts) WithKeepTTL() *SetOpts {
	o.KeepTTL = true
	o.EX = 0
	o.PX = 0
	o.EXAT = nil
	o.PXAT = nil
	return o
}

func NewSetOpts() *SetOpts {
	return &SetOpts{}
}

func defaultOptions() Opts {
	return Opts{
		host:            "localhost",
		port:            6379,
		libname:         "redis-client-go",
		libversion:      "0.0.1",
		protocol:        2,
		decodeResponses: false,
		db:              0,
	}
}

func WithCustomConnectFunc(connFunc RedisConnectFunc) OptsFunc {
	return func(o *Opts) {
		o.redisConnectFunc = connFunc
	}
}

func DecodeResponses(shouldDecode bool) OptsFunc {
	return func(o *Opts) {
		o.decodeResponses = shouldDecode
	}
}

func WithClientName(clienName string) OptsFunc {
	return func(o *Opts) {
		o.clientname = clienName
	}
}

func WithUsername(username string) OptsFunc {
	return func(o *Opts) {
		o.username = username
	}
}

func WithPassword(password string) OptsFunc {
	return func(o *Opts) {
		o.password = password
	}
}

func encodeCommand(args []string) []byte {
	buf := []byte{'*'}
	buf = strconv.AppendInt(buf, int64(len(args)), 10)
	buf = append(buf, '\r', '\n')
	for _, arg := range args {
		buf = append(buf, '$')
		buf = strconv.AppendInt(buf, int64(len(arg)), 10)
		buf = append(buf, '\r', '\n')
		buf = append(buf, arg...)
		buf = append(buf, '\r', '\n')
	}

	return buf
}

func (r *RedisClient) Close() error {
	return r.conn.Close()
}

func (r *RedisClient) sendCommand(args []string) error {
	encoded_command := encodeCommand(args)
	_, err := r.conn.Write(encoded_command)
	if err != nil {
		return err
	}
	return err
}

func (r *RedisClient) Do(args ...string) (interface{}, error) {
	if err := r.sendCommand(args); err != nil {
		return nil, fmt.Errorf("error while sending command %w", err)
	}
	return r.readResponse()
}

func (r *RedisClient) Get(key string) (interface{}, error) {
	resp, err := r.Do("GET", key)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (r *RedisClient) SetWithOptions(key string, val string, opts *SetOpts) (interface{}, error) {
	args := []string{"SET", key, val}

	if opts.NX {
		args = append(args, "NX")
	} else if opts.XX {
		args = append(args, "XX")
	}

	if opts.Get {
		args = append(args, "GET")
	}

	if opts.EX != 0 {
		args = append(args, "EX", strconv.Itoa(int(opts.EX))) // set amount of seconds
	} else if opts.PX != 0 {
		args = append(args, "PX", strconv.FormatInt(int64(opts.PX), 10)) // set amount of milliseconds
	} else if opts.EXAT != nil {
		args = append(args, "EXAT", strconv.FormatInt(opts.EXAT.Unix(), 10)) //UNIX timestamp with seconds
	} else if opts.PXAT != nil {
		args = append(args, "PXAT", strconv.FormatInt(opts.PXAT.UnixNano()/int64(time.Millisecond), 10)) //UNIX timestamp with milliseconds
	} else if opts.KeepTTL {
		args = append(args, "KEEPTTL")
	}

	resp, err := r.Do(args...)

	if err != nil {
		return nil, err
	}

	// Handle the GET option
	if opts.Get {
		if resp == nil {
			return nil, nil // Key didn't exist before
		}
		str, ok := resp.(string)
		if !ok {
			return nil, fmt.Errorf("unexpected response type for SET with GET: %T", resp)
		}
		return str, nil
	}

	// Handle normal SET response
	if resp == nil {
		return nil, nil // SET NX/XX condition not met
	}
	_, ok := resp.(string)
	if !ok {
		return nil, fmt.Errorf("unexpected response type for SET: %T", resp)
	}
	return resp, nil
}

func (r *RedisClient) Exists(args ...string) (interface{}, error) {
	command_args := make([]string, 0, len(args)+1)
	command_args = append(command_args, "EXISTS")
	command_args = append(command_args, args...)
	resp, err := r.Do(command_args...)
	if err != nil {
		return nil, fmt.Errorf("unknown command %w", err)
	}
	return resp, nil
}

func (r *RedisClient) Echo(val string) (interface{}, error) {
	resp, err := r.Do("ECHO", val)
	if err != nil {
		return nil, fmt.Errorf("error while sending echo command %w", err)
	}
	return resp, nil
}

// Lpop is unfinished
func (r *RedisClient) LPop(listname string) (interface{}, error) {
	resp, err := r.Do("LPOP", listname)
	if err != nil {
		return nil, fmt.Errorf("erorr while sending lpop command: %w", err)
	}
	return resp, nil
}
func (r *RedisClient) Set(key string, val string) (interface{}, error) {
	resp, err := r.Do("SET", val)
	if err != nil {
		return nil, fmt.Errorf("error while sending set command %w", err)
	}
	return resp, nil
}

func (r *RedisClient) SetNx(key string, val string) (interface{}, error) {
	resp, err := r.Do("SETNX", key, val)
	if err != nil {
		return nil, fmt.Errorf("erorr while sending setnx command: %w", err)
	}
	return resp, nil
}

func (r *RedisClient) SetRange(key string, offset int, val string) (interface{}, error) {
	resp, err := r.Do("SETRANGE", key, strconv.Itoa(offset), val)
	if err != nil {
		return nil, fmt.Errorf("erorr while sending setrange command: %w", err)
	}
	return resp, nil
}

func (r *RedisClient) StrLen(key string) (interface{}, error) {
	resp, err := r.Do("STRLEN", key)
	if err != nil {
		return nil, fmt.Errorf("erorr while sending strlen command: %w", err)
	}
	return resp, nil
}

func (r *RedisClient) Sadd(args ...string) (interface{}, error) {
	command_args := []string{"SADD"}
	command_args = append(command_args, args...)
	resp, err := r.Do(command_args...)
	if err != nil {
		return nil, fmt.Errorf("erorr while sending sadd command: %w", err)
	}
	return resp, nil
}

func (r *RedisClient) Scard(key string) (interface{}, error) {
	resp, err := r.Do("SCARD", key)
	if err != nil {
		return nil, fmt.Errorf("erorr while sending scard command: %w", err)
	}
	return resp, nil
}

func (r *RedisClient) Sdiff(keys ...string) (interface{}, error) {
	commands_args := []string{"SDIFF"}
	commands_args = append(commands_args, keys...)
	resp, err := r.Do(commands_args...)
	if err != nil {
		return nil, fmt.Errorf("erorr while sending sdiff command: %w", err)
	}
	return resp, nil
}

func (r *RedisClient) SdiffStore(destination string, keys ...string) (interface{}, error) {
	commands_args := []string{"SDIFFSTORE"}
	commands_args = append(commands_args, keys...)
	resp, err := r.Do(commands_args...)
	if err != nil {
		return nil, fmt.Errorf("erorr while sending sdiffstore command: %w", err)
	}
	return resp, nil
}

func (r *RedisClient) Sinter(key string, keys ...string) (interface{}, error) {
	commands_args := []string{"SINTER", key}
	commands_args = append(commands_args, keys...)
	resp, err := r.Do(commands_args...)
	if err != nil {
		return nil, fmt.Errorf("erorr while sending sinter command: %w", err)
	}
	return resp, nil
}

// Need to think of some better way of encapsulating method signature
func (r *RedisClient) SinterCard(numcard int, keys []string, limit ...int) (interface{}, error) {
	commands_args := []string{"SINTERCARD", strconv.Itoa(numcard)}
	commands_args = append(commands_args, keys...)
	if len(limit) > 0 {
		commands_args = append(commands_args, "LIMIT", strconv.Itoa(limit[0]))
	}
	resp, err := r.Do(commands_args...)
	if err != nil {
		return nil, fmt.Errorf("erorr while sending sintercard command: %w", err)
	}
	return resp, nil
}

// gotta test it further
func (r *RedisClient) SinterStore(keys ...string) (interface{}, error) {
	commands_args := []string{"SINTERSTORE"}
	commands_args = append(commands_args, keys...)
	resp, err := r.Do(commands_args...)
	if err != nil {
		return nil, fmt.Errorf("erorr while sending sinterstore command: %w", err)
	}
	return resp, nil
}

func (r *RedisClient) SisMember(key string, member string) (interface{}, error) {
	resp, err := r.Do("SISMEMBER", key, member)
	if err != nil {
		return nil, fmt.Errorf("erorr while sending sismember command: %w", err)
	}
	return resp, nil
}

func (r *RedisClient) Smembers(key string) (interface{}, error) {
	resp, err := r.Do("SMEMBERS", key)
	if err != nil {
		return nil, fmt.Errorf("erorr while sending smembers command: %w", err)
	}
	return resp, nil
}

func (r *RedisClient) Smismember(key string, members ...string) (interface{}, error) {
	command_args := []string{"SMISMEMBER", key}
	command_args = append(command_args, members...)
	resp, err := r.Do(command_args...)
	if err != nil {
		return nil, fmt.Errorf("erorr while sending smismember command: %w", err)
	}
	return resp, nil
}

func (r *RedisClient) Smove(source string, destination string, member string) (interface{}, error) {
	resp, err := r.Do("SMOVE", source, destination, member)
	if err != nil {
		return nil, fmt.Errorf("erorr while sending smove command: %w", err)
	}
	return resp, nil
}

func (r *RedisClient) readResponse() (interface{}, error) {
	line, err := r.reader.ReadString('\n')

	if err != nil {
		return nil, err
	}
	switch line[0] {
	case '+':
		return strings.TrimSpace(line[1:]), nil
	case '-':
		return nil, errors.New(strings.TrimSpace(line[1:]))
	case ':':
		return strconv.ParseInt(strings.TrimSpace(line[1:]), 10, 64)
	case '$':
		return r.readBulkString(line)
	case '*':
		return r.readArray(line)
	default:
		return nil, fmt.Errorf("unknown response type: %s", string(line[0]))
	}
}

func (r *RedisClient) readBulkString(line string) (interface{}, error) {
	length, err := strconv.Atoi(strings.TrimSpace(line[1:]))
	if err != nil {
		return nil, fmt.Errorf("error while reading bulk string %w", err)
	}
	if length == -1 {
		return nil, nil
	}
	buf := make([]byte, length+2) // +2 for \r\n
	_, err = r.reader.Read(buf)
	if err != nil {
		return "", fmt.Errorf("error while reading bulk string %w", err)
	}
	return string(buf[:length]), nil
}

func (r *RedisClient) readArray(line string) ([]interface{}, error) {
	count, err := strconv.Atoi(strings.TrimSpace(line[1:]))
	if err != nil {
		return nil, err
	}
	if count == -1 {
		return nil, nil
	}
	array := make([]interface{}, count)
	for i := 0; i < count; i++ {
		array[i], err = r.readResponse()
		if err != nil {
			return nil, fmt.Errorf("error while reading array response %w", err)
		}
	}
	return array, nil
}

func (c *RedisClient) authenticate() error {

	if c.username != "" && c.password != "" {
		return c.authCommand("AUTH", c.username, c.password)
	} else if c.password != "" {
		return c.authCommand("AUTH", c.password)
	}
	return nil
}

func (c *RedisClient) authCommand(args ...string) error {
	resp, err := c.Do(args...)
	if err != nil {
		return fmt.Errorf("authentication failed: %v", err)
	}

	if str, ok := resp.(string); ok && str == "OK" {
		return nil
	}

	return fmt.Errorf("unexpected authentication response: %v", resp)
}
