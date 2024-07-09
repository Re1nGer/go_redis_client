package main

import (
	"fmt"
	"net"
	"strconv"
)

type RedisClient struct {
	conn     net.Conn
	host     string
	port     int
	username string
	password string
	ssl      bool
	//for now it suffices to have just these fields
}

func NewClient(host string, port int) *RedisClient {
	//for now let's just ignore potential error
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", host, port))
	if err != nil {
		fmt.Errorf("%s", err.Error())
		return nil
	}
	return &RedisClient{
		host: host,
		port: port,
		conn: conn,
	}
}

func (r *RedisClient) Get(key string) []byte {
	arr := r.BuildArray(2)
	arr = r.BuildGet(arr, key)
	_, err := r.conn.Write(arr)

	if err != nil {
		fmt.Println(err)
	}
	return arr
}

func (r *RedisClient) Set(key string, val string) []byte {
	return make([]byte, 0)
}

func (r *RedisClient) BuildArray(count int64) []byte {
	arr := make([]byte, 0)
	arr = append(arr, '*')
	arr = strconv.AppendInt(arr, count, 10)
	return append(arr, '\r', '\n')
}

func (r *RedisClient) BuildGet(arr []byte, key string) []byte {
	arr = append(arr, []byte("3$")...)
	arr = append(arr, '\r', '\n')
	l := len(key)
	arr = append(arr, '$')
	arr = strconv.AppendInt(arr, int64(l), 10)
	arr = append(arr, '\r', '\n')
	arr = append(arr, []byte(key)...)
	arr = append(arr, '\r', '\n')
	fmt.Println(string(arr))
	return arr
}

func main() {
	c := NewClient("localhost", 6379)
	arr := c.Get("test")
	fmt.Println(c, string(arr))
	buf := make([]byte, 512)
	n, _ := c.conn.Read(buf)
	fmt.Println("response", string(buf[:n]))
}
