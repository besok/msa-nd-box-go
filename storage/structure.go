package storage

import "fmt"

type Records []string
type Line interface {
	equal(other Line) bool
}
type Lines interface {
	fromString(records Records)
	ToString() Records
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
func (s StrLine) equal(other Line) bool {
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
