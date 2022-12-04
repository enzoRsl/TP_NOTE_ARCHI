package mqtt

import (
	"fmt"
	"log"
	"os"
	"time"

	utils "github.com/enzoRsl/TP_NOTE_ARCHI/internal/utils"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

func init() {
	utils.InitEnvVariables()
}

// CreateClientOptions creates the client options for the MQTT client
func createClientOptions(brokerURI string, clientId string, username string, password string) *mqtt.ClientOptions {
	opts := mqtt.NewClientOptions()

	// AddBroker adds a broker URI to the list of brokers to be used.
	// The format should be "scheme://host:port"
	opts.AddBroker(brokerURI)

	if username != "" && password != "" {
		opts.SetUsername(username)
		opts.SetPassword(password)
	}
	opts.SetClientID(clientId)
	return opts
}

// Connect connects to the MQTT broker
func connect(brokerURI string, clientId string) mqtt.Client {
	fmt.Printf("Trying to connect to %s with clientId %s\n", brokerURI, clientId)

	username := os.Getenv("MQTT_USERNAME")
	password := os.Getenv("MQTT_PASSWORD")

	opts := createClientOptions(brokerURI, clientId, username, password)
	client := mqtt.NewClient(opts)
	token := client.Connect()

	for !token.WaitTimeout(3 * time.Second) {
	}
	if err := token.Error(); err != nil {
		log.Fatal(err)
	}

	return client
}

// Generates the MQTT client with the .env variables
func GetMqttClient(sensorName string, sensorId string) mqtt.Client {
	client := connect(os.Getenv("MQTT_URI"), sensorName+sensorId)
	return client
}
