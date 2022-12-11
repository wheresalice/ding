package main

import (
	"fmt"
	paho "github.com/eclipse/paho.mqtt.golang"
	"github.com/urfave/cli/v2"
	"github.com/wheresalice/ding/lib"
	mqtt "github.com/wind-c/comqtt/server"
	"github.com/wind-c/comqtt/server/events"
	"github.com/wind-c/comqtt/server/listeners"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

var username string

func server() {
	sigs := make(chan os.Signal, 1)
	done := make(chan bool, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigs
		done <- true
	}()

	server := mqtt.NewServer(nil)
	tcp := listeners.NewTCP("t1", ":1883")
	err := server.AddListener(tcp, nil)
	if err != nil {
		log.Fatalln(err)
	}

	go func() {
		err := server.Serve()
		if err != nil {
			log.Fatalln(err)
		}
		log.Println("Listening on port 1883")
	}()

	server.Events.OnProcessMessage = func(client events.Client, packet events.Packet) (events.Packet, error) {
		if string(packet.Payload) != "ding" {
			log.Println("bad message")
			return events.Packet{}, mqtt.ErrRejectPacket
		} else {
			log.Printf("ding for %s", packet.TopicName)
			return packet, nil
		}
	}
	<-done
}

func listen() {
	if len(os.Args) != 3 {
		fmt.Println("Usage: ding listen username@servername")
		os.Exit(1)
	}
	input := strings.Split(os.Args[2], "@")
	username = input[0]
	server := input[1]

	sigs := make(chan os.Signal, 1)
	done := make(chan bool, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigs
		done <- true
	}()

	opts := paho.NewClientOptions()
	opts.AddBroker(fmt.Sprintf("tcp://%s:%d", server, 1883))
	opts.SetClientID(username)
	opts.SetDefaultPublishHandler(lib.MessagePubHandler)
	opts.OnConnect = lib.ConnectHandler
	opts.OnConnectionLost = lib.ConnectLostHandler
	client := paho.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	<-done
}

func ding() {
	if len(os.Args) != 3 {
		fmt.Println("Usage: ding listen username@servername")
		os.Exit(1)
	}
	input := strings.Split(os.Args[2], "@")
	username := input[0]
	server := input[1]

	opts := paho.NewClientOptions()
	opts.AddBroker(fmt.Sprintf("tcp://%s:%d", server, 1883))
	client := paho.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}
	token := client.Publish(username, 0, false, "ding")
	token.Wait()
}

func main() {
	app := &cli.App{
		Name:  "ding",
		Usage: `send and receive dings`,
		Commands: []*cli.Command{
			{
				Name:    "server",
				Aliases: []string{"s"},
				Usage:   "run a ding server",
				Action: func(cCtx *cli.Context) error {
					server()
					return nil
				},
			},
			{
				Name:    "listen",
				Aliases: []string{"l"},
				Usage:   "listen for dings",
				Action: func(cCtx *cli.Context) error {
					listen()
					return nil
				},
			},
			{
				Name:    "ding",
				Aliases: []string{"d"},
				Usage:   "send a ding",
				Action: func(cCtx *cli.Context) error {
					ding()
					return nil
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
