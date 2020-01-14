package types

import "fmt"

const TestUser = "419032e5-d2b4-4711-b83d-77e0aed0e832"
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
