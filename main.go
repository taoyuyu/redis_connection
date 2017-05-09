package main

import (
	"fmt"
	"runtime"
	"strconv"
	"sync"

	"redis_connection/redis_connection_pool"
)

func main() {
	runtime.GOMAXPROCS(4)
	err := redis_connection_pool.SetSize(8)
	if err != nil {
		fmt.Println(err)
		return
	}

	err = redis_connection_pool.InitConnection("127.0.0.1", 6379)
	if err != nil {
		fmt.Println(err)
		return
	}

	wc := new(sync.WaitGroup)
	wc.Add(2)
	go func() {
		connect := redis_connection_pool.GetConnection()
		defer wc.Done()
		defer redis_connection_pool.ReturnConnection(connect)
		for i := 0; i < 100; i++ {
			fmt.Println("func 1:" + strconv.Itoa(i))
			_, err1 := (*connect).Do("set", "taoyu"+strconv.Itoa(i), strconv.Itoa(i))
			if err1 != nil {
				fmt.Println(err1)
				return
			}
		}

	}()
	go func() {
		connect := redis_connection_pool.GetConnection()
		defer wc.Done()
		defer redis_connection_pool.ReturnConnection(connect)
		for i := 100; i < 200; i++ {
			fmt.Println("func 2:" + strconv.Itoa(i))
			_, err1 := (*connect).Do("set", "taoyu"+strconv.Itoa(i), strconv.Itoa(i))
			if err1 != nil {
				fmt.Println(err1)
				return
			}
		}
	}()

	wc.Wait()
	fmt.Println("wait done")

	redis_connection_pool.CloseConnection()
	return
}
