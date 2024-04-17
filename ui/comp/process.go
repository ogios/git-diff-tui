package comp

import (
	"fmt"
	"slices"
	"strings"
	"unicode/utf8"

	"github.com/mattn/go-runewidth"
	"github.com/ogios/merge-repo/api"
)

type ANSITable struct {
	Sub  *ANSITable
	Data []byte

	// rune index
	Bound [2]int
}

type ANSITableList struct {
	L []*ANSITable
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
			// using utf8 rune function
			// but maybe just byte(v) is enough since ansi only contains rune of one byte?
			byteData := []byte{}
			byteData = utf8.AppendRune(byteData, v)
			ansiStack = append(ansiStack, &ANSIStackItem{
				startIndex: i,
				data:       byteData,
			})
		} else {
			// in ansi sequence content
			if ansi {
				last := ansiStack[len(ansiStack)-1]
				last.data = utf8.AppendRune(last.data, v)
				// last.data = append(last.data, byte(v))
				// end of an ansi sequence. terminate
				if v == 'm' {
					ansi = false
					// clip cap
					last.data = slices.Clip(last.data)
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
		L: tables,
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
	if len(a.L) == 0 {
		return a.L
	}
	var start, end int
	temp := search(a.L, startIndex)
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

	temp = search(a.L, endIndex)
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

	return a.L[start : end+1]
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
	// Data  []byte
	Data  *RuneDataList
	Bound [2]int
}

type RuneDataList struct {
	L          []*RuneData
	TotalWidth int
}

func (r *RuneDataList) Init(s string) *RuneDataList {
	r.L = make([]*RuneData, len(s))
	for i, v := range s {
		bs := []byte{}
		utf8.AppendRune(bs, v)
		w := runewidth.RuneWidth(v)
		r.L[i] = &RuneData{
			Byte:  slices.Clip(bs),
			Width: w,
		}
		r.TotalWidth += w
	}
	return r
}

type RuneData struct {
	Byte  []byte
	Width int
}

const LINE_SPLIT = "\n"

func ClipView(s string, x, y, width, height int) string {
	atablelist, raw := GetANSIs(s)
	rawlines := strings.Split(raw, LINE_SPLIT)
	sublines := make([]*SubLine, len(rawlines))
	index := 0
	for i, v := range rawlines {
		data := (&RuneDataList{}).Init(v)
		lastIndex := index + len(data.L)
		sublines[i] = &SubLine{
			Bound: [2]int{index, lastIndex},
			Data:  data,
		}
		index = lastIndex + 1
	}
	lines := api.SliceFrom(sublines, y, y+height)
	var buf strings.Builder
	buf.Grow((width + 1) * height)
	// lines
	for _, sl := range lines {
		if len(sl.Data)-1 >= x {
			index := 0
			// lineSlice := sl.Data[x:]
			// atable slice
			atables := atablelist.GetSlice(sl.Bound[0], sl.Bound[1])
			// every table
			for _, a := range atables {
				// table's sub tables
				temp := a
				endIndex := temp.Bound[1] - sl.Bound[0]
				for temp != nil {
					startIndex := temp.Bound[0] - sl.Bound[0]
					// before table startIndex
					if startIndex > index {
						normalSubString := sl.Data[index:startIndex]
						buf.Write(normalSubString)
						index += len(normalSubString)
					}
					// ansi insert
					buf.Write(temp.Data)
					// assign sub table
					temp = a.Sub
				}
				// add rest
				normalSubString := sl.Data[index:endIndex]
				buf.Write(normalSubString)
				index += len(normalSubString)
				// add end escape
				buf.WriteString(ESCAPE_SEQUENCE_END)
			}
			// add rest
			if index < len(sl.Data)-1 {
				buf.Write(sl.Data[index:])
			}
		}
		// line break
		buf.WriteString(LINE_SPLIT)
	}
	return buf.String()
}

func clipLines(lines []*SubLine, atablelist *ANSITableList, x, y, width, height int) {
	var buf strings.Builder
	buf.Grow((width + 1) * height)
	// lines
	for _, sl := range lines {
		// if x is within the width of line
		if sl.Data.TotalWidth-1 >= x {
			// for range every rune and count width

			// NOTE: ignored code down below
			index := 0
			// lineSlice := sl.Data[x:]
			// atable slice
			atables := atablelist.GetSlice(sl.Bound[0], sl.Bound[1])
			// every table
			for _, a := range atables {
				// table's sub tables
				temp := a
				endIndex := temp.Bound[1] - sl.Bound[0]
				for temp != nil {
					startIndex := temp.Bound[0] - sl.Bound[0]
					// before table startIndex
					if startIndex > index {
						normalSubString := sl.Data[index:startIndex]
						buf.Write(normalSubString)
						index += len(normalSubString)
					}
					// ansi insert
					buf.Write(temp.Data)
					// assign sub table
					temp = a.Sub
				}
				// add rest
				normalSubString := sl.Data[index:endIndex]
				buf.Write(normalSubString)
				index += len(normalSubString)
				// add end escape
				buf.WriteString(ESCAPE_SEQUENCE_END)
			}
			// add rest
			if index < len(sl.Data)-1 {
				buf.Write(sl.Data[index:])
			}
		}
		// line break
		buf.WriteString(LINE_SPLIT)
	}
}
