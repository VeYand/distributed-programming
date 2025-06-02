package service

type OperationType int

const (
	SetOperation OperationType = iota
	GetOperation
	ListKeysOperation
)

type Command struct {
	Type  OperationType
	Key   string
	Value int
	Reply chan Response
}

type Response struct {
	Value int
	Keys  []string
	Err   error
}
