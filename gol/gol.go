package gol

// Params provides the details of how to run the Game of Life and which image to load.
type Params struct {
	Turns       int
	Threads     int
	ImageWidth  int
	ImageHeight int
}

// Run starts the processing of Game of Life. It should initialise channels and goroutines.
func Run(p Params, events chan<- Event, keyPresses <-chan rune) {
	<-keyPresses
	//	TODO: Put the missing channels in here.

	ioCommand := make(chan ioCommand)
	ioIdle := make(chan bool)
	ioFilename := make(chan string)
	ioOutput := make(chan uint8)
	ioInput := make(chan uint8)

	ioChannels := ioChannels{
		command:  ioCommand,
		idle:     ioIdle,
		filename: ioFilename, // ioCommand 에서 채널이 만들어 졌긴 하지만 값이 nil -> 수정 (o)
		output:   ioOutput,
		input:    ioInput,
	}
	go startIo(p, ioChannels) // 1. IO 고루틴 이 시작 되는곳!, input argument p == Params 값들, 그리고 채널이 인풋 값.

	distributorChannels := distributorChannels{
		events:     events,
		ioCommand:  ioCommand,
		ioIdle:     ioIdle,
		ioFilename: ioFilename, // 맞지?
		ioOutput:   ioOutput,   //?
		ioInput:    ioInput,    // ?
	}
	distributor(p, distributorChannels)

}
