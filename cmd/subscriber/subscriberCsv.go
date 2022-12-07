package main

import (
	m "github.com/eclipse/paho.mqtt.golang"
	"github.com/enzoRsl/TP_NOTE_ARCHI/internal/csv"
	"github.com/enzoRsl/TP_NOTE_ARCHI/internal/mqtt"
	"strings"
	"sync"
)

func main() {
	client := mqtt.GetMqttClient("subscriberCsv")
	waitGroup := sync.WaitGroup{}
	waitGroup.Add(1)
	_, csvWriter := csv.CreateCsvWriter("database.csv")

	token := client.Subscribe("airport/#", 2, func(client m.Client, message m.Message) {
		s := strings.Split(message.Topic(), "/")
		var value []string
		for i := 1; i < len(s); i += 2 {
			value = append(value, s[i])
		}
		value = append(value, string(message.Payload()))
		csv.WriteInCsv(csvWriter, value)
		defer csvWriter.Flush()
	})
	token.Wait()
	if err := token.Error(); err != nil {
		panic(err)
	}
	waitGroup.Wait()
}
