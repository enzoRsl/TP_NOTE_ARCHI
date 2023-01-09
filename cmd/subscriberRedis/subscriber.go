package main

import (
	"os"
	"strings"
	"sync"

	m "github.com/eclipse/paho.mqtt.golang"
	"github.com/enzoRsl/TP_NOTE_ARCHI/internal/mqtt"
	"github.com/enzoRsl/TP_NOTE_ARCHI/internal/redis"
	"github.com/enzoRsl/TP_NOTE_ARCHI/internal/utils"
)

func init() {
	utils.InitEnvVariables()
}

func main() {
	client := mqtt.GetMqttClient("subscriberRedis")
	waitGroup := sync.WaitGroup{}
	waitGroup.Add(1)
	qos := byte(os.Getenv("MQTT_QOS")[0] - '0')
	token := client.Subscribe("airport/#", qos, func(client m.Client, message m.Message) {
		conn := redis.EstablishConnectionToRedis("tcp", os.Getenv("REDIS_URI"))
		topicSeparated := strings.Split(message.Topic(), "/")
		airport := topicSeparated[1]
		datatype := topicSeparated[3]
		redis.PushDataInList(conn, "airport", airport, true)
		redis.PushDataInList(conn, "datatype", datatype, true)
		redis.PushDataInList(conn, message.Topic(), string(message.Payload()), false)
	})
	token.Wait()
	if err := token.Error(); err != nil {
		panic(err)
	}
	waitGroup.Wait()

}
