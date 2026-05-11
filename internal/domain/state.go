package domain

type State string

const (
	StateSuccess State = "SUCCESS"
	StateFail    State = "FAIL"
	StateDisqual State = "DISQUAL"
)
