package gol

import (
	"uk.ac.bris.cs/gameoflife/util"
)

type distributorChannels struct {
	events     chan<- Event
	ioCommand  chan<- ioCommand
	ioIdle     <-chan bool
	ioFilename chan<- string
	ioOutput   chan<- uint8
	ioInput    <-chan uint8
} //이중 어떤 애들은 사용하기전에 wired up 해야함

// distributor divides the work between workers and interacts with other goroutines.
func distributor(p Params, c distributorChannels) {
	// TODO: Create a 2D slice to store the world.

	// 파일 이름 선정 by p,만약 256X256 파일이 들어온다면, 스트링은 만들 수 있고, 알맞는 채널을 통해 (위에 있는 채널들 중) 보낸다
	//그리고 이미지를 바이트바이바이트로 받아서(get)해서 이 2D World 에 Store!

	if p.ImageWidth == 16 && p.ImageHeight == 16 {
		c.ioFilename <- "16x16"
	} else if p.ImageWidth == 128 && p.ImageHeight == 128 {
		c.ioFilename <- "128x128"
	} else if p.ImageWidth == 256 && p.ImageHeight == 256 {
		c.ioFilename <- "256x256"
	} else if p.ImageWidth == 512 && p.ImageHeight == 512 {
		c.ioFilename <- "512x512"
	}

	// TODO: Execute all turns of the Game of Life. // 포문 이용, 두개의 2D 슬라이스 필요
	turn := 0

	world := make([][]byte, p.ImageHeight)
	for i := range world {
		world[i] = make([]byte, p.ImageWidth)
	}
	var cells []util.Cell

	for range world {
		for y := 0; y < p.ImageHeight; y++ {
			for x := 0; x < p.ImageWidth; x++ {
				val := <-c.ioInput // 중요, 그리고 작성하고 싶던 그 이미지를 바이트바이바이트로 이 채널로 보냄, 그리고 배열에 put,
				world[y][x] = val

				cells = append(cells, util.Cell{
					X: y,
					Y: x,
				})
			}
		}
		turn++
	}

	// TODO: Report the final state using FinalTurnCompleteEvent.

	c.events <- FinalTurnComplete{turn, cells}

	// Make sure that the Io has finished any output before exiting.
	c.ioCommand <- ioCheckIdle
	<-c.ioIdle

	c.events <- StateChange{turn, Quitting}

	// Close the channel to stop the SDL goroutine gracefully. Removing may cause deadlock.
	close(c.events)
}
