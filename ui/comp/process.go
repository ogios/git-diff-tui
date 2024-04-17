package comp

import (
	"fmt"
	"strings"

	"github.com/ogios/merge-repo/api"
)

type ANSITable struct {
	Sub   *ANSITable `json:"sub"`
	Data  []byte     `json:"data"`
	Bound [2]int     `json:"bound"`
}

type ANSITableList struct {
	l []*ANSITable
}

type ANSIStackItem struct {
	data       []byte
	startIndex int
}

const (
	ESCAPE_SEQUENCE     = '\x1b'
	ESCAPE_SEQUENCE_END = string(ESCAPE_SEQUENCE) + "[0m"
)

func GetANSIs(s string) (ANSITableList, string) {
	// preserve normal string
	var normalString strings.Builder
	normalString.Grow(len(s))

	// preserve ansi string and position
	tables := make([]*ANSITable, 0)
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
							table := stackToTable(ansiStack[:len(ansiStack)-1], i)
							tables = append(tables, table)
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
	return ANSITableList{
		l: tables,
	}, normalString.String()
}

func stackToTable(stack []*ANSIStackItem, endIndex int) *ANSITable {
	first := stack[0]
	root := &ANSITable{
		Bound: [2]int{
			first.startIndex,
			endIndex,
		},
		Data: first.data,
	}
	temp := root
	for _, v := range stack[1:] {
		temp.Sub = &ANSITable{
			Bound: [2]int{
				v.startIndex,
				endIndex,
			},
			Data: v.data,
		}
		temp = temp.Sub
	}
	return root
}

const (
	MODE_START = 0
	MODE_END   = 1
)

var EMPTY_ANSITABLELIST = make([]*ANSITable, 0)

func (a *ANSITableList) GetSlice(startIndex, endIndex int) []*ANSITable {
	if len(a.l) == 0 {
		return a.l
	}
	var start, end int
	temp := search(a.l, startIndex)
	fmt.Println("start temp:", temp)
	if len(temp) == 1 {
		start = temp[0]
	} else if len(temp) == 2 {
		if temp[1] == -1 {
			return EMPTY_ANSITABLELIST
		} else {
			start = temp[1]
		}
	}

	temp = search(a.l, endIndex)
	fmt.Println("end temp:", temp)
	if len(temp) == 1 {
		end = temp[0]
	} else if len(temp) == 2 {
		if temp[0] == -1 {
			return EMPTY_ANSITABLELIST
		} else {
			end = temp[0]
		}
	}
	fmt.Println("start and end index:", start, end)

	return a.l[start : end+1]
}

func search(list []*ANSITable, pos int) []int {
	listLen := len(list)
	i := (listLen - 1) / 2
	step := i
	halfStep := func() {
		step /= 2
		if step == 0 {
			step = 1
		}
	}
	halfStep()
	// 1. index between one bounds start & end
	// or
	// 2.between ( last end and current start ) or ( current end and next start )
	for {
		// fmt.Println("for round:", i)
		v := list[i]
		// between bounds
		if v.Bound[0] <= pos && v.Bound[1] > pos {
			return []int{i}
		} else {
			// smaller than start
			if pos < v.Bound[0] {
				if i > 0 {
					// not the first one
					prev := list[i-1]
					if prev.Bound[1] <= pos {
						// i bigger than prev end and i smaller than current start means circumstance 2
						return []int{i - 1, i}
					} else {
						// i smaller than prev end means still space to go left
						// i = i - int(math.Floor(float64(i)/2))
						i -= step
						halfStep()
					}
				} else {
					// first one and i smaller than first start means circumstance 2
					return []int{-1, i}
				}
			} else if pos >= v.Bound[1] {
				// bigger than end
				if i < listLen-1 {
					// not the last one
					next := list[i+1]
					if pos < next.Bound[0] {
						// i bigger than current end and smaller than next start means circumstance 2
						return []int{i, i + 1}
					} else {
						// i bigger or equal to next start means still space to go right
						// i = i + int(math.Ceil(float64(i)/2))
						// i = i + int(math.Floor(float64(listLen-i)/2))
						i += step
						halfStep()
					}
				} else {
					// last one and i bigger than end means circumstance 2
					// return []int{i, i + 1}
					return []int{i, -1}
				}
			}
		}
	}
}

type SubLine struct {
	Data  string
	Bound [2]int
}

func ClipView(s string, x, y, width, height int) string {
	atablelist, raw := GetANSIs(s)
	lines := strings.Split(raw, "\n")
	sublines := make([]*SubLine, len(lines))
	index := 0
	for i, v := range lines {
		lastIndex := index + len(v)
		sublines[i] = &SubLine{
			Bound: [2]int{index, lastIndex},
			Data:  v,
		}
		index = lastIndex + 1
	}
	slice := api.SliceFrom(sublines, y, y+height)
	var buf strings.Builder
	buf.Grow((width + 1) * height)
	for _, sl := range slice {
		atables := atablelist.GetSlice(sl.Bound[0], sl.Bound[1])
		index := 0
		for _, a := range atables {
			startIndex := a.Bound[0] - sl.Bound[0]
			endIndex := a.Bound[1] - sl.Bound[0]
		}
		fmt.Println(atables)
	}
	return ""
}
