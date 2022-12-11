package lib

import (
	"fmt"
	paho "github.com/eclipse/paho.mqtt.golang"
	"log"
)

var MessagePubHandler paho.MessageHandler = func(client paho.Client, msg paho.Message) {
	log.Printf("Received %s for %s\n", msg.Payload(), msg.Topic())
	fmt.Printf("\a")
}

var ConnectHandler paho.OnConnectHandler = func(client paho.Client) {
	log.Println("Connected")
	opts := client.OptionsReader()
	topic := opts.ClientID()
	token := client.Subscribe(topic, 1, nil)
	token.Wait()
	log.Printf("Subscribed as %v\n", topic)
}

var ConnectLostHandler paho.ConnectionLostHandler = func(client paho.Client, err error) {
	log.Printf("Connect lost: %v\n", err)
}
