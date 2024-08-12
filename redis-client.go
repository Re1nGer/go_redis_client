package redisclient

import (
	"bufio"
	"context"
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
	host                    string
	port                    int
	db                      int
	username                string
	password                string
	ssl                     bool
	socketTimeout           time.Duration //timeout for
	socketConnectTimeout    time.Duration //timeout for connecting to indicated host and port
	socketKeepAliveInterval time.Duration
	connectionPool          *ConnectionPool
	encoding                string
	encodingErrors          string
	charset                 string
	errors                  []error
	decodeResponses         bool
	retryOnTimeout          bool
	retryOnError            []error
	sslKeyFile              string
	sslCertFile             string
	sslCertReqs             string
	sslCaCerts              string
	sslCaPath               string
	sslCaData               []byte
	sslCheckhostname        bool
	sslPassword             string
	sslValidateOcsp         bool
	sslValidateOcspstapled  bool
	sslOcspContext          interface{} // this might need a more specific type
	sslocspexpectedcert     []byte
	maxconnections          int
	singleConnectionClient  bool
	healthcheckInterval     time.Duration
	clientname              string
	libname                 string
	libversion              string
	retry                   *RetryOptions
	redisConnectFunc        RedisConnectFunc
	credentialProvider      CredentialProvider
	protocol                int
}

// TODO: Fill in
type RetryOptions struct{}
type CredentialProvider struct{}
type ConnectionPool struct{}
type RedisConnectFunc func() (net.Conn, error)

type CommandListOpts struct {
	filterby string
}

type BgsaveOpts struct {
	schedule bool
}

type ClientUnblockOpts struct {
	timeout bool
	error   bool
}

type BitcountOpts struct {
	start int
	end   int
	bit   string
}

type ClientUnblockOptsFunc func(*ClientUnblockOpts)
type ClientTrackingOption func(*clientTrackingOptions)

type AclCatOpts struct {
	cat string
}

type AclCatOptsFunc func(*AclCatOpts)

type clientTrackingOptions struct {
	on       bool
	redirect int
	prefix   []string
	bcast    bool
	optin    bool
	optout   bool
	noLoop   bool
}

type LCSOptions struct {
	LEN          bool
	IDX          bool
	MINMATCHLEN  int
	WITHMATCHLEN bool
}

// LCSOptsFunc is a function type for setting LCS options
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

type HExpireOpts struct {
	NX       bool
	XX       bool
	GT       bool
	LT       bool
	Duration *time.Duration
}

type LPopOpts struct {
	count int64
}

type SScanOpts struct {
	Match string
	Count int
}

type LPosOpts struct {
	rank   int64
	count  int64
	maxlen int64
}

type RPopOpts struct {
	count int64
}
type GetexOpts struct {
	EX      *time.Duration
	PX      *time.Duration
	EXAT    *time.Time
	PXAT    *time.Time
	Persist bool
}
type BitfieldOperation struct {
	Op       string // GET, SET, or INCRBY
	Encoding string // i for signed, u for unsigned, followed by bit count
	Offset   int    // Bit offset
	Value    int64  // Used for SET and INCRBY operations
}

type BitfieldRoOperation struct {
	Encoding string // i for signed, u for unsigned, followed by bit count
	Offset   int    // Bit offset
	Value    int64  // Used for SET and INCRBY operations
}

// BitfieldOptions represents additional options for the BITFIELD command
type BitfieldOptions struct {
	Overflow string // WRAP, SAT, or FAIL
}

type LPosOptsFunc func(*LPosOpts)
type LPopFunc func(*LPopOpts)
type RPopOptsFunc func(*RPopOpts)
type OptsFunc func(*Opts)
type SetOptsFunc func(*SetOpts)
type HExpireOptsFunc func(*HExpireOpts)
type SScanOptsFunc func(*SScanOpts)
type GetexOptsFunc func(*GetexOpts)
type LCSOptsFunc func(*LCSOptions)
type BitcountOptsFunc func(*BitcountOpts)
type BgsaveOptsFunc func(*BgsaveOpts)
type CommandListOptsFunc func(*CommandListOpts)

func WithFilterByCommandList(modifier string) CommandListOptsFunc {
	return func(opts *CommandListOpts) {
		opts.filterby = modifier
	}
}

func BitfieldGet(encoding string, offset int) BitfieldOperation {
	return BitfieldOperation{Op: "GET", Encoding: encoding, Offset: offset}
}

func BitfieldSet(encoding string, offset int, value int64) BitfieldOperation {
	return BitfieldOperation{Op: "SET", Encoding: encoding, Offset: offset, Value: value}
}

func BitfieldIncrBy(encoding string, offset int, increment int64) BitfieldOperation {
	return BitfieldOperation{Op: "INCRBY", Encoding: encoding, Offset: offset, Value: increment}
}

func WithTimeout() ClientUnblockOptsFunc {
	return func(opts *ClientUnblockOpts) {
		opts.timeout = true
	}
}

func WithError() ClientUnblockOptsFunc {
	return func(opts *ClientUnblockOpts) {
		opts.error = true
	}
}

func WithStartEnd(start int, end int) BitcountOptsFunc {
	return func(opts *BitcountOpts) {
		opts.start = start
		opts.end = end
	}
}
func WithOn() ClientTrackingOption {
	return func(o *clientTrackingOptions) {
		o.on = true
	}
}

func WithOff() ClientTrackingOption {
	return func(o *clientTrackingOptions) {
		o.on = false
	}
}

func WithRedirect(clientID int) ClientTrackingOption {
	return func(o *clientTrackingOptions) {
		o.redirect = clientID
	}
}

func WithPrefix(prefix ...string) ClientTrackingOption {
	return func(o *clientTrackingOptions) {
		o.prefix = append(o.prefix, prefix...)
	}
}

func WithBcast() ClientTrackingOption {
	return func(o *clientTrackingOptions) {
		o.bcast = true
	}
}

func WithOptin() ClientTrackingOption {
	return func(o *clientTrackingOptions) {
		o.optin = true
	}
}

func WithOptout() ClientTrackingOption {
	return func(o *clientTrackingOptions) {
		o.optout = true
	}
}

func WithNoLoop() ClientTrackingOption {
	return func(o *clientTrackingOptions) {
		o.noLoop = true
	}
}

func WithBit(bit string) BitcountOptsFunc {
	return func(opts *BitcountOpts) {
		opts.bit = bit
	}
}

func WithLen() LCSOptsFunc {
	return func(opts *LCSOptions) {
		opts.LEN = true
	}
}

func WithIDX() LCSOptsFunc {
	return func(opts *LCSOptions) {
		opts.IDX = true
	}
}

func WithMinMatchLen(length int) LCSOptsFunc {
	return func(opts *LCSOptions) {
		opts.MINMATCHLEN = length
	}
}

func WithMatchLen() LCSOptsFunc {
	return func(opts *LCSOptions) {
		opts.WITHMATCHLEN = true
	}
}

func WithEX(seconds time.Duration) GetexOptsFunc {
	return func(opts *GetexOpts) {
		opts.EX = &seconds
	}
}

func WithPX(milliseconds time.Duration) GetexOptsFunc {
	return func(opts *GetexOpts) {
		opts.PX = &milliseconds
	}
}

