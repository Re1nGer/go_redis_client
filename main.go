package main

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

type OptsFunc func(*Opts)

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
		return nil, fmt.Errorf("error while authenticating %w", auth_err)
	}

	return client, nil
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

func (r *RedisClient) Set(key string, val string) (interface{}, error) {
	resp, err := r.Do("SET", key, val)
	if err != nil {
		return nil, err
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

func (r *RedisClient) LPop(listname string) (interface{}, error) {
	resp, err := r.Do("LPOP", listname)
	if err != nil {
		return nil, fmt.Errorf("erorr while sending lpop command: %w", err)
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

func main() {
	c, err := NewClient("redis-15358.c92.us-east-1-3.ec2.redns.redis-cloud.com", 15358, WithPassword(""))
	if err != nil {
		fmt.Printf("error while connecting to redis %s", err)
		return
	}
	arr1, _ := c.Set("test", "vallllll")

	arr, err := c.Get("test")

	arr2, err2 := c.Exists("123", "3123", "124123")

	arr3, err3 := c.LPop("nonexistant_list")

	fmt.Println("response", arr1, arr, arr2, err2, err, arr3, err3)
}
