package gol

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

	turn := 0
	// TODO: Execute all turns of the Game of Life.

	// TODO: Report the final state using FinalTurnCompleteEvent.

	// Make sure that the Io has finished any output before exiting.
	c.ioCommand <- ioCheckIdle
	<-c.ioIdle

	c.events <- StateChange{turn, Quitting}

	// Close the channel to stop the SDL goroutine gracefully. Removing may cause deadlock.
	close(c.events)
}
