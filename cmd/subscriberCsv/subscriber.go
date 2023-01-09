package main

import (
	"os"
	"strings"
	"sync"

	m "github.com/eclipse/paho.mqtt.golang"
	"github.com/enzoRsl/TP_NOTE_ARCHI/internal/csv"
	"github.com/enzoRsl/TP_NOTE_ARCHI/internal/mqtt"
)

func main() {
	client := mqtt.GetMqttClient("subscriberCsv")
	waitGroup := sync.WaitGroup{}
	waitGroup.Add(1)
	qos := byte(os.Getenv("MQTT_QOS")[0] - '0')
	_, csvWriter := csv.CreateCsvWriter("database.csv")
	headers := []string{"airport", "datatype", "sensor", "date", "hour", "datetime", "value"}
	csv.WriteInCsv(csvWriter, headers)
	token := client.Subscribe("airport/#", qos, func(client m.Client, message m.Message) {
		s := strings.Split(message.Topic(), "/")
		var value []string
		for i := 1; i < len(s); i += 2 {
			value = append(value, s[i])
		}
		messagePayloadToArrayOfString := strings.Split(string(message.Payload()), ":")
		for _, element := range messagePayloadToArrayOfString {
			value = append(value, element)
		}
		csv.WriteInCsv(csvWriter, value)
		defer csvWriter.Flush()
	})
	token.Wait()
	if err := token.Error(); err != nil {
		panic(err)
	}
	waitGroup.Wait()
}