func WithEXAT(timestamp time.Time) GetexOptsFunc {
	return func(opts *GetexOpts) {
		opts.EXAT = &timestamp
	}
}

func WithPXAT(timestamp time.Time) GetexOptsFunc {
	return func(opts *GetexOpts) {
		opts.PXAT = &timestamp
	}
}
func WithPersist() GetexOptsFunc {
	return func(opts *GetexOpts) {
		opts.Persist = true
	}
}

func WithNXSet() SetOptsFunc {
	return func(o *SetOpts) {
		o.NX = true
		o.XX = false // NX and XX are mutually exclusive
	}
}

func WithXXSet() SetOptsFunc {
	return func(o *SetOpts) {
		o.XX = true
		o.NX = false // NX and XX are mutually exclusive
	}
}

func WithGet() SetOptsFunc {
	return func(o *SetOpts) {
		o.Get = true
	}
}

func WithEXSet(seconds int64) SetOptsFunc {
	return func(o *SetOpts) {
		o.EX = seconds
		o.PX = 0
		o.EXAT = nil
		o.PXAT = nil
		o.KeepTTL = false
	}
}

func WithPXSet(milliseconds int64) SetOptsFunc {
	return func(o *SetOpts) {
		o.PX = milliseconds
		o.EX = 0
		o.EXAT = nil
		o.PXAT = nil
		o.KeepTTL = false
	}
}

func WithEXATSet(timestamp time.Time) SetOptsFunc {
	return func(o *SetOpts) {
		o.EXAT = &timestamp
		o.EX = 0
		o.PX = 0
		o.PXAT = nil
		o.KeepTTL = false
	}
}

func WithPXATSet(timestamp time.Time) SetOptsFunc {
	return func(o *SetOpts) {
		o.PXAT = &timestamp
		o.EX = 0
		o.PX = 0
		o.EXAT = nil
		o.KeepTTL = false
	}
}

func WithKeepTTL() SetOptsFunc {
	return func(o *SetOpts) {
		o.KeepTTL = true
		o.EX = 0
		o.PX = 0
		o.EXAT = nil
		o.PXAT = nil
	}
}

func WithMatchSScan(pattern string) SScanOptsFunc {
	return func(opts *SScanOpts) {
		opts.Match = pattern
	}
}

func WithCountSScan(count int) SScanOptsFunc {
	return func(opts *SScanOpts) {
		opts.Count = count
	}
}

func NewSetOpts() *SetOpts {
	return &SetOpts{}
}

func NewHExpireOpts() HExpireOpts {
	return HExpireOpts{}
}

func WithNX() HExpireOptsFunc {
	return func(o *HExpireOpts) {
		o.NX = true
	}
}

func WithXX() HExpireOptsFunc {
	return func(o *HExpireOpts) {
		o.XX = true
	}
}

func WithGT() HExpireOptsFunc {
	return func(o *HExpireOpts) {
		o.GT = true
	}
}

func WithLT() HExpireOptsFunc {
	return func(o *HExpireOpts) {
		o.LT = true
	}
}

func WithCountLPop(count int64) LPopFunc {
	return func(o *LPopOpts) {
		o.count = count
	}
}

func WithCountLPos(count int64) LPosOptsFunc {
	return func(o *LPosOpts) {
		o.count = count
	}
}

func WithRankLPos(rank int64) LPosOptsFunc {
	return func(o *LPosOpts) {
		o.rank = rank
	}
}

func WithMaxLenLPos(maxlen int64) LPosOptsFunc {
	return func(o *LPosOpts) {
		o.maxlen = maxlen
	}
}

func WithCountRPop(count int64) RPopOptsFunc {
	return func(o *RPopOpts) {
		o.count = count
	}
}

func defaultLPosOpts() LPosOpts {
	return LPosOpts{
		count: -1, // since count can be set to 0 to find all matches in the list, default is set to -1
		//see https://redis.io/docs/latest/commands/lpos/
	}
}

func defaultLPopOpts() LPopOpts {
	return LPopOpts{}
}

