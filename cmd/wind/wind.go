package main

import (
	"fmt"
	"log"
	"os"
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
	client := mqtt.GetMqttClient("windspeed" + sensorId)

	// Start a job to get the wind speed every time given by the interval variable found in the .env file
	for {
		windSpeed := opendatameteo.GetWindSpeed(airport.Latitude, airport.Longitude, utils.GetCurrentDate())

		topic := fmt.Sprintf(
			"airport/%s/datatype/windspeed/date/%s/hour/%s/",
			airport.AirportIataCode, utils.GetCurrentDate(), utils.GetCurrentHour(),
		)
		message := fmt.Sprintf(
			"%s:%s:%f",
			sensorId, utils.GetCurrentDateAndHour(), windSpeed,
		)

		qos := byte(os.Getenv("MQTT_QOS")[0] - '0')
		client.Publish(topic, qos, false, message)
		log.Default().Println("Published message: ", message, " on topic: ", topic)

		time.Sleep(publicationInterval * time.Second)
	}
}
