package redisClient

type RedisClient struct {
	host     string
	port     int
	username string
	password string
	ssl      bool
	//for now it suffices to have just these fields
}

func NewClient(host string, port int) *RedisClient {
	return &RedisClient{
		host: host,
		port: port,
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