func defaultRPopOpts() RPopOpts {
	return RPopOpts{}
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

func WithSocketConnectTimeout(duration time.Duration) OptsFunc {
	return func(o *Opts) {
		o.socketConnectTimeout = duration
	}
}

func WithSocketTimeout(duration time.Duration) OptsFunc {
	return func(o *Opts) {
		o.socketTimeout = duration
	}
}

func WithSocketKeepAlive(enabled bool, interval time.Duration) OptsFunc {
	return func(o *Opts) {
		o.socketKeepAliveInterval = interval
	}
}

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
		if defaultOpts.socketConnectTimeout > 0 && defaultOpts.socketTimeout > 0 {
			ctx, cancel := context.WithTimeout(context.Background(), defaultOpts.socketTimeout)
			defer cancel()
			dialer := &net.Dialer{
				Timeout: defaultOpts.socketConnectTimeout,
			}
			if defaultOpts.socketKeepAliveInterval > 0 {
				dialer.KeepAlive = defaultOpts.socketKeepAliveInterval
			}
			conn, err = dialer.DialContext(ctx, "tcp", fmt.Sprintf("%s:%d", host, port))
		} else {
			conn, err = net.Dial("tcp", fmt.Sprintf("%s:%d", host, port))
		}
	}

	if err != nil {
		return nil, fmt.Errorf("failed to connect to redis %w", err)
	}

	if defaultOpts.socketTimeout > 0 {
		err = conn.SetDeadline(time.Now().Add(defaultOpts.socketTimeout))
		if err != nil {
			conn.Close()
			return nil, fmt.Errorf("failed to set connection deadline: %w", err)
		}
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
		return nil, fmt.Errorf("%w", err)
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

// refactor
func (r *RedisClient) Set(key string, val string, opts ...SetOptsFunc) (interface{}, error) {
	args := []string{"SET", key, val}

	defaultOpts := &SetOpts{}

	for _, opt := range opts {
		opt(defaultOpts)
	}

	if defaultOpts.NX {
		args = append(args, "NX")
	} else if defaultOpts.XX {
		args = append(args, "XX")
	}

	if defaultOpts.Get {
		args = append(args, "GET")
	}

	if defaultOpts.EX != 0 {
		args = append(args, "EX", strconv.Itoa(int(defaultOpts.EX))) // set amount of seconds
	} else if defaultOpts.PX != 0 {
		args = append(args, "PX", strconv.FormatInt(int64(defaultOpts.PX), 10)) // set amount of milliseconds
	} else if defaultOpts.EXAT != nil {
		args = append(args, "EXAT", strconv.FormatInt(defaultOpts.EXAT.Unix(), 10)) //UNIX timestamp with seconds
	} else if defaultOpts.PXAT != nil {
		args = append(args, "PXAT", strconv.FormatInt(defaultOpts.PXAT.UnixNano()/int64(time.Millisecond), 10)) //UNIX timestamp with milliseconds
	} else if defaultOpts.KeepTTL {
		args = append(args, "KEEPTTL")
	}

	resp, err := r.Do(args...)

	if err != nil {
		return nil, err
	}

	if defaultOpts.Get {
		if resp == nil {
			return nil, nil
		}
		str, ok := resp.(string)
		if !ok {
			return nil, fmt.Errorf("unexpected response type for SET with GET: %T", resp)
		}
		return str, nil
	}

	if resp == nil {
		return nil, nil // SET NX/XX condition not met
	}
	_, ok := resp.(string)
	if !ok {
		return nil, fmt.Errorf("unexpected response type for SET: %T", resp)
	}
	return resp, nil
}

func (r *RedisClient) Append(key string, val string) (interface{}, error) {
	resp, err := r.Do("APPEND", key, val)
	if err != nil {
		return nil, fmt.Errorf("error while sending append command %w", err)
	}
	return resp, nil
}

func (r *RedisClient) Decr(key string) (interface{}, error) {
	resp, err := r.Do("DECR", key)
	if err != nil {
		return nil, fmt.Errorf("error while sending decr command %w", err)
	}
	return resp, nil
}

func (r *RedisClient) Decrby(key string, decrement int64) (interface{}, error) {
	resp, err := r.Do("DECRBY", key, strconv.Itoa(int(decrement)))
	if err != nil {
		return nil, fmt.Errorf("error while sending decrby command %w", err)
	}
	return resp, nil
}

func (r *RedisClient) Getdel(key string) (interface{}, error) {
	resp, err := r.Do("GETDEL", key)
	if err != nil {
		return nil, fmt.Errorf("error while sending getdel command %w", err)
	}
	return resp, nil
}

func (r *RedisClient) Getex(key string, opts ...GetexOptsFunc) (interface{}, error) {

	options := &GetexOpts{}

	for _, opt := range opts {
		opt(options)
	}

	args := []string{"GETEX", key}

	if options.EX != nil {
		args = append(args, "EX", strconv.FormatInt(int64(options.EX.Seconds()), 10))
	} else if options.PX != nil {
		args = append(args, "PX", strconv.FormatInt(options.PX.Milliseconds(), 10))
	} else if options.EXAT != nil {
		args = append(args, "EXAT", strconv.FormatInt(options.EXAT.Unix(), 10))
	} else if options.PXAT != nil {
		args = append(args, "PXAT", strconv.FormatInt(options.PXAT.UnixNano()/int64(time.Millisecond), 10))
	} else if options.Persist {
		args = append(args, "PERSIST")
	}

	resp, err := r.Do(args...)

	if err != nil {
		return nil, fmt.Errorf("error while sending GETEX command: %w", err)
	}
	return resp, nil
}

func (r *RedisClient) Getrange(key string, start int64, end int64) (interface{}, error) {
	resp, err := r.Do("GETRANGE", key, strconv.Itoa(int(start)), strconv.Itoa(int(end)))
	if err != nil {
		return nil, fmt.Errorf("error while sending getrange command %w", err)
	}
	return resp, nil
}

// As of redis v6.20 it's deprecated, see https://redis.io/docs/latest/commands/getset/
func (r *RedisClient) Getset(key string, value string) (interface{}, error) {
	resp, err := r.Do("GETSET", key, value)
	if err != nil {
		return nil, fmt.Errorf("error while sending getset command %w", err)
	}
	return resp, nil
}

func (r *RedisClient) Incr(key string) (interface{}, error) {
	resp, err := r.Do("INCR", key)
	if err != nil {
		return nil, fmt.Errorf("error while sending incr command %w", err)
	}
	return resp, nil
}

func (r *RedisClient) Incrby(key string, increment int64) (interface{}, error) {
	resp, err := r.Do("INCRBY", key, strconv.Itoa(int(increment)))
	if err != nil {
		return nil, fmt.Errorf("error while sending incrby command %w", err)
	}
	return resp, nil
}

func (r *RedisClient) Incrbyfloat(key string, increment int64) (interface{}, error) {
	resp, err := r.Do("INCRBYFLOAT", key, strconv.Itoa(int(increment)))
	if err != nil {
		return nil, fmt.Errorf("error while sending incrbyfloat command %w", err)
	}
	return resp, nil
}

// not type safe
func (r *RedisClient) LCS(key1 string, key2 string, opts ...LCSOptsFunc) (interface{}, error) {
	options := &LCSOptions{}
	for _, opt := range opts {
		opt(options)
	}

	args := []string{"LCS", key1, key2}

	if options.LEN {
		args = append(args, "LEN")
	}
	if options.IDX {
		args = append(args, "IDX")
	}
	if options.MINMATCHLEN > 0 {
		args = append(args, "MINMATCHLEN", strconv.Itoa(options.MINMATCHLEN))
	}
	if options.WITHMATCHLEN {
		args = append(args, "WITHMATCHLEN")
	}

	resp, err := r.Do(args...)

	if err != nil {
		return nil, fmt.Errorf("error while sending lcs command: %w", err)
	}

	return resp, nil
}

func (r *RedisClient) MGet(key string, keys ...string) (interface{}, error) {
	command_args := []string{"MGET", key}
	command_args = append(command_args, keys...)
	resp, err := r.Do(command_args...)
	if err != nil {
		return nil, fmt.Errorf("error while sending mget command %w", err)
	}
	return resp, nil
}

func (r *RedisClient) MSet(key string, value string, keyvalues ...string) (interface{}, error) {
	command_args := []string{"MSET", key, value}
	command_args = append(command_args, keyvalues...)
	resp, err := r.Do(command_args...)
	if err != nil {
		return nil, fmt.Errorf("error while sending mset command %w", err)
	}
	return resp, nil
}

func (r *RedisClient) Msetnx(key string, value string, keyvalues ...string) (interface{}, error) {
	command_args := []string{"MSETNX", key, value}
	command_args = append(command_args, keyvalues...)
	resp, err := r.Do(command_args...)
	if err != nil {
		return nil, fmt.Errorf("error while sending msetnx command %w", err)
	}
	return resp, nil
}

// deprecated
func (r *RedisClient) Psetex(key string, milliseconds int64, value string) (interface{}, error) {
	command_args := []string{"PSETEX", key, strconv.Itoa(int(milliseconds)), value}
	resp, err := r.Do(command_args...)
	if err != nil {
		return nil, fmt.Errorf("error while sending psetex command %w", err)
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

// deprecated
func (r *RedisClient) SetNx(key string, val string) (interface{}, error) {
	resp, err := r.Do("SETNX", key, val)
	if err != nil {
		return nil, fmt.Errorf("erorr while sending setnx command: %w", err)
	}
	return resp, nil
}

func (r *RedisClient) SetRange(key string, offset int64, val string) (interface{}, error) {
	resp, err := r.Do("SETRANGE", key, strconv.Itoa(int(offset)), val)
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

func (r *RedisClient) Substr(key string, start int64, end int64) (interface{}, error) {
	resp, err := r.Do("STRLEN", key, strconv.Itoa(int(start)), strconv.Itoa(int(end)))
	if err != nil {
		return nil, fmt.Errorf("erorr while sending substr command: %w", err)
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

// requires multiple args
func (r *RedisClient) Spop(key string) (interface{}, error) {
	resp, err := r.Do("SPOP", key)
	if err != nil {
		return nil, fmt.Errorf("erorr while sending spop command: %w", err)
	}
	return resp, nil
}

func (r *RedisClient) SpopWithCount(key string, count int) (interface{}, error) {
	resp, err := r.Do("SPOP", key, strconv.Itoa(count))
	if err != nil {
		return nil, fmt.Errorf("erorr while sending spop with count command: %w", err)
	}
	return resp, nil
}

func (r *RedisClient) SRandMember(key string) (interface{}, error) {
	resp, err := r.Do("SRANDMEMBER ", key)
	if err != nil {
		return nil, fmt.Errorf("erorr while sending srandmember command: %w", err)
	}
	return resp, nil
}

func (r *RedisClient) SRandMemberWithCount(key string, count int) (interface{}, error) {
	resp, err := r.Do("SRANDMEMBER ", key, strconv.Itoa(count))
	if err != nil {
		return nil, fmt.Errorf("erorr while sending srandmember command: %w", err)
	}
	return resp, nil
}

func (r *RedisClient) Srem(key string, members ...string) (interface{}, error) {
	commands_args := []string{"SREM", key}
	commands_args = append(commands_args, members...)
	resp, err := r.Do(commands_args...)

	if err != nil {
		return nil, fmt.Errorf("erorr while sending srem command: %w", err)
	}
	return resp, nil
}

func (r *RedisClient) Sscan(key string, cursor string, opts ...SScanOptsFunc) (interface{}, error) {

	options := &SScanOpts{}

	for _, opt := range opts {
		opt(options)
	}

	commands_args := []string{"SSCAN", key, cursor}

	if options.Match != "" {
		commands_args = append(commands_args, "MATCH", options.Match)
	}
	if options.Count > 0 {
		commands_args = append(commands_args, "COUNT", strconv.Itoa(options.Count))
	}

	resp, err := r.Do(commands_args...)

	if err != nil {
		return nil, fmt.Errorf("erorr while sending sscan command: %w", err)
	}
	return resp, nil
}

func (r *RedisClient) SUnion(key string, keys ...string) (interface{}, error) {
	commands_args := []string{"SUNION", key}
	commands_args = append(commands_args, keys...)
	resp, err := r.Do(commands_args...)

	if err != nil {
		return nil, fmt.Errorf("erorr while sending sunion command: %w", err)
	}
	return resp, nil
}

func (r *RedisClient) SUnionStore(destination string, key string, keys ...string) (interface{}, error) {
	commands_args := []string{"SUNIONSTORE", destination, key}
	commands_args = append(commands_args, keys...)
	resp, err := r.Do(commands_args...)

	if err != nil {
		return nil, fmt.Errorf("erorr while sending sunionstore command: %w", err)
	}
	return resp, nil
}

func (r *RedisClient) Del(key string, keys ...string) (interface{}, error) {
	commands_args := []string{"DEL", key}
	commands_args = append(commands_args, keys...)
	resp, err := r.Do(commands_args...)

	if err != nil {
		return nil, fmt.Errorf("erorr while sending delete ommand: %w", err)
	}
	return resp, nil
}

func (r *RedisClient) Dump(key string) (interface{}, error) {
	resp, err := r.Do("DUMP", key)
	if err != nil {
		return nil, fmt.Errorf("erorr while sending dump command: %w", err)
	}
	return resp, nil
}

func (r *RedisClient) HDel(key string, field string, fields ...string) (interface{}, error) {
	commands_args := []string{"HDEL", key, field}
	commands_args = append(commands_args, fields...)
	resp, err := r.Do(commands_args...)

	if err != nil {
		return nil, fmt.Errorf("erorr while sending hdel command: %w", err)
	}

	return resp, nil
}

func (r *RedisClient) HSet(key string, field string, value string, fieldvalueargs ...string) (interface{}, error) {
	commands_args := []string{"HSET", key, field, value}
	commands_args = append(commands_args, fieldvalueargs...)
	resp, err := r.Do(commands_args...)

	if err != nil {
		return nil, fmt.Errorf("erorr while sending hset command: %w", err)
	}

	return resp, nil
}

func (r *RedisClient) HGet(key string, field string) (interface{}, error) {
	resp, err := r.Do("HGET", key, field)
	if err != nil {
		return nil, fmt.Errorf("erorr while sending hget command: %w", err)
	}
	return resp, nil
}

func (r *RedisClient) HExists(key string, field string) (interface{}, error) {
	resp, err := r.Do("HEXISTS", key, field)
	if err != nil {
		return nil, fmt.Errorf("erorr while sending hexists command: %w", err)
	}
	return resp, nil
}

func (r *RedisClient) HGetAll(key string) (interface{}, error) {
	resp, err := r.Do("HGETALL", key)
	if err != nil {
		return nil, fmt.Errorf("erorr while sending hgetall command: %w", err)
	}
	return resp, nil
}

func (r *RedisClient) HLen(key string) (interface{}, error) {
	resp, err := r.Do("HLEN", key)
	if err != nil {
		return nil, fmt.Errorf("erorr while sending hlen command: %w", err)
	}
	return resp, nil
}

func (r *RedisClient) HExpireTime(key string, numFields int, field string, fields ...string) (interface{}, error) {
	commands_args := []string{"HEXPIRETIME", key, "FIELDS", strconv.Itoa(numFields), field}
	commands_args = append(commands_args, fields...)
	resp, err := r.Do(commands_args...)
	if err != nil {
		return nil, fmt.Errorf("erorr while sending hexpire time command: %w", err)
	}
	return resp, nil
}

func (r *RedisClient) HIncrby(key string, field string, increment string) (interface{}, error) {
	resp, err := r.Do("HINCRBY", key, field, increment)
	if err != nil {
		return nil, fmt.Errorf("erorr while sending hincrby command: %w", err)
	}
	return resp, nil
}

func (r *RedisClient) HIncrbyfloat(key string, field string, increment string) (interface{}, error) {
	resp, err := r.Do("HINCRBYFLOAT", key, field, increment)
	if err != nil {
		return nil, fmt.Errorf("erorr while sending hincrbyfloat command: %w", err)
	}
	return resp, nil
}

func (r *RedisClient) HMGet(key string, field string, fields ...string) (interface{}, error) {
	commands_args := []string{"HMGET", key, field}
	commands_args = append(commands_args, fields...)
	resp, err := r.Do(commands_args...)
	if err != nil {
		return nil, fmt.Errorf("erorr while sending hmget command: %w", err)
	}
	return resp, nil
}

func (r *RedisClient) HPersist(key string, numFields int, field string, fields ...string) (interface{}, error) {
	commands_args := []string{"HPERSIST", key, "FIELDS", strconv.Itoa(numFields), field}
	commands_args = append(commands_args, fields...)
	resp, err := r.Do(commands_args...)
	if err != nil {
		return nil, fmt.Errorf("erorr while sending hpersist command: %w", err)
	}
	return resp, nil
}

func (r *RedisClient) HStrlen(key string, field string) (interface{}, error) {
	resp, err := r.Do("HSTRLEN", key, field)
	if err != nil {
		return nil, fmt.Errorf("erorr while sending hstrlen command: %w", err)
	}
	return resp, nil
}

// command doesn't exist ?
func (r *RedisClient) HExpire(key string, seconds int64, fields []string, opts ...HExpireOptsFunc) (interface{}, error) {

	var args []string

	defaultOpts := NewHExpireOpts()

	for _, fn := range opts {
		fn(&defaultOpts)
	}

	args = append(args, "HEXPIRE", key)

	args = append(args, strconv.FormatInt(seconds, 10))

	if defaultOpts.NX {
		args = append(args, "NX")
	}
	if defaultOpts.XX {
		args = append(args, "XX")
	}
	if defaultOpts.GT {
		args = append(args, "GT")
	}
	if defaultOpts.LT {
		args = append(args, "LT")
	}

	args = append(args, "FIELDS", strconv.Itoa(len(fields)))

	args = append(args, fields...)

	resp, err := r.Do(args...)

	if err != nil {
		return 0, fmt.Errorf("error while sending HEXPIRE command: %w", err)
	}

	return resp, nil
}

// needs futher investigation
func (r *RedisClient) BLMove(source string, destination string, wherefrom string, whereto string, timeout int64) (interface{}, error) {
	resp, err := r.Do("BLMOVE", source, destination, wherefrom, whereto, strconv.Itoa(int(timeout)))
	if err != nil {
		return nil, fmt.Errorf("erorr while sending blmove command: %w", err)
	}
	return resp, nil
}

func (r *RedisClient) LIndex(key string, index string) (interface{}, error) {
	resp, err := r.Do("LINDEX", key, index)
	if err != nil {
		return nil, fmt.Errorf("erorr while sending lindex command: %w", err)
	}
	return resp, nil
}

func (r *RedisClient) LInsert(key string, where string, pivot string, element string) (interface{}, error) {
	resp, err := r.Do("LINSERT", key, where, pivot, element)
	if err != nil {
		return nil, fmt.Errorf("erorr while sending linsert command: %w", err)
	}
	return resp, nil
}

func (r *RedisClient) LLen(key string) (interface{}, error) {
	resp, err := r.Do("LLEN", key)
	if err != nil {
		return nil, fmt.Errorf("erorr while sending llen command: %w", err)
	}
	return resp, nil
}

func (r *RedisClient) LMove(source string, destination string, wheresource string, wheredestination string) (interface{}, error) {
	resp, err := r.Do("LMOVE", source, destination, wheresource, wheredestination)
	if err != nil {
		return nil, fmt.Errorf("erorr while sending lmove command: %w", err)
	}
	return resp, nil
}

// PUT ON HOLD, LOTS OF OPTIONAL PARAMETERS
func (r *RedisClient) LMPop(source string, destination string, wheresource string, wheredestination string) (interface{}, error) {
	resp, err := r.Do("LMPOP", source, destination, wheresource, wheredestination)
	if err != nil {
		return nil, fmt.Errorf("erorr while sending lmove command: %w", err)
	}
	return resp, nil
}

func (r *RedisClient) LPop(listname string, opts ...LPopFunc) (interface{}, error) {

	defaultOpts := defaultLPopOpts()

	for _, fn := range opts {
		fn(&defaultOpts)
	}

	commands_args := []string{"LPOP", listname}

	if defaultOpts.count > 0 {
		commands_args = append(commands_args, strconv.Itoa(int(defaultOpts.count)))
	}

	resp, err := r.Do(commands_args...)

	if err != nil {
		return nil, fmt.Errorf("erorr while sending lpop command: %w", err)
	}
	return resp, nil
}

func (r *RedisClient) LPos(key string, element string, opts ...LPosOptsFunc) (interface{}, error) {

	defaultOpts := defaultLPosOpts()

	for _, fn := range opts {
		fn(&defaultOpts)
	}

	commands_args := []string{"LPOS", key, element}

	if defaultOpts.rank != 0 {
		commands_args = append(commands_args, "RANK", strconv.Itoa(int(defaultOpts.rank)))
	}
	if defaultOpts.count != -1 {
		commands_args = append(commands_args, "COUNT", strconv.Itoa(int(defaultOpts.count)))
	}
	if defaultOpts.maxlen != 0 {
		commands_args = append(commands_args, "MAXLEN", strconv.Itoa(int(defaultOpts.maxlen)))
	}

	resp, err := r.Do(commands_args...)

	if err != nil {
		return nil, fmt.Errorf("erorr while sending lpos command: %w", err)
	}
	return resp, nil
}

func (r *RedisClient) LPush(key string, element string, elements ...string) (interface{}, error) {
	commands_args := []string{"LPUSH", key, element}

	commands_args = append(commands_args, elements...)

	resp, err := r.Do(commands_args...)

	if err != nil {
		return nil, fmt.Errorf("erorr while sending lpush command: %w", err)
	}
	return resp, nil
}

func (r *RedisClient) LPushx(key string, element string, elements ...string) (interface{}, error) {
	commands_args := []string{"LPUSHX", key, element}

	commands_args = append(commands_args, elements...)

	resp, err := r.Do(commands_args...)

	if err != nil {
		return nil, fmt.Errorf("erorr while sending lpushx command: %w", err)
	}
	return resp, nil
}

func (r *RedisClient) LRange(key string, start string, stop string) (interface{}, error) {
	commands_args := []string{"LRANGE", key, start, stop}

	resp, err := r.Do(commands_args...)

	if err != nil {
		return nil, fmt.Errorf("erorr while sending lrange command: %w", err)
	}
	return resp, nil
}

func (r *RedisClient) LRem(key string, start int64, element string) (interface{}, error) {
	commands_args := []string{"LREM", key, strconv.Itoa(int(start)), element}

	resp, err := r.Do(commands_args...)

	if err != nil {
		return nil, fmt.Errorf("erorr while sending lrem command: %w", err)
	}
	return resp, nil
}

func (r *RedisClient) LSet(key string, index int64, element string) (interface{}, error) {
	commands_args := []string{"LSET", key, strconv.Itoa(int(index)), element}

	resp, err := r.Do(commands_args...)

	if err != nil {
		return nil, fmt.Errorf("erorr while sending lset command: %w", err)
	}
	return resp, nil
}

func (r *RedisClient) LTrim(key string, start int64, stop int64) (interface{}, error) {
	commands_args := []string{"LTRIM", key, strconv.Itoa(int(start)), strconv.Itoa(int(stop))}

	resp, err := r.Do(commands_args...)

	if err != nil {
		return nil, fmt.Errorf("erorr while sending ltrim command: %w", err)
	}
	return resp, nil
}

func (r *RedisClient) RPop(key string, opts ...RPopOptsFunc) (interface{}, error) {

	defaultOpts := defaultRPopOpts()

	for _, fn := range opts {
		fn(&defaultOpts)
	}

	commands_args := []string{"RPOP", key}

	if defaultOpts.count > 0 {
		commands_args = append(commands_args, strconv.Itoa(int(defaultOpts.count)))
	}

	resp, err := r.Do(commands_args...)

	if err != nil {
		return nil, fmt.Errorf("erorr while sending rpop command: %w", err)
	}
	return resp, nil
}

// deprecated as of redis v6.2.0, see https://redis.io/docs/latest/commands/rpoplpush/
func (r *RedisClient) RPopLPush(source string, destination string) (interface{}, error) {
	commands_args := []string{"RPOPLPUSH", source, destination}

	resp, err := r.Do(commands_args...)

	if err != nil {
		return nil, fmt.Errorf("erorr while sending rpoplpush command: %w", err)
	}
	return resp, nil
}

func (r *RedisClient) RPush(key string, element string, elements ...string) (interface{}, error) {
	commands_args := []string{"RPUSH", key, element}

	commands_args = append(commands_args, elements...)

	resp, err := r.Do(commands_args...)

	if err != nil {
		return nil, fmt.Errorf("erorr while sending rpush command: %w", err)
	}
	return resp, nil
}

func (r *RedisClient) RPushx(key string, element string, elements ...string) (interface{}, error) {
	commands_args := []string{"RPUSHX", key, element}

	commands_args = append(commands_args, elements...)

	resp, err := r.Do(commands_args...)

	if err != nil {
		return nil, fmt.Errorf("erorr while sending rpushx command: %w", err)
	}
	return resp, nil
}

func (r *RedisClient) Bitcount(key string, opts ...BitcountOptsFunc) (interface{}, error) {

	defaultOpts := BitcountOpts{
		start: -1,
		end:   -1,
	}

	for _, fn := range opts {
		fn(&defaultOpts)
	}

	commands_args := []string{"BITCOUNT", key}

	if defaultOpts.start != -1 && defaultOpts.end != -1 {
		commands_args = append(commands_args, strconv.Itoa(defaultOpts.start), strconv.Itoa(defaultOpts.end))
	}

	if defaultOpts.bit == "BYTE" {
		commands_args = append(commands_args, "BYTE")
	}

	if defaultOpts.bit == "BITE" {
		commands_args = append(commands_args, "BITE")
	}

	//handle bit/byte

	resp, err := r.Do(commands_args...)

	if err != nil {
		return nil, fmt.Errorf("erorr while sending rpushx command: %w", err)
	}
	return resp, nil
}

func (r *RedisClient) Bitfield(key string, operations []BitfieldOperation, opts *BitfieldOptions) (interface{}, error) {
	args := []string{"BITFIELD", key}

	for _, op := range operations {
		args = append(args, op.Op, op.Encoding, strconv.Itoa(op.Offset))
		if op.Op == "SET" || op.Op == "INCRBY" {
			args = append(args, strconv.FormatInt(op.Value, 10))
		}
	}

	if opts != nil && opts.Overflow != "" {
		args = append(args, "OVERFLOW", opts.Overflow)
	}

	resp, err := r.Do(args...)
	if err != nil {
		return nil, fmt.Errorf("error while sending BITFIELD command: %w", err)
	}

	return resp, nil
}

func (r *RedisClient) Bitfieldro(key string, operations []BitfieldRoOperation) (interface{}, error) {
	args := []string{"BITFIELD_RO", key}

	for _, op := range operations {
		args = append(args, op.Encoding, strconv.Itoa(op.Offset))
	}

	resp, err := r.Do(args...)
	if err != nil {
		return nil, fmt.Errorf("error while sending bitfield_ro command: %w", err)
	}

	return resp, nil
}

// bitwiseop <AND | OR | XOR | NOT>, see https://redis.io/docs/latest/commands/bitop/
func (r *RedisClient) Bitop(bitwiseop string, keys ...string) (interface{}, error) {
	args := []string{"BITOP", bitwiseop}

	args = append(args, keys...)

	resp, err := r.Do(args...)

	if err != nil {
		return nil, fmt.Errorf("error while sending bitop command: %w", err)
	}

	return resp, nil
}

func (r *RedisClient) Getbit(key string, offset int64) (interface{}, error) {

	resp, err := r.Do("GETBIT", key, strconv.Itoa(int(offset)))

	if err != nil {
		return nil, fmt.Errorf("error while sending getbit command: %w", err)
	}

	return resp, nil
}

func (r *RedisClient) Setbit(key string, offset int64, value string) (interface{}, error) {

	resp, err := r.Do("SETBIT", key, strconv.Itoa(int(offset)), value)

	if err != nil {
		return nil, fmt.Errorf("error while sending setbit command: %w", err)
	}

	return resp, nil
}

// TODO: implement BITPOS
func (r *RedisClient) Copy(key string, destination string) (interface{}, error) {

	resp, err := r.Do("COPY", key, destination)

	if err != nil {
		return nil, fmt.Errorf("error while sending copy command: %w", err)
	}

	return resp, nil
}

func (r *RedisClient) Keys(pattern string) (interface{}, error) {

	resp, err := r.Do("KEYS", pattern)

	if err != nil {
		return nil, fmt.Errorf("error while sending keys command: %w", err)
	}

	return resp, nil
}

func (r *RedisClient) Move(key string, db string) (interface{}, error) {

	resp, err := r.Do("MOVE", key, db)

	if err != nil {
		return nil, fmt.Errorf("error while sending move command: %w", err)
	}

	return resp, nil
}

func (r *RedisClient) Persist(key string) (interface{}, error) {

	resp, err := r.Do("PERSIST", key)

	if err != nil {
		return nil, fmt.Errorf("error while sending persist command: %w", err)
	}

	return resp, nil
}

func (r *RedisClient) Clientcaching(caching string) (interface{}, error) {

	resp, err := r.Do("CLIENT", "CACHING", caching)

	if err != nil {
		return nil, fmt.Errorf("error while sending clientcaching command: %w", err)
	}

	return resp, nil
}

func (r *RedisClient) Clientgetname() (interface{}, error) {

	args := []string{"CLIENT", "GETNAME"}

	resp, err := r.Do(args...)

	if err != nil {
		return nil, fmt.Errorf("error while sending clientgetname command: %w", err)
	}

	return resp, nil
}

func (r *RedisClient) Clientgetredir() (interface{}, error) {

	resp, err := r.Do("CLIENT", "GETREDIR")

	if err != nil {
		return nil, fmt.Errorf("error while sending clientgetredir command: %w", err)
	}

	return resp, nil
}

func (r *RedisClient) Clientid() (interface{}, error) {

	resp, err := r.Do("CLIENT", "ID")

	if err != nil {
		return nil, fmt.Errorf("error while sending clientid command: %w", err)
	}

	return resp, nil
}

func (r *RedisClient) Clientinfo() (interface{}, error) {

	resp, err := r.Do("CLIENT", "INFO")

	if err != nil {
		return nil, fmt.Errorf("error while sending clientinfo command: %w", err)
	}

	return resp, nil
}

func (r *RedisClient) Clientkill(args ...string) (interface{}, error) {

	resp, err := r.Do("CLIENT", "KILL")

	if err != nil {
		return nil, fmt.Errorf("error while sending clientkill command: %w", err)
	}

	return resp, nil
}

func (r *RedisClient) Clientlist(args ...string) (interface{}, error) {

	resp, err := r.Do("CLIENT", "LIST")

	if err != nil {
		return nil, fmt.Errorf("error while sending clientlist command: %w", err)
	}

	return resp, nil
}

func (r *RedisClient) Clientnoevict(evict string) (interface{}, error) {

	resp, err := r.Do("CLIENT", "NO-EVICT", evict)

	if err != nil {
		return nil, fmt.Errorf("error while sending clientnoevict command: %w", err)
	}

	return resp, nil
}

func (r *RedisClient) Clientnotouch(touch string) (interface{}, error) {

	resp, err := r.Do("CLIENT", "NO-TOUCH", touch)

	if err != nil {
		return nil, fmt.Errorf("error while sending clientnotouch command: %w", err)
	}

	return resp, nil
}

// investigate
func (r *RedisClient) Clientpause(mode string) (interface{}, error) {

	resp, err := r.Do("CLIENT", "PAUSE", "timeout", mode)

	if err != nil {
		return nil, fmt.Errorf("error while sending clientpause command: %w", err)
	}

	return resp, nil
}

func (r *RedisClient) Clientreply(mode string) (interface{}, error) {

	resp, err := r.Do("CLIENT", "REPLY", mode)

	if err != nil {
		return nil, fmt.Errorf("error while sending clientreply command: %w", err)
	}

	return resp, nil
}

func (r *RedisClient) Clientsetinfo(attribute string) (interface{}, error) {

	resp, err := r.Do("CLIENT", "SETINFO", attribute)

	if err != nil {
		return nil, fmt.Errorf("error while sending clientsetinfo command: %w", err)
	}

	return resp, nil
}

func (r *RedisClient) Clientsetname(connectionname string) (interface{}, error) {

	resp, err := r.Do("CLIENT", "SETNAME", connectionname)

	if err != nil {
		return nil, fmt.Errorf("error while sending connectionname command: %w", err)
	}

	return resp, nil
}

func (r *RedisClient) Clienttracking(opts ...ClientTrackingOption) (interface{}, error) {

	options := &clientTrackingOptions{}
	for _, opt := range opts {
		opt(options)
	}

	command_args := []string{"CLIENT", "TRACKING"}

	if options.on {
		command_args = append(command_args, "ON")
	} else {
		command_args = append(command_args, "OFF")
	}

	if options.redirect != 0 {
		command_args = append(command_args, "REDIRECT", strconv.Itoa(options.redirect))
	}

	for _, prefix := range options.prefix {
		command_args = append(command_args, "PREFIX", prefix)
	}

	if options.bcast {
		command_args = append(command_args, "BCAST")
	}

	if options.optin {
		command_args = append(command_args, "OPTIN")
	}

	if options.optout {
		command_args = append(command_args, "OPTOUT")
	}

	if options.noLoop {
		command_args = append(command_args, "NOLOOP")
	}

	resp, err := r.Do(command_args...)

	if err != nil {
		return nil, fmt.Errorf("error while sending clienttracking command: %w", err)
	}

	return resp, nil
}

func (r *RedisClient) Clienttrackinginfo() (interface{}, error) {

	resp, err := r.Do("CLIENT", "TRACKINGINFO")

	if err != nil {
		return nil, fmt.Errorf("error while sending clienttrackinginfo command: %w", err)
	}

	return resp, nil
}

func (r *RedisClient) Clientunblock(clientid int64, opts ...ClientUnblockOptsFunc) (interface{}, error) {

	defaultOpts := &ClientUnblockOpts{}

	for _, opt := range opts {
		opt(defaultOpts)
	}

	commands_args := []string{"CLIENT", "UNBLOCK ", strconv.Itoa(int(clientid))}

	if defaultOpts.error {
		commands_args = append(commands_args, "ERROR")
	}

	if defaultOpts.timeout {
		commands_args = append(commands_args, "TIMEOUT")
	}

	resp, err := r.Do(commands_args...)

	if err != nil {
		return nil, fmt.Errorf("error while sending clientunblock command: %w", err)
	}

	return resp, nil
}

func (r *RedisClient) Clientunpause() (interface{}, error) {

	resp, err := r.Do("CLIENT", "UNPAUSE")

	if err != nil {
		return nil, fmt.Errorf("error while sending clientunpause command: %w", err)
	}

	return resp, nil
}

func (r *RedisClient) Ping() (interface{}, error) {

	commands_args := []string{"PING"}

	resp, err := r.Do(commands_args...)

	if err != nil {
		return nil, fmt.Errorf("error while sending ping command: %w", err)
	}

	return resp, nil
}

// As of Redis v7.2.0 the command is deprecated, see https://redis.io/docs/latest/commands/quit/
func (r *RedisClient) Quit(message string) (interface{}, error) {

	resp, err := r.Do("QUIT", message)

	if err != nil {
		return nil, fmt.Errorf("error while sending quit command: %w", err)
	}

	return resp, nil
}

// Transaction commands
func (r *RedisClient) Multi() error {
	_, err := r.Do("MULTI")
	return err
}

func (r *RedisClient) Exec() (interface{}, error) {
	return r.Do("EXEC")
}

func (r *RedisClient) Discard() error {
	_, err := r.Do("DISCARD")
	return err
}

func (r *RedisClient) Watch(keys ...string) error {
	args := append([]string{"WATCH"}, keys...)
	_, err := r.Do(args...)
	return err
}

func (r *RedisClient) Unwatch() error {
	_, err := r.Do("UNWATCH")
	return err
}

func (r *RedisClient) TxPipeline() *TxPipeline {
	return &TxPipeline{client: r}
}

type TxPipeline struct {
	client *RedisClient
	cmds   [][]string
}

func (p *TxPipeline) Queue(args ...string) {
	p.cmds = append(p.cmds, args)
}

func (p *TxPipeline) Exec() ([]interface{}, error) {
	if err := p.client.Multi(); err != nil {
		return nil, err
	}

	for _, cmd := range p.cmds {
		if _, err := p.client.Do(cmd...); err != nil {
			p.client.Discard()
			return nil, err
		}
	}

	result, err := p.client.Exec()

	if err != nil {
		return nil, err
	}

	//for now that's just an array of any
	results, ok := result.([]interface{})
	if !ok {
		return nil, fmt.Errorf("unexpected result type from EXEC")
	}

	return results, nil
}

// stream commands
func (r *RedisClient) Xack(key string, group string, id string, ids ...string) (interface{}, error) {

	commands_args := []string{"XACK", key, group, id}

	commands_args = append(commands_args, ids...)

	resp, err := r.Do(commands_args...)

	if err != nil {
		return nil, fmt.Errorf("error while sending xack command: %w", err)
	}

	return resp, nil
}

func (r *RedisClient) Xackdel(key string, group string, id string, ids ...string) (interface{}, error) {

	commands_args := []string{"XDEL", key, group, id}

	commands_args = append(commands_args, ids...)

	resp, err := r.Do(commands_args...)

	if err != nil {
		return nil, fmt.Errorf("error while sending xdel command: %w", err)
	}

	return resp, nil
}

//addX xtrim
//update

//server management

func (r *RedisClient) Aclcat(opts ...AclCatOptsFunc) (interface{}, error) {

	options := &AclCatOpts{}

	for _, opt := range opts {
		opt(options)
	}

	command_args := []string{"ACL", "CAT"}

	if options.cat != "" {
		command_args = append(command_args, options.cat)
	}

	resp, err := r.Do(command_args...)

	if err != nil {
		return nil, fmt.Errorf("error while sending aclcat command: %w", err)
	}

	return resp, nil
}

func (r *RedisClient) Acldeluser(username string, usernames ...string) (interface{}, error) {

	command_args := []string{"ACL", "DELUSER", username}

	command_args = append(command_args, usernames...)

	resp, err := r.Do(command_args...)

	if err != nil {
		return nil, fmt.Errorf("error while sending acldeluser command: %w", err)
	}

	return resp, nil
}

func (r *RedisClient) Acldryrun(username string, command string, args ...string) (interface{}, error) {

	command_args := []string{"ACL", "DRYRUN", username, command}

	command_args = append(command_args, args...)

	resp, err := r.Do(command_args...)

	if err != nil {
		return nil, fmt.Errorf("error while sending acldryrun command: %w", err)
	}

	return resp, nil
}

func (r *RedisClient) Aclgenpass(bits string) (interface{}, error) {

	command_args := []string{"ACL", "GENPASS", bits}

	resp, err := r.Do(command_args...)

	if err != nil {
		return nil, fmt.Errorf("error while sending aclgenpass command: %w", err)
	}

	return resp, nil
}

//TODO: add ci/cd tests on every build

func (r *RedisClient) Aclgetuser(username string) (interface{}, error) {

	command_args := []string{"ACL", "GETUSER", username}

	resp, err := r.Do(command_args...)

	if err != nil {
		return nil, fmt.Errorf("error while sending aclgetuser command: %w", err)
	}

	return resp, nil
}

func (r *RedisClient) Acllist() (interface{}, error) {

	command_args := []string{"ACL", "LIST"}

	resp, err := r.Do(command_args...)

	if err != nil {
		return nil, fmt.Errorf("error while sending acllist command: %w", err)
	}

	return resp, nil
}

func (r *RedisClient) Aclload() (interface{}, error) {

	command_args := []string{"ACL", "LOAD"}

	resp, err := r.Do(command_args...)

	if err != nil {
		return nil, fmt.Errorf("error while sending aclload command: %w", err)
	}

	return resp, nil
}

func (r *RedisClient) Aclsave() (interface{}, error) {

	command_args := []string{"ACL", "SAVE"}

	resp, err := r.Do(command_args...)

	if err != nil {
		return nil, fmt.Errorf("error while sending aclsave command: %w", err)
	}

	return resp, nil
}

func (r *RedisClient) Aclsetuser(username string, rules ...string) (interface{}, error) {

	command_args := []string{"ACL", "SETUSER", username}

	command_args = append(command_args, rules...)

	resp, err := r.Do(command_args...)

	if err != nil {
		return nil, fmt.Errorf("error while sending aclsetuser command: %w", err)
	}

	return resp, nil
}

func (r *RedisClient) Aclusers() (interface{}, error) {

	command_args := []string{"ACL", "USERS"}

	resp, err := r.Do(command_args...)

	if err != nil {
		return nil, fmt.Errorf("error while sending aclusers command: %w", err)
	}

	return resp, nil
}

func (r *RedisClient) Aclwhoami() (interface{}, error) {

	command_args := []string{"ACL", "WHOAMI"}

	resp, err := r.Do(command_args...)

	if err != nil {
		return nil, fmt.Errorf("error while sending aclwhoami command: %w", err)
	}

	return resp, nil
}

func (r *RedisClient) Bgrewriteaof() (interface{}, error) {

	resp, err := r.Do("BGREWRITEAOF")

	if err != nil {
		return nil, fmt.Errorf("error while sending BGREWRITEAOF command: %w", err)
	}

	return resp, nil
}

func (r *RedisClient) Bgsave(opts ...BgsaveOptsFunc) (interface{}, error) {

	options := &BgsaveOpts{}

	for _, opt := range opts {
		opt(options)
	}

	command_args := []string{"BGSAVE"}

	if options.schedule {
		command_args = append(command_args, "SCHEDULE")
	}

	resp, err := r.Do(command_args...)

	if err != nil {
		return nil, fmt.Errorf("error while sending bgsave command: %w", err)
	}

	return resp, nil
}

func (r *RedisClient) Command(subcommands ...string) (interface{}, error) {

	command_args := []string{"COMMAND"}

	command_args = append(command_args, subcommands...)

	resp, err := r.Do(command_args...)

	if err != nil {
		return nil, fmt.Errorf("error while sending command command: %w", err)
	}

	return resp, nil
}

func (r *RedisClient) Commandcount() (interface{}, error) {

	command_args := []string{"COMMAND", "COUNT"}

	resp, err := r.Do(command_args...)

	if err != nil {
		return nil, fmt.Errorf("error while sending command count command: %w", err)
	}

	return resp, nil
}

func (r *RedisClient) Commanddocs(commands ...string) (interface{}, error) {

	command_args := []string{"COMMAND", "DOCS"}

	command_args = append(command_args, commands...)

	resp, err := r.Do(command_args...)

	if err != nil {
		return nil, fmt.Errorf("error while sending command docs command: %w", err)
	}

	return resp, nil
}

func (r *RedisClient) Commandgetkeys(command string, args ...string) (interface{}, error) {

	command_args := []string{"COMMAND", "GETKEYS", command}

	command_args = append(command_args, args...)

	resp, err := r.Do(command_args...)

	if err != nil {
		return nil, fmt.Errorf("error while sending command getkeys docs command: %w", err)
	}

	return resp, nil
}

func (r *RedisClient) Commandgetkeyssandflags(command string, args ...string) (interface{}, error) {

	command_args := []string{"COMMAND", "GETKEYSANDFLAGS", command}

	command_args = append(command_args, args...)

	resp, err := r.Do(command_args...)

	if err != nil {
		return nil, fmt.Errorf("error while sending command getkeyssandflags docs command: %w", err)
	}

	return resp, nil
}

func (r *RedisClient) Commandinfo(commandnames ...string) (interface{}, error) {

	command_args := []string{"COMMAND", "INFO"}

	command_args = append(command_args, commandnames...)

	resp, err := r.Do(command_args...)

	if err != nil {
		return nil, fmt.Errorf("error while sending command command info docs command: %w", err)
	}

	return resp, nil
}

func (r *RedisClient) Commandlist(opts ...CommandListOptsFunc) (interface{}, error) {

	options := &CommandListOpts{}

	for _, opt := range opts {
		opt(options)
	}

	command_args := []string{"COMMAND", "LIST"}

	if options.filterby != "" {
		command_args = append(command_args, "FILTERBY", options.filterby)
	}

	resp, err := r.Do(command_args...)

	if err != nil {
		return nil, fmt.Errorf("error while sending command command list docs command: %w", err)
	}

	return resp, nil
}

func (r *RedisClient) Configget(parameter string, parameters ...string) (interface{}, error) {

	command_args := []string{"CONFIG", "GET", parameter}

	command_args = append(command_args, parameters...)

	resp, err := r.Do(command_args...)

	if err != nil {
		return nil, fmt.Errorf("error while sending command config get docs command: %w", err)
	}

	return resp, nil
}

func (r *RedisClient) Configresetstat() (interface{}, error) {

	command_args := []string{"CONFIG", "RESETSTAT"}

	resp, err := r.Do(command_args...)

	if err != nil {
		return nil, fmt.Errorf("error while sending command config resetstat docs command: %w", err)
	}

	return resp, nil
}

func (r *RedisClient) Configrewrite() (interface{}, error) {

	command_args := []string{"CONFIG", "REWRITE"}

	resp, err := r.Do(command_args...)

	if err != nil {
		return nil, fmt.Errorf("error while sending command config rewrite docs command: %w", err)
	}

	return resp, nil
}

func (r *RedisClient) Configset(parameter string, value string, parameters ...string) (interface{}, error) {

	command_args := []string{"CONFIG", "SET", parameter, value}

	command_args = append(command_args, parameters...)

	resp, err := r.Do(command_args...)

	if err != nil {
		return nil, fmt.Errorf("error while sending command config set  docs command: %w", err)
	}

	return resp, nil
}
