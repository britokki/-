package gol

import (
	"fmt"
	"uk.ac.bris.cs/gameoflife/util"
)

// Event represents any Game of Life event that needs to be communicated to the user.
type Event interface {
	// Stringer allows each event to be printed by the GUI
	fmt.Stringer
	// GetCompletedTurns should return the number of fully completed turns.
	// If the 0th turn is finished, this should return 1.

	// 이 인터페이스는 스트링을 반환하는 하나의 메소드가 있다 (Stringer) -> 따라서 밑에 있는 Any 타입들을 지원하는 스트링 메소드를
	// 구현해야 함

	GetCompletedTurns() int // 얘를 구현해야 함!

}

// AliveCellsCount is an Event notifying the user about the number of currently alive cells.
// This Event should be sent every 2s.
type AliveCellsCount struct { // implements Event
	CompletedTurns int
	CellsCount     int // 얘네 타입 지원하는 스트링 메소드
}

// ImageOutputComplete is an Event notifying the user about the completion of output.
// This Event should be sent every time an image has been saved.
type ImageOutputComplete struct { // implements Event
	CompletedTurns int
	Filename       string
}

// State represents a change in the state of execution.
type State int

const (
	Paused State = iota
	Executing
	Quitting
)

// StateChange is an Event notifying the user about the change of state of execution.
// This Event should be sent every time the execution is paused, resumed or quit.
type StateChange struct { // implements Event
	CompletedTurns int
	NewState       State
}

// CellFlipped is an Event notifying the GUI about a change of state of a single cell.
// This even should be sent every time a cell changes state.
// Make sure to send this event for all cells that are alive when the image is loaded in.
type CellFlipped struct { // implements Event
	CompletedTurns int
	Cell           util.Cell
}

// TurnComplete is an Event notifying the GUI about turn completion.
// SDL will render a frame when this event is sent.
// All CellFlipped events must be sent *before* TurnComplete.
type TurnComplete struct { // implements Event
	CompletedTurns int
}

// FinalTurnComplete is an Event notifying the testing framework about the new world state after execution finished.
// The data included with this Event is used directly by the tests.
// SDL closes the window when this Event is sent.
type FinalTurnComplete struct {
	CompletedTurns int
	Alive          []util.Cell // alive cells! , X,Y 좌표 포함 하고 있음! , Cell 눌러서 읽는걸 추천(시온피셜
}

// String methods allow the different types of Events and States to be printed.

func (state State) String() string {
	switch state {
	case Paused:
		return "Paused"
	case Executing:
		return "Executing"
	case Quitting:
		return "Quitting"
	default:
		return "Incorrect State"
	}
}

func (event StateChange) String() string {
	return fmt.Sprintf("%v", event.NewState)
}

func (event StateChange) GetCompletedTurns() int {
	return event.CompletedTurns
}

func (event AliveCellsCount) String() string {
	return fmt.Sprintf("Alive Cells %v", event.CellsCount)
}

func (event AliveCellsCount) GetCompletedTurns() int {
	return event.CompletedTurns
}

func (event ImageOutputComplete) String() string {
	return fmt.Sprintf("File %v output complete", event.Filename)
}

func (event ImageOutputComplete) GetCompletedTurns() int {
	return event.CompletedTurns
}

func (event CellFlipped) String() string {
	return fmt.Sprintf("CellFlipped%v", event.Cell) // 확인 요망
}

func (event CellFlipped) GetCompletedTurns() int {
	return event.CompletedTurns
}

func (event TurnComplete) String() string {
	return fmt.Sprintf("Turncmpleated%v", event.CompletedTurns) //확인 요망
}

func (event TurnComplete) GetCompletedTurns() int {
	return event.CompletedTurns
}

// Stringer allows each event to be printed by the GUI

// GetCompletedTurns should return the number of fully completed turns.
// If the 0th turn is finished, this should return 1.

// 이 인터페이스는 스트링을 반환하는 하나의 메소드가 있다 (Stringer) -> 따라서 밑에 있는 Any 타입들을 지원하는 스트링 메소드를
// 구현해야 함
func (event FinalTurnComplete) String() string {
	return fmt.Sprintf("FinalTurnCompleted %v", event.CompletedTurns)
	//이 Event가  distributor에서 testgol 로 전송되어야 함.
	// 확인 요망
}

func (event FinalTurnComplete) GetCompletedTurns() int {
	return event.CompletedTurns //확인 요망(건드린건 없음)

}

// This might all seem like weird syntax to you...
// You have however seen something similar to it before in first year.

// In the Go code an Interface called Event is created, this provides a set of methods that
// need to be defined for something to have the type Event.

// This is a similar concept to typeclasses in Haskell. A typeclass called Event could be defined.
// It would require two methods to be implemented: string and getCompletedTurns. Note the
// similarities between the type signatures of the Go and Haskell functions.

/*
> class Event event where
>   string :: event -> String
>   getCompletedTurns :: event -> Int
*/

// A new data type called ImageOutputComplete can then be created, just like in Go.

/*
> data ImageOutputComplete = ImageOutputComplete Int String
*/

// Now in the Go code extension methods are created for the ImageOutputComplete so that it
// provides the methods required for the Event Inteface. Similarly, in Haskell, an instance
// of the typeclass Event can be created.

/*
> instance Event ImageOutputComplete where
>   string (ImageOutputComplete t f) = concat ["Turn ", show t, " - File ", f, " output complete"]
>   getCompletedTurns (ImageOutputComplete t f) = t
*/
