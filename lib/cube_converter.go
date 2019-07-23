package lib

import (
	"cube_adapter/lib/util"
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

const rgb565plane = 2

type cubeConverterImpl struct {
	Width, Height, Depth, DataLen int
}

func newLedCubeConverter(width, height, depth int) CubeConverter {
	return &cubeConverterImpl{
		Width:   width,
		Height:  height,
		Depth:   depth,
		DataLen: width * height * depth * rgb565plane}
}

var ledConverter16x32x8 CubeConverter
var ledConverter15x50x15 CubeConverter
var ledConverter30x100x30 CubeConverter

func NewLedCubeConverter(dataLength int) CubeConverter {
	switch dataLength {
	case 16 * 32 * 8 * rgb565plane:
		if ledConverter16x32x8 == nil {
			ledConverter16x32x8 = newLedCubeConverter(16, 32, 8) // data len 8,192
		}
		return ledConverter16x32x8
	case 15 * 50 * 15 * rgb565plane:
		if ledConverter15x50x15 == nil {
			ledConverter15x50x15 = newLedCubeConverter(15, 50, 15) // data len 22,500
		}
		return ledConverter15x50x15
	case 30 * 100 * 30 * rgb565plane:
		if ledConverter30x100x30 == nil {
			ledConverter30x100x30 = newLedCubeConverter(30, 100, 30) // data len 180,000
		}
		return ledConverter30x100x30
	default:
		return nil
	}
}

func (c *cubeConverterImpl) GetCubeWidth() int {
	return c.Width
}

func (c *cubeConverterImpl) GetCubeHeight() int {
	return c.Height

}
func (c *cubeConverterImpl) GetCubeDepth() int {
	return c.Depth
}

func (c *cubeConverterImpl) ConvertToSudare(cube []byte) []byte {

	if len(cube) != c.DataLen {
		log.Warn("Invalid datalength ", len(cube))
		return nil
	}

	sudareBuf := make([]byte, LedCylinderRadius*LedCylinderHeight*LedCylinderCount*2)

	offsetCubeX := float64(c.Width) / 2.0
	offsetCubeZ := float64(c.Depth) / 2.0

	cylinderUnitDegree := (2 * math.Pi) / LedCylinderCount

	xscale := float64(c.Width) / float64(LedCylinderDiameter)
	yscale := float64(c.Height) / float64(LedCylinderHeight)

	util.ConcurrentEnum(0, LedCylinderCount, func(cylinder int) {
		sin := math.Sin(cylinderUnitDegree * float64(cylinder))
		cos := math.Cos(cylinderUnitDegree * float64(cylinder))

		for r := 0; r < LedCylinderRadius; r++ {
			cubex := int(math.Round(offsetCubeX + cos*float64(r)*xscale))
			cubez := int(math.Round(offsetCubeZ + sin*float64(r)*xscale))
			for y := 0; y < LedCylinderHeight; y++ {
				cubey := (c.Height - 1) - int(math.Round(float64(y)*yscale))

				if cubex >= 0 && cubey >= 0 && cubez >= 0 &&
					cubex < c.Width && cubez < c.Depth && cubey < c.Height {
					idxS := ((LedCylinderRadius*cylinder+r)*LedCylinderHeight + y) * 2

					idxC := (cubez + cubey*c.Depth + cubex*c.Height*c.Depth) * 2
					sudareBuf[idxS] = cube[idxC]
					sudareBuf[idxS+1] = cube[idxC+1]
				}
			}
		}
	})
	return sudareBuf
}
