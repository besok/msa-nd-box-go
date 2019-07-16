package storage

import (
	"fmt"
	"log"
	"strconv"
	"strings"
)

type Event string
type Name string

const (
	Put       Event = "Put"
	Get             = "Get"
	GetValue        = "GetValue"
	RemoveKey       = "RemoveKey"
	RemoveVal       = "RemoveVal"
	Init            = "Init"
	Clean           = "Clean"
)

type Listener func(event Event, storageName Name, key string, value Line)
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
func (h *ListenerHandler) Handle(event Event, storage Name, key string, value Line) {
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
	Get(idx int) (Line, bool)
	Remove(line Line) bool
	Size() int
}

type LBStrategy string

const (
	Robin  LBStrategy = "robin"
	Random LBStrategy = "random"
)

type LBLines struct {
	lines []LBLine
}

func (l *LBLines) fromString(records Records) {
	l.lines = make([]LBLine, len(records))
	for i, v := range records {
		split := strings.Split(v, ":")

		serviceName := split[0]
		strategy := LBStrategy(split[1])
		el, e := strconv.Atoi(split[2])
		if e != nil {
			el = 0
		}
		l.lines[i] = LBLine{serviceName, strategy, el}
	}
}

func (l *LBLines) ToString() Records {
	lines := make(Records, len(l.lines))
	for i, v := range l.lines {
		lines[i] = fmt.Sprintf("%s:%s:%d", v.Service, v.Strategy, v.Idx)
	}
	return lines
}

func (*LBLines) Equal(left Line, right Line) bool {
	return left.Equal(right)
}

func (l *LBLines) Add(line Line) bool {
	l.lines = append(l.lines, line.(LBLine))
	return true
}

func (l *LBLines) Get(idx int) (Line, bool) {
	sz := len(l.lines)
	if idx > sz-1 || idx < 0 {
		return nil, false
	}
	return l.lines[idx], true
}

func (l *LBLines) Remove(line Line) bool {
	idx := -1
	values := l.lines
	for i, el := range values {
		if l.Equal(el, line) {
			idx = i
			break
		}
	}
	if idx < 0 {
		return false
	}
	ln := len(values)
	values[ln-1], values[idx] = values[idx], values[ln-1]
	l.lines = values[:ln-1]
	return true
}

func (l *LBLines) Size() int {
	return len(l.lines)
}

type LBLine struct {
	Service  string
	Strategy LBStrategy
	Idx      int
}

func (l LBLine) Equal(other Line) bool {
	otherL := other.(LBLine)
	return l.Service == otherL.Service
}

type EmptyLine string

func (*EmptyLine) Equal(other Line) bool {
	return true
}

type EmptyLines struct {
	lines []EmptyLine
}

func (*EmptyLines) Get(idx int) (Line, bool) {
	return nil, false
}

func (*EmptyLines) fromString(records Records) {
}

func (*EmptyLines) ToString() Records {
	return []string{}
}

func (*EmptyLines) Equal(left Line, right Line) bool {
	return true
}

func (*EmptyLines) Add(line Line) bool {
	return true
}

func (*EmptyLines) Remove(line Line) bool {
	return true
}

func (*EmptyLines) Size() int {
	return 0
}

func CreateEmptyLines() *Lines {
	var lines Lines
	lines = new(EmptyLines)
	return &lines
}

func CreateEmptyRecords() Records {
	return make(Records, 0)
}
func CreateStringLines() Lines {
	return new(StringLines)
}

func CreateCBLines() Lines {
	return new(CBLines)
}
func CreateLBLines() Lines {
	return new(LBLines)
}
func CreateReloadLines() Lines {
	return new(ReloadLines)
}

// simple base impl
type StringLine struct {
	Value string
}
type StringLines struct {
	Lines []StringLine
}

func (l *StringLines) Get(idx int) (Line, bool) {
	if idx > l.Size()-1 || idx < 0 {
		return nil, false
	}
	return l.Lines[idx], true
}

func (l *StringLines) Size() int {
	return len(l.Lines)
}

func (l *StringLines) Add(line Line) bool {
	l.Lines = append(l.Lines, line.(StringLine))
	return true
}

