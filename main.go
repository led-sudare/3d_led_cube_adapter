package main

import (
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

const LedCylinderWidth = 30
const LedCylinderHeight = 100
const LedCylinderCount = 60

func rgb565to888(c565 uint32) []byte {
	r := (byte)((c565 & 0xF800) >> 8)
	g := (byte)((c565 & 0x7E0) >> 3)
	b := (byte)((c565 & 0x1F) << 3)
	return []byte{r, g, b}
}

func makeSudare(cube []byte) []byte {
	sudareBuf := make([]byte, 360000)

	offsetCubeX := LedCubeWidth / 2.0
	offsetCubeZ := LedCubeDepth / 2.0

	cylinderUnitDegree := math.Pi / LedCylinderCount

	xscale := float64(LedCubeWidth) / float64(LedCylinderWidth)
	yscale := float64(LedCubeHeight) / float64(LedCylinderHeight)

	for cylinder := 0; cylinder < LedCylinderCount; cylinder++ {
		sin := math.Sin(cylinderUnitDegree * float64(cylinder))
		cos := math.Cos(cylinderUnitDegree * float64(cylinder))

		for r := 0; r < LedCylinderWidth; r++ {
			cubex := int(math.Round(offsetCubeX + cos*float64(r-(LedCylinderWidth/2))*xscale))
			cubez := int(math.Round(offsetCubeZ + sin*float64(r-(LedCylinderWidth/2))*xscale))
			for y := LedCylinderHeight - 1; 0 <= y; y-- {
				cubey := int(math.Round(float64(y) * yscale))

				if cubex >= 0 && cubey >= 0 && cubez >= 0 &&
					cubex < LedCubeWidth && cubez < LedCubeDepth && cubey < LedCubeHeight {
					idxS := ((LedCylinderHeight * LedCylinderWidth * cylinder) +
						(LedCylinderWidth * y) +
						r) * 2

					idxC := (cubez + cubey*LedCubeDepth + cubex*LedCubeHeight*LedCubeDepth) * 2
					sudareBuf[idxS] = cube[idxC]
					sudareBuf[idxS+1] = cube[idxC+1]

				}
			}
		}
	}
	return sudareBuf
}

var (
	port         = flag.String("p", ":2345", "http service port")
	logVerbose   = flag.Bool("v", false, "output detailed log.")
	optInputPort = flag.String("r", "127.0.0.1:5563", "Specify IP and port of server main_realsense_serivce.py running.")
)

func main() {
	flag.Parse()

	fmt.Println("Server is Running at localhost:9001")
	conn, _ := net.ListenPacket("udp", "localhost:9001")
	defer conn.Close()

	endpoint := "tcp://" + *optInputPort
	zmqsock, err := zmq.NewPub(endpoint)
	if err != nil {
		panic(err)
	}
	err = zmqsock.Connect(endpoint)
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
