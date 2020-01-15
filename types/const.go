package types

import "fmt"

const (
	TestUser      = "419032e5-d2b4-4711-b83d-77e0aed0e832"
	ServerConnStr = "postgres://docker:docker@127.0.0.1/test1?sslmode=disable"
)

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
	ErrorTransactionExist  = fmt.Errorf("transaction already exist")
	ErrorUserNotExist      = fmt.Errorf("user does not exist")
)
