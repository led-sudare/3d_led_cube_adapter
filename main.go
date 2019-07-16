package main

import (
	"cube_adapter/lib"
	"cube_adapter/lib/util"
	"flag"
	"time"

	log "github.com/cihub/seelog"
	zmq "github.com/zeromq/goczmq"
)

type Configs struct {
	XProxySubBind string `json:"XProxySubBind"`
}

func NewConfigs() Configs {
	return Configs{
		XProxySubBind: "0.0.0.0:5510",
	}
}

func main() {
	configs := NewConfigs()
	util.ReadConfig(&configs)

	var (
		xproxySubBind = flag.String("r", configs.XProxySubBind, "Specify IP and port of server zeromq SUB running.")
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

	endpoint := "tcp://" + *xproxySubBind
	zmqPubSock := zmq.NewSock(zmq.Pub)
	err = zmqPubSock.Connect(endpoint)
	if err != nil {
		panic(err)
	}
	log.Infof("Start ZMQ Pub %v", endpoint)
	defer zmqPubSock.Destroy()

	ticker := util.NewInlineTicker(2 * time.Second)

	for {
		buffer, _, err := zmqSubSock.RecvFrame()
		if err != nil {
			continue
		}
		converter := lib.NewLedCubeConverter(len(buffer))
		if converter == nil {
			ticker.DoIfFire(func() {
				log.Warn("Invalid datalength ", len(buffer))
			})
			continue
		}
		ticker.DoIfFire(func() {
			log.Infof("Received last data len: %v", len(buffer))
		})

		if sudareData := converter.ConvertToSudare(buffer); sudareData != nil {
			zmqPubSock.SendFrame(sudareData, zmq.FlagNone)
		}
	}
}
