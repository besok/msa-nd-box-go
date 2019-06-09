package storage

import (
	"log"
)

type StorageEvent string
type StorageName string

const (
	Put       StorageEvent = "Put"
	Get                    = "Get"
	RemoveKey              = "RemoveKey"
	RemoveVal              = "RemoveVal"
	Init                   = "Init"
	Clean                  = "Clean"
)

type Listener func(event StorageEvent, storageName StorageName, key string, value Line)
type ListenerHandler struct {
	listeners []Listener
}

func CreateListenerHandler() ListenerHandler {
	return ListenerHandler{make([]Listener, 0)}
}

func (h *ListenerHandler) AddListener(l Listener) {
	log.Println("add listener ")
	h.listeners = append(h.listeners, l)
}
func (h *ListenerHandler) Handle(event StorageEvent, storage StorageName, key string, value Line) {
	for _, listener := range h.listeners {
		listener(event, storage, key, value)
	}
}

type Records []string
type Line interface {
	Equal(other Line) bool
}
type Lines interface {
	fromString(records Records)
	ToString() Records
	Equal(left Line, right Line) bool
	Add(line Line) bool
	Remove(line Line) bool
	Size() int
}

func CreateEmptyRecords() Records {
	return make(Records, 0)
}
func CreateStringLines() Lines {
	return new(StringLines)
}

// simple base impl
type StringLine struct {
	Value string
}
type StringLines struct {
	lines []StringLine
}

func (l *StringLines) Size() int {
	return len(l.lines)
}

func (l *StringLines) Add(line Line) bool {
	l.lines = append(l.lines, line.(StringLine))
	return true
}

func (l *StringLines) Remove(line Line) bool {
	idx := -1
	lines := l.lines
	for i, el := range lines {
		if l.Equal(el, line) {
			idx = i
			break
		}
	}
	if idx < 0 {
		return false
	}
	ln := len(lines)
	lines[ln-1], lines[idx] = lines[idx], lines[ln-1]
	l.lines = lines[:ln-1]
	return true
}

func (l *StringLines) Equal(left Line, right Line) bool {
	return left.Equal(right)
}

func (l *StringLines) fromString(records Records) {
	l.lines = make([]StringLine, 0)
	for _, v := range records {
		l.lines = append(l.lines, StringLine{v})
	}
}
func (l *StringLines) ToString() Records {
	records := make(Records, 0)
	for _, v := range l.lines {
		records = append(records, v.Value)
	}
	return records
}
func (s StringLine) Equal(other Line) bool {
	oth, ok := other.(StringLine)
	if !ok {
		return false
	}
	return s.Value == oth.Value
}

func PrintLines(lines Lines) {
	if lines != nil {
		records := lines.ToString()
		for _, r := range records {
			log.Printf(" record: %s \n", r)
		}
	} else {
		log.Println(" lines are a empty ")
	}
}
