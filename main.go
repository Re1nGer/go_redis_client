package main

import (
	"bufio"
	"errors"
	"fmt"
	"net"
	"strconv"
	"strings"
)

type RedisClient struct {
	conn     net.Conn
	reader   bufio.Reader
	host     string
	port     int
	username string
	password string
	ssl      bool
	//for now it suffices to have just these fields
}

func NewClient(host string, port int) (*RedisClient, error) {

	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", host, port))
	if err != nil {
		return nil, errors.New("failed to connect to redis")
	}
	return &RedisClient{
		host:   host,
		port:   port,
		conn:   conn,
		reader: *bufio.NewReader(conn),
	}, nil
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

func (r *RedisClient) sendCommand(comamnd []string) error {
	encoded_command := encodeCommand(comamnd)
	_, err := r.conn.Write(encoded_command)
	if err != nil {
		return err
	}
	return err
}

func (r *RedisClient) Do(command string, args ...string) (interface{}, error) {
	if err := r.sendCommand(args); err != nil {
		return nil, err
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
	resp, err := r.Do("EXISTS", args...)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (r *RedisClient) BuildArray(count int64) []byte {
	arr := make([]byte, 0)
	arr = append(arr, '*')
	arr = strconv.AppendInt(arr, count, 10)
	return append(arr, '\r', '\n')
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
		return nil, err
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

func main() {
	c, err := NewClient("localhost", 6379)
	if err != nil {
		return
	}
	arr1, err1 := c.Set("test", "val")
	if err1 != nil {
		return
	}

	arr, err := c.Get("test")
	if err != nil {
		return
	}
	fmt.Println("response", arr1, arr)
}
