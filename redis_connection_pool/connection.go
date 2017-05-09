package redis_connection_pool

import (
	"errors"
	"log"
	"strconv"

	"github.com/garyburd/redigo/redis"
)

var max_size = 10
var default_size = 5
var redisPool chan *redis.Conn

func SetSize(size int) (err error) {
	if size <= 0 {
		err = errors.New("size <= 0 error")
		return
	}
	if size <= max_size {
		default_size = size
	} else {
		default_size = max_size
	}
	return
}

func InitConnection(host string, port int) error {
	if redisPool != nil {
		err := errors.New("Connection already initialized")
		return err
	}
	redisPool = make(chan *redis.Conn, default_size)

	hostName := host + ":" + strconv.Itoa(port)

	for i := 0; i < default_size; i++ {
		rs, err := redis.Dial("tcp", hostName)
		if err != nil {
			log.Println("INFO: Redis connect error: ", err)
			return err
		}
		redisPool <- &rs
	}
	return nil
}

func GetConnection() *redis.Conn {
	return <-redisPool
}

func ReturnConnection(rc *redis.Conn) error {
	if rc == nil {
		return errors.New("Can't return nil connection")
	}
	redisPool <- rc
	return nil
}

func CloseConnection() error {
	//close(channel)后只能读不能写,range完成直接退出
	close(redisPool)
	for rs := range redisPool {
		err := (*rs).Close()
		if err != nil {
			return err
		}
		log.Println("one connection closed succeed")
	}
	return nil
}
