package types

import "fmt"

const (
	Win = OperationState(iota + 1)
	Lost
)

var (
	StateToString = map[OperationState]string{
		Win:  "win",
		Lost: "lost"}

	StringToState = map[string]OperationState{
		"win":  Win,
		"lost": Lost}
)

var (
	ErrorUnknownState      = fmt.Errorf("unknown state")
	ErrorUnknownSourceType = fmt.Errorf("unknown Source-Type")
)