func (l *StringLines) Remove(line Line) bool {
	idx := -1
	lines := l.Lines
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
	l.Lines = lines[:ln-1]
	return true
}

func (*StringLines) Equal(left Line, right Line) bool {
	return left.Equal(right)
}

func (l *StringLines) fromString(records Records) {
	l.Lines = make([]StringLine, 0)
	for _, v := range records {
		l.Lines = append(l.Lines, StringLine{v})
	}
}
func (l *StringLines) ToString() Records {
	records := make(Records, 0)
	for _, v := range l.Lines {
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
		log.Println(" Lines are a empty ")
	}
}

type CBLine struct {
	Address string
	Active  bool
}
type CBLines struct {
	Lines []CBLine
}

func (l *CBLines) Get(idx int) (Line, bool) {
	sz := len(l.Lines)
	if idx > sz-1 || idx < 0 {
		return nil, false
	}
	return l.Lines[idx], true
}

func (l *CBLines) fromString(records Records) {
	l.Lines = make([]CBLine, len(records))
	for i, v := range records {
		res := strings.Split(v, "=")
		flag := false
		if res[1] == "true" {
			flag = true
		}
		l.Lines[i] = CBLine{Address: res[0], Active: flag}
	}
}

func (l *CBLines) ToString() Records {
	records := make([]string, l.Size())
	for i, v := range l.Lines {
		records[i] = fmt.Sprintf("%s=%s", v.Address, strconv.FormatBool(v.Active))
	}
	return records
}

func (CBLines) Equal(left Line, right Line) bool {
	return left.Equal(right)
}

func (l *CBLines) Add(line Line) bool {
	l.Lines = append(l.Lines, line.(CBLine))
	return true
}

func (l *CBLines) Remove(line Line) bool {
	idx := -1
	values := l.Lines
	for i, el := range values {
		if l.Equal(el, line) {
			idx = i
			break
		}
	}
	if idx < 0 {
		return false
	}
	ln := len(values)
	values[ln-1], values[idx] = values[idx], values[ln-1]
	l.Lines = values[:ln-1]
	return true
}

func (l *CBLines) Size() int {
	return len(l.Lines)
}

func (v CBLine) Equal(other Line) bool {
	otherV := other.(CBLine)
	return v.Address == otherV.Address
}

type ReloadLine struct {
	Service string
	Address string
	Path    string
	Limit   int
	Count   int
}

type ReloadLines struct {
	Lines []ReloadLine
}

func (l *ReloadLines) Get(idx int) (Line, bool) {
	sz := len(l.Lines)
	if idx > sz-1 || idx < 0 {
		return nil, false
	}
	return l.Lines[idx], true
}
func (l *ReloadLines) fromString(records Records) {
	l.Lines = make([]ReloadLine, len(records))
	for i, v := range records {
		vals := strings.Split(v, "=")

		count, _ := strconv.Atoi(vals[3])
		limit, _ := strconv.Atoi(vals[4])
		rl := ReloadLine{
			Service: vals[0],
			Address: vals[1],
			Path:    vals[2],
			Count:   count,
			Limit:   limit,
		}
		l.Lines[i] = rl
	}
}

func (l *ReloadLines) ToString() Records {
	records := make([]string, l.Size())
	for i, v := range l.Lines {
		records[i] = fmt.Sprintf("%s=%s=%s=%d=%d", v.Service, v.Address, v.Path, v.Count, v.Limit)
	}
	return records
}
func (ReloadLines) Equal(left Line, right Line) bool {
	return left.Equal(right)
}
func (l *ReloadLines) Add(line Line) bool {
	l.Lines = append(l.Lines, line.(ReloadLine))
	return true
}
func (l *ReloadLines) Remove(line Line) bool {
	idx := -1
	values := l.Lines
	for i, el := range values {
		if l.Equal(el, line) {
			idx = i
			break
		}
	}
	if idx < 0 {
		return false
	}
	ln := len(values)
	values[ln-1], values[idx] = values[idx], values[ln-1]
	l.Lines = values[:ln-1]
	return true
}

func (l *ReloadLines) Size() int {
	return len(l.Lines)
}
func (v ReloadLine) Equal(other Line) bool {
	otherV := other.(ReloadLine)
	return v.Service == otherV.Service
}
