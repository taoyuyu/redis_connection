package redis_connection_pool

import (
	"errors"
	"log"
	"strconv"

	"github.com/garyburd/redigo/redis"
)

var MAX_SIZE = 10
var DEFAULT_SIZE = 5
var redisPool chan *redis.Conn

func SetSize(size int) (err error) {
	if size <= 0 {
		err = errors.New("size <= 0 error")
		return
	}
	if size <= MAX_SIZE {
		DEFAULT_SIZE = size
	} else {
		DEFAULT_SIZE = MAX_SIZE
	}
	return
}

func InitConnection(host string, port int) error {
	if redisPool != nil {
		err := errors.New("Connection already initialized")
		return err
	}
	redisPool = make(chan *redis.Conn, DEFAULT_SIZE)

	hostName := host + ":" + strconv.Itoa(port)

	for i := 0; i < DEFAULT_SIZE; i++ {
		rs, err := redis.Dial("tcp", hostName)
		if err != nil {
			log.Println("INFO: Redis connect error: ", err)
			return err
		}
		// select 0 in redis
		rs.Do("SELECT", 0)
		redisPool <- &rs
	}
	return nil
}

func GetConnection() *redis.Conn {
	return <-redisPool
}

func ReturnConnection(rc *redis.Conn) {
	redisPool <- rc
	return
}

func CloseConnection() {
	for rs := range redisPool {
		(*rs).Close()
	}
	close(redisPool)
}
