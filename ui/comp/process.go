package comp

import (
	"strings"
)

type ANSITable struct {
	sub   *ANSITable
	data  []byte
	bound [2]int
}

type ANSIStackItem struct {
	data       []byte
	startIndex int
}

type ANSITableList []*ANSITable

const (
	ESCAPE_SEQUENCE     = '\x1b'
	ESCAPE_SEQUENCE_END = string(ESCAPE_SEQUENCE) + "[0m"
)

func GetANSIs(s string) (ANSITableList, string) {
	// preserve normal string
	var normalString strings.Builder
	normalString.Grow(len(s))

	// preserve ansi string and position
	tables := make(ANSITableList, 0)
	ansiStack := make([]*ANSIStackItem, 0)
	ansi := false
	i := 0
	for _, v := range s {
		// met `esc` char
		if v == ESCAPE_SEQUENCE {
			ansi = true
			ansiStack = append(ansiStack, &ANSIStackItem{
				startIndex: i,
				data:       []byte{byte(v)},
			})
		} else {
			// in ansi sequence content
			if ansi {
				last := ansiStack[len(ansiStack)-1]
				last.data = append(last.data, byte(v))
				// end of an ansi sequence. terminate
				if v == 'm' {
					ansi = false
					// ends all ansi sequences in stack
					if string(last.data) == ESCAPE_SEQUENCE_END {
						if len(ansiStack) > 1 {
							tables = append(tables, stackToTable(ansiStack, i))
						}
						ansiStack = make([]*ANSIStackItem, 0)
					}
				}
			} else {
				// normal content
				normalString.WriteRune(v)
				i++
			}
		}
	}
	return tables, normalString.String()
}

func stackToTable(stack []*ANSIStackItem, endIndex int) *ANSITable {
	first := stack[0]
	root := &ANSITable{
		bound: [2]int{
			first.startIndex,
			endIndex,
		},
		data: first.data,
	}
	temp := root
	for _, v := range stack[1:] {
		temp.sub = &ANSITable{
			bound: [2]int{
				v.startIndex,
				endIndex,
			},
			data: v.data,
		}
		temp = temp.sub
	}
	return root
}

// func ClipView(s string, block [4]int) string {
// }
