package message

type Status int

const (
	Fail Status = iota
	Ready
	Run
	Done
)


type ServerMessage struct {
	Service Service
}

type Service struct {
	Address string
	Service string
}
