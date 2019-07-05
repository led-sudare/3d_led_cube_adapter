package main

import (
	"3d_led_cube_adapter/lib"
	"3d_led_cube_adapter/lib/util"
	"flag"

	log "github.com/cihub/seelog"
	zmq "github.com/zeromq/goczmq"
)

type Configs struct {
	ZmqTarget string `json:"zmqTarget"`
}

func NewConfigs() Configs {
	return Configs{
		ZmqTarget: "0.0.0.0:5510",
	}
}

func main() {
	configs := NewConfigs()
	util.ReadConfig(&configs)

	var (
		optInputPort = flag.String("r", configs.ZmqTarget, "Specify IP and port of server zeromq SUB running.")
	)

	flag.Parse()

	zmqSubSock := zmq.NewSock(zmq.Sub)
	zmqSubSock.SetSubscribe("")

	_, err := zmqSubSock.Bind("tcp://*:5520")
	if err != nil {
		panic(err)
	}
	log.Infof("Start ZMQ Sub tcp://*:5520")
	defer zmqSubSock.Destroy()

	endpoint := "tcp://" + *optInputPort
	zmqPubSock := zmq.NewSock(zmq.Pub)
	err = zmqPubSock.Connect(endpoint)
	if err != nil {
		panic(err)
	}
	log.Infof("Start ZMQ Pub %v", endpoint)
	defer zmqPubSock.Destroy()

	for {
		buffer, _, err := zmqSubSock.RecvFrame()
		if err != nil {
			continue
		}
		converter := lib.NewLedCubeConverter(len(buffer))
		log.Infof("Received : %v\n", len(buffer))

		if sudareData := converter.ConvertToSudare(buffer); sudareData != nil {
			zmqPubSock.SendFrame(sudareData, zmq.FlagNone)
		}
	}
}
