package main

import (
	"3d_led_cube_adapter/lib/util"
	"flag"
	"fmt"
	"math"
	"net"

	log "github.com/cihub/seelog"
	zmq "github.com/zeromq/goczmq"
)

const LedCubeWidth = 16
const LedCubeHeight = 32
const LedCubeDepth = 8

const LedCylinderDiameter = 30
const LedCylinderRadius = LedCylinderDiameter / 2
const LedCylinderHeight = 100
const LedCylinderCount = 60

func rgb565to888(c565 uint32) []byte {
	r := (byte)((c565 & 0xF800) >> 8)
	g := (byte)((c565 & 0x7E0) >> 3)
	b := (byte)((c565 & 0x1F) << 3)
	return []byte{r, g, b}
}

func makeSudare(cube []byte) []byte {
	sudareBuf := make([]byte, LedCylinderRadius*LedCylinderHeight*LedCylinderCount*2)

	offsetCubeX := LedCubeWidth / 2.0
	offsetCubeZ := LedCubeDepth / 2.0

	cylinderUnitDegree := (2 * math.Pi) / LedCylinderCount

	xscale := float64(LedCubeWidth) / float64(LedCylinderDiameter)
	yscale := float64(LedCubeHeight) / float64(LedCylinderHeight)

	util.ConcurrentEnum(0, LedCylinderCount, func(cylinder int) {
		sin := math.Sin(cylinderUnitDegree * float64(cylinder))
		cos := math.Cos(cylinderUnitDegree * float64(cylinder))

		for r := 0; r < LedCylinderRadius; r++ {
			cubex := int(math.Round(offsetCubeX + cos*float64(r)*xscale))
			cubez := int(math.Round(offsetCubeZ + sin*float64(r)*xscale))
			for y := 0; y < LedCylinderHeight; y++ {
				cubey := (LedCubeHeight - 1) - int(math.Round(float64(y)*yscale))

				if cubex >= 0 && cubey >= 0 && cubez >= 0 &&
					cubex < LedCubeWidth && cubez < LedCubeDepth && cubey < LedCubeHeight {
					idxS := ((LedCylinderHeight * LedCylinderRadius * cylinder) +
						(LedCylinderRadius * y) +
						r) * 2

					idxC := (cubez + cubey*LedCubeDepth + cubex*LedCubeHeight*LedCubeDepth) * 2
					sudareBuf[idxS] = cube[idxC]
					sudareBuf[idxS+1] = cube[idxC+1]
				}
			}
		}
	})
	return sudareBuf
}

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

	fmt.Println("Server is Running at 0.0.0.0:9001")
	conn, _ := net.ListenPacket("udp", "0.0.0.0:9001")
	defer conn.Close()

	endpoint := "tcp://" + *optInputPort
	zmqsock := zmq.NewSock(zmq.Pub)
	err := zmqsock.Connect(endpoint)
	if err != nil {
		panic(err)
	}
	defer zmqsock.Destroy()

	buffer := make([]byte, 8192)
	for {
		length, remoteAddr, _ := conn.ReadFrom(buffer)
		conv := makeSudare(buffer)
		zmqsock.SendFrame(conv, zmq.FlagNone)
		log.Infof("Received from %v: %v\n", remoteAddr, length)
	}

}
