package storage

import "fmt"

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
}

func CreateEmptyRecords() Records {
	return make(Records, 0)
}
func CreateStrLines() Lines {
	return new(StrLines)
}

// simple base impl
type StrLine struct {
	value string
}
type StrLines struct {
	lines []StrLine
}

func (l *StrLines) Add(line Line) bool{
	l.lines = append(l.lines,line.(StrLine))
	return true
}

func (l *StrLines) Remove(line Line) bool {
	idx := -1
	lines := l.lines
	for i, el := range lines {
		if l.Equal(el, line) {
			idx = i
			break
		}
	}
	if idx < 0{
		return false
	}
	ln := len(lines)
	lines[ln-1], lines[idx] = lines[idx], lines[ln-1]
	l.lines = lines[:ln-1]
	return true
}

func (l *StrLines) Equal(left Line, right Line) bool {
	return left.Equal(right)
}

func (l *StrLines) fromString(records Records) {
	l.lines = make([]StrLine, 0)
	for _, v := range records {
		l.lines = append(l.lines, StrLine{v})
	}
}
func (l *StrLines) ToString() Records {
	records := make(Records, 0)
	for _, v := range l.lines {
		records = append(records, v.value)
	}
	return records
}
func (s StrLine) Equal(other Line) bool {
	oth, ok := other.(StrLine)
	if !ok {
		return false
	}
	return s.value == oth.value
}

func PrintLines(lines Lines) {
	records := lines.ToString()
	for _, r := range records {
		fmt.Printf(" record: %s \n", r)
	}
}
