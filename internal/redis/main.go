package redis

import (
	"github.com/gomodule/redigo/redis"
	"log"
)

func EstablishConnectionToRedis(network string, address string) redis.Conn {
	conn, err := redis.Dial(network, address)
	if err != nil {
		log.Fatal(err)
	}

	return conn
}

func PushDataInList(conn redis.Conn, listName string, dataName string) {
	_, err := conn.Do("RPUSH", listName, dataName)
	if err != nil {
		log.Fatal(err)
	}
}
