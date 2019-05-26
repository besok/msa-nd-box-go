package message

const (
	fail  = "fail"
	run   = "run"
	ready = "ready"
	done  = "done"
)

type ServerMessage struct {
	Service Service
}

type Service struct {
	Address string
	Service string
}
