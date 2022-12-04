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
	client := mqtt.GetMqttClient("pressure", sensorId)

	// Start a job to get the pressure every time given by the interval variable found in the .env file
	for {
		pressure := opendatameteo.GetAirPressure(airport.Latitude, airport.Longitude, utils.GetCurrentDate())

		topic := fmt.Sprintf(
			"airport/%s/datatype/pressure/sensor/%s/date/%s/hour/%s/",
			airport.AirportIataCode, sensorId, utils.GetCurrentDate(), utils.GetCurrentHour(),
		)
		message := fmt.Sprintf(
			"%s:%f",
			utils.GetCurrentDateAndHour(), pressure,
		)

		client.Publish(topic, 0, false, message)
		log.Default().Println("Published message: ", message, " on topic: ", topic)

		time.Sleep(publicationInterval * time.Second)
	}
}
