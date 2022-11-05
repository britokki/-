package gol

import (
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"uk.ac.bris.cs/gameoflife/util"
)

type ioChannels struct {
	command  <-chan ioCommand
	idle     chan<- bool
	filename <-chan string
	output   <-chan uint8
	input    chan<- uint8
}

// ioState is the internal ioState of the io goroutine.
type ioState struct {
	params   Params
	channels ioChannels
}

// ioCommand allows requesting behaviour from the io (pgm) goroutine.
type ioCommand uint8

// This is a way of creating enums in Go.
// It will evaluate to:
//
//	ioOutput 	= 0
//	ioInput 	= 1
//	ioCheckIdle = 2
const (
	ioOutput ioCommand = iota
	ioInput
	ioCheckIdle
)

// writePgmImage receives an array of bytes and writes it to a pgm file.
func (io *ioState) writePgmImage() {
	_ = os.Mkdir("out", os.ModePerm)

	// Request a filename from the distributor.
	filename := <-io.channels.filename

	file, ioError := os.Create("out/" + filename + ".pgm")
	util.Check(ioError)
	defer file.Close()

	_, _ = file.WriteString("P5\n")
	//_, _ = file.WriteString("# PGM file writer by pnmmodules (https://github.com/owainkenwayucl/pnmmodules).\n")
	_, _ = file.WriteString(strconv.Itoa(io.params.ImageWidth))
	_, _ = file.WriteString(" ")
	_, _ = file.WriteString(strconv.Itoa(io.params.ImageHeight))
	_, _ = file.WriteString("\n")
	_, _ = file.WriteString(strconv.Itoa(255))
	_, _ = file.WriteString("\n")

	world := make([][]byte, io.params.ImageHeight)
	for i := range world {
		world[i] = make([]byte, io.params.ImageWidth)
	}

	for y := 0; y < io.params.ImageHeight; y++ {
		for x := 0; x < io.params.ImageWidth; x++ {
			val := <-io.channels.output // 중요, 그리고 작성하고 싶던 그 이미지를 바이트바이바이트로 이 채널로 보냄, 그리고 배열에 put,
			if val != 0 {
				fmt.Println(x, y)
			}
			world[y][x] = val
		}
	}

	for y := 0; y < io.params.ImageHeight; y++ {
		for x := 0; x < io.params.ImageWidth; x++ {
			_, ioError = file.Write([]byte{world[y][x]}) // 그리고~여기서 작성!
			util.Check(ioError)
		}
	}

	ioError = file.Sync()
	util.Check(ioError)

	fmt.Println("File", filename, "output done!")
}

// readPgmImage opens a pgm file and sends its data as an array of bytes.
func (io *ioState) readPgmImage() {

	// Request a filename from the distributor.
	filename := <-io.channels.filename // 밑에서 한거를 appropriate 한 이 채널로 다운 받고

	data, ioError := ioutil.ReadFile("images/" + filename + ".pgm")
	// 파일 이름 저장 하는법, +filename 으로 파일 이름 보내고,
	// 윗 줄에서 채널로 받은 후, 옳은 경로로 파일 지정,파일 네임,.pgm 작성.

	util.Check(ioError)

	fields := strings.Fields(string(data))

	if fields[0] != "P5" {
		panic("Not a pgm file")
	}

	width, _ := strconv.Atoi(fields[1])
	if width != io.params.ImageWidth {
		panic("Incorrect width")
	}

	height, _ := strconv.Atoi(fields[2])
	if height != io.params.ImageHeight {
		panic("Incorrect height")
	}

	maxval, _ := strconv.Atoi(fields[3])
	if maxval != 255 {
		panic("Incorrect maxval/bit depth")
	}

	image := []byte(fields[4])

	for _, b := range image {
		io.channels.input <- b // 중요한 부분! , 윗줄에서 받은 이미지를 읽고(take), 바이트바이바이트로 채널로 보낸다
		// 그리고 그 채널은 distributor 고루틴으로 전송
	}

	fmt.Println("File", filename, "input done!")
}

// startIo should be the entrypoint of the io goroutine.
// IO 고루틴 의 ENTRYPOINT
func startIo(p Params, c ioChannels) {
	io := ioState{
		params:   p,
		channels: c,
	}

	for {
		select {
		// Block and wait for requests from the distributor
		case command := <-io.channels.command:
			switch command {
			case ioInput:
				io.readPgmImage() //만약 readPgmImage 함수 호출 -> 위에 함수 선언 으로
			case ioOutput:
				io.writePgmImage()
			case ioCheckIdle:
				io.channels.idle <- true // (프로그램이 닫혀도 되는지 확인) 체크 아이들은 pgm이 프로그램이 끝내기 전까지 본인 일을 제대로 했는지 확인.
			}
		}
	}
}
