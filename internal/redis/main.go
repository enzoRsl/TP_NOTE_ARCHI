package redis

import (
	"log"

	"github.com/gomodule/redigo/redis"
)

func EstablishConnectionToRedis(network string, address string) redis.Conn {
	conn, err := redis.Dial(network, address)
	if err != nil {
		log.Fatal(err)
	}

	return conn
}

func dataAlreadyInList(conn redis.Conn, listName string, data string) bool {
	datas, err := redis.Strings(conn.Do("LRANGE", listName, 0, -1))
	if err != nil {
		log.Fatal(err)
	}
	for _, element := range datas {
		if element == data {
			return true
		}
	}
	return false
}

func PushDataInList(conn redis.Conn, listName string, data string, pushIfDataNotInList bool) {
	if pushIfDataNotInList && dataAlreadyInList(conn, listName, data) {
		return
	}
	_, err := conn.Do("RPUSH", listName, data)
	if err != nil {
		log.Fatal(err)
	}
}
