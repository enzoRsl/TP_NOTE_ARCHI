package main

import (
	"fmt"
	"log"
	"time"

	mqtt "github.com/enzoRsl/TP_NOTE_ARCHI/internal/mqtt"
	opendatameteo "github.com/enzoRsl/TP_NOTE_ARCHI/internal/opendatameteo"
	utils "github.com/enzoRsl/TP_NOTE_ARCHI/internal/utils"
)

func init() {
	utils.InitEnvVariables()
}

func main() {
	publicationInterval := utils.GetPublicationIntervalInSeconds()
	sensorId, airport := utils.GetCliParams()

	// Initialise the MQTT client
	client := mqtt.GetMqttClient(sensorId)

	// Start a job to get the temperature every time given by the interval variable found in the .env file
	for {
		temperature := opendatameteo.GetTemperature(airport.Latitude, airport.Longitude, utils.GetCurrentDate())

		topic := fmt.Sprintf(
			"airport/%s/datatype/temperature/sensor/%s/date/%s/hour/%s/",
			airport.AirportIataCode, sensorId, utils.GetCurrentDate(), utils.GetCurrentHour(),
		)
		message := fmt.Sprintf("%s:%f", utils.GetCurrentDate(), temperature)

		client.Publish(topic, 0, false, message)
		log.Default().Println("Published message: ", message, " on topic: ", topic)

		time.Sleep(publicationInterval * time.Second)
	}
}
