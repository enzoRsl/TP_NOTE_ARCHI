package main

import (
	m "github.com/eclipse/paho.mqtt.golang"
	"github.com/enzoRsl/TP_NOTE_ARCHI/internal/mqtt"
	"github.com/enzoRsl/TP_NOTE_ARCHI/internal/redis"
	"os"
	"sync"
)

func main() {
	client := mqtt.GetMqttClient("subscriberRedis")
	waitGroup := sync.WaitGroup{}
	waitGroup.Add(1)
	token := client.Subscribe("airport/#", 2, func(client m.Client, message m.Message) {
		conn := redis.EstablishConnectionToRedis("tcp", os.Getenv("REDIS_URI"))
		redis.PushDataInList(conn, message.Topic(), string(message.Payload()))
	})
	token.Wait()
	if err := token.Error(); err != nil {
		panic(err)
	}
	waitGroup.Wait()

}
