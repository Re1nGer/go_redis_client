package main

import (
	"bufio"
	"errors"
	"fmt"
	"net"
	"strconv"
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
		return nil, errors.New("Failed To Connect To Redis")
	}
	return &RedisClient{
		host:   host,
		port:   port,
		conn:   conn,
		reader: *bufio.NewReader(conn),
	}, nil
}

func (r *RedisClient) Get(key string) []byte {
	arr := r.BuildArray(2)
	return r.BuildGet(arr, key)
}

func (r *RedisClient) Set(key string, val string) []byte {
	arr := r.BuildArray(3)
	return r.BuildSet(arr, key, val)
}

func (r *RedisClient) BuildArray(count int64) []byte {
	arr := make([]byte, 0)
	arr = append(arr, '*')
	arr = strconv.AppendInt(arr, count, 10)
	return append(arr, '\r', '\n')
}

// gotta refactor this shit
func (r *RedisClient) BuildGet(arr []byte, key string) []byte {
	arr = AppendGetCommand(arr)
	l := len(key)
	arr = append(arr, '$')
	arr = strconv.AppendInt(arr, int64(l), 10)
	arr = append(arr, '\r', '\n')
	arr = append(arr, []byte(key)...)
	arr = append(arr, '\r', '\n')
	fmt.Println(string(arr))
	return arr
}

func (r *RedisClient) BuildSet(arr []byte, key string, value string) []byte {
	arr = AppendSetCommand(arr)
	arr = append(arr, '$')
	arr = strconv.AppendInt(arr, int64(len(key)), 10)
	arr = append(arr, '\r', '\n')
	arr = append(arr, []byte(key)...)
	arr = append(arr, '\r', '\n')
	arr = append(arr, '$')
	v_len := len(value)
	arr = strconv.AppendInt(arr, int64(v_len), 10)
	arr = append(arr, '\r', '\n')
	arr = append(arr, []byte(value)...)
	arr = append(arr, '\r', '\n')
	fmt.Println(string(arr))
	return arr
}

func AppendSetCommand(arr []byte) []byte {
	arr = append(arr, '$')
	arr = strconv.AppendInt(arr, int64(3), 10)
	arr = append(arr, '\r', '\n')
	arr = append(arr, []byte("SET")...)
	return append(arr, '\r', '\n')
}
func AppendGetCommand(arr []byte) []byte {
	arr = append(arr, '$')
	arr = strconv.AppendInt(arr, int64(3), 10)
	arr = append(arr, '\r', '\n')
	arr = append(arr, []byte("GET")...)
	return append(arr, '\r', '\n')
}

func AppendBulkString(arr []byte)   {}
func AppendSimpleString(arr []byte) {}
func AppendError(arr []byte)        {}
func AppendNumber(arr []byte)       {}

func main() {
	c, err := NewClient("localhost", 6379)
	if err != nil {
		return
	}
	arr := c.Set("test", "value")
	buf := make([]byte, 512)
	c.conn.Write(arr)
	n, _ := c.conn.Read(buf)
	fmt.Println("response", string(buf[:n]))
}
