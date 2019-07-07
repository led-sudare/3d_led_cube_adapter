package lib

import (
	"3d_led_cube_adapter/lib/util"
	"math"

	log "github.com/cihub/seelog"
)

type CubeConverter interface {
	GetCubeWidth() int
	GetCubeHeight() int
	GetCubeDepth() int
	ConvertToSudare(cube []byte) []byte
}

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

type ledCubeConverter struct {
	width, height, depth, datalen int
}

func NewLedCubeConverter(dataLength int) CubeConverter {
	switch dataLength {
	case 8192:
		return &ledCubeConverter{16, 32, 8, 8192}
	case 22500:
		return &ledCubeConverter{15, 50, 15, 22500}
	default:
		return nil
	}
}

func (c *ledCubeConverter) GetCubeWidth() int {
	return c.width
}

func (c *ledCubeConverter) GetCubeHeight() int {
	return c.height

}
func (c *ledCubeConverter) GetCubeDepth() int {
	return c.depth
}

func (c *ledCubeConverter) ConvertToSudare(cube []byte) []byte {

	if len(cube) != c.datalen {
		log.Warn("Invalid datalength ", len(cube))
		return nil
	}

	sudareBuf := make([]byte, LedCylinderRadius*LedCylinderHeight*LedCylinderCount*2)

	offsetCubeX := float64(c.width) / 2.0
	offsetCubeZ := float64(c.depth) / 2.0

	cylinderUnitDegree := (2 * math.Pi) / LedCylinderCount

	xscale := float64(c.width) / float64(LedCylinderDiameter)
	yscale := float64(c.height) / float64(LedCylinderHeight)

	util.ConcurrentEnum(0, LedCylinderCount, func(cylinder int) {
		sin := math.Sin(cylinderUnitDegree * float64(cylinder))
		cos := math.Cos(cylinderUnitDegree * float64(cylinder))

		for r := 0; r < LedCylinderRadius; r++ {
			cubex := int(math.Round(offsetCubeX + cos*float64(r)*xscale))
			cubez := int(math.Round(offsetCubeZ + sin*float64(r)*xscale))
			for y := 0; y < LedCylinderHeight; y++ {
				cubey := (c.height - 1) - int(math.Round(float64(y)*yscale))

				if cubex >= 0 && cubey >= 0 && cubez >= 0 &&
					cubex < c.width && cubez < c.depth && cubey < c.height {
					idxS := ((LedCylinderHeight * LedCylinderRadius * cylinder) +
						(LedCylinderRadius * y) +
						r) * 2

					idxC := (cubez + cubey*c.depth + cubex*c.height*c.depth) * 2
					sudareBuf[idxS] = cube[idxC]
					sudareBuf[idxS+1] = cube[idxC+1]
				}
			}
		}
	})
	return sudareBuf
}
