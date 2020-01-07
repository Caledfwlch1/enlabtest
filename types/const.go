package types

const (
	Undefine = OperationState(iota)
	Win
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

const (
	OperationOk  = "OK"
	UnknownState = "unknown state"
)

const (
	UserId = "752ca952-c89e-4f3a-9d31-8478f6b8c9c8"
)
