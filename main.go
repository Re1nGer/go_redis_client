package redisClient

import (
	"fmt"
	"net"
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
	return make([]byte, 0)
}

func (r *RedisClient) Set(key string, val string) []byte {
	return make([]byte, 0)
}

func (r *RedisClient) BuildArray(count int) []byte {
	return make([]byte, 0)
}

func main() {
	c := NewClient("localhost", 6379)
	fmt.Println(c)
}
