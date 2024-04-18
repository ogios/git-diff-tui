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
	Sub   *ANSITable
	Data  []byte
	Bound [2]int // rune index
}

// implement `BoundsStruct` for search
func (a *ANSITable) getBounds() [2]int {
	return a.Bound
}

type ANSITableList struct {
	L []BoundsStruct
}

type ANSIQueueItem struct {
	data       []byte
	startIndex int
}

const (
	ESCAPE_SEQUENCE     = '\x1b'
	ESCAPE_SEQUENCE_END = string(ESCAPE_SEQUENCE) + "[0m"
)

// split `string with ansi` into `ansi sequences` and `raw string`
func GetANSIs(s string) (*ANSITableList, string) {
	// preserve normal string
	var normalString strings.Builder
	normalString.Grow(len(s))

	// preserve ansi string and position
	tables := make([]BoundsStruct, 0)
	ansiQueue := make([]*ANSIQueueItem, 0)
	ansi := false
	// NOTE: do not use `for i := range string` index since it's not i+=1 but i+=byte_len
	// solution: transform s into []rune or use custom variable for index
	i := 0
	for _, v := range s {
		// meet `esc` char
		if v == ESCAPE_SEQUENCE {
			// enable ansi mode until meet 'm'
			ansi = true
			// using utf8 rune function
			// but maybe just byte(v) is enough since ansi only contains rune of one byte?
			byteData := []byte{}
			byteData = utf8.AppendRune(byteData, v)
			ansiQueue = append(ansiQueue, &ANSIQueueItem{
				startIndex: i,
				data:       slices.Clip(byteData),
			})
		} else {
			// in ansi sequence content mode
			if ansi {
				last := ansiQueue[len(ansiQueue)-1]
				last.data = utf8.AppendRune(last.data, v)
				// end of an ansi sequence. terminate
				if v == 'm' {
					ansi = false
					// clip cap
					last.data = slices.Clip(last.data)
					// ends all ansi sequences in queue and create ansi table
					if string(last.data) == ESCAPE_SEQUENCE_END {
						// skip if ansi queue only contain "[0m", which means no ansi actually working
						if len(ansiQueue) > 1 {
							table := queueToTable(ansiQueue[:len(ansiQueue)-1], i)
							tables = append(tables, table)
						}
						// reset queue
						ansiQueue = make([]*ANSIQueueItem, 0)
					}
				}
			} else {
				// normal content
				normalString.WriteRune(v)
				i++
			}
		}
	}
	return &ANSITableList{
		L: slices.Clip(tables),
	}, normalString.String()
}

// transform queue into ansi table which contains all ansi sequences from start to end
func queueToTable(queue []*ANSIQueueItem, endIndex int) *ANSITable {
	first := queue[0]
	root := &ANSITable{
		Bound: [2]int{
			first.startIndex,
			endIndex,
		},
		Data: first.data,
	}

	// add to sub
	temp := root
	for _, v := range queue[1:] {
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

// var EMPTY_ANSITABLELIST = make([]*ANSITable, 0)
var EMPTY_ANSITABLELIST = make([]BoundsStruct, 0)

// get a slice of ansi table, it will find all tables between `startIndex` and `endIndex`
func (a *ANSITableList) GetSlice(startIndex, endIndex int) []BoundsStruct {
	if len(a.L) == 0 {
		return a.L
	}
	var start, end int
	temp := search(a.L, startIndex)
	fmt.Println("start temp:", temp)
	// len == 1 means index within a specific table
	// len == 2 means index between two tables, and for `startIndex` we only need the tables after `startIndex`
	// temp[1] == -1 means already at the front of tablelist and no matchs
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
	// len == 1 means index within a specific table
	// len == 2 means index between two tables, and for `endIndex` we only need the tables before `endIndex`
	// temp[1] == -1 means already at the front of tablelist and no matchs
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

	// get slice of tablelist between start and end
	return a.L[start : end+1]
}

type SubLine struct {
	Data  *RuneDataList
	Bound [2]int
}

type RuneDataList struct {
	L          []BoundsStruct
	TotalWidth int
}

// init RuneData list given runes
//
// RuneDataList can only be set with this function, no more process allowed afterwards
func (r *RuneDataList) Init(s []rune) *RuneDataList {
	r.L = make([]BoundsStruct, len(s))
	visibleIndex := 0
	// for every rune, get its width, start and end index refers to the visible line
	// and save rune data into bytes
	for i, v := range s {
		bs := []byte{}
		bs = utf8.AppendRune(bs, v)
		w := runewidth.RuneWidth(v)
		r.L[i] = &RuneData{
			Byte:  slices.Clip(bs),
			Bound: [2]int{visibleIndex, visibleIndex + w},
		}
		r.TotalWidth += w
		visibleIndex += w
	}
	return r
}

type RuneData struct {
	Byte  []byte
	Bound [2]int // refers to the visible width
}

func (r *RuneData) getBounds() [2]int {
	return r.Bound
}

const LINE_SPLIT = "\n"

var SPACE_HODLER = []byte(" ")

func ClipView(s string, x, y, width, height int) string {
	atablelist, raw := GetANSIs(s)
	for _, v := range atablelist.L {
		fmt.Println("table: ", v)
	}
	rawlines := strings.Split(raw, LINE_SPLIT)
	sublines := make([]*SubLine, len(rawlines))
	index := 0
	for i, v := range rawlines {
		data := (&RuneDataList{}).Init([]rune(v))
		lastIndex := index + len(data.L)
		sublines[i] = &SubLine{
			Bound: [2]int{index, lastIndex},
			Data:  data,
		}
		index = lastIndex + 1
	}
	lines := api.SliceFrom(sublines, y, y+height)
	fmt.Println(lines, lines[0].Data)
	clipLines(lines, atablelist, x, y, width, height)
	return ""
}

func clipLines(lines []*SubLine, atablelist *ANSITableList, x, y, width, height int) {
	var buf strings.Builder
	buf.Grow((width + 1) * height)
	// lines
	for lineIndex, sl := range lines {
		// if x is within the width of line
		if sl.Data.TotalWidth-1 >= x {
			// (x) for range every rune and count width
			// (✓) binary search for a range of rune
			var start, end int
			temp := search(sl.Data.L, x)
			start = temp[0]
			temp = search(sl.Data.L, x+width)
			if len(temp) == 1 {
				end = temp[0]
			} else if len(temp) == 2 {
				end = temp[0]
			}
			lineRunes := sl.Data.L[start : end+1]
			fmt.Println("indexes:", start, end)
			fmt.Println("lineRunes:", lineRunes, lineRunes[0])

			// start from lineRunes start
			index := 0
			// atable slice
			atables := atablelist.GetSlice(start+sl.Bound[0], end+sl.Bound[0])
			fmt.Println("tables:", atables)
			// every table
			for _, a := range atables {
				// table's sub tables
				temp := a.(*ANSITable)
				fmt.Println("temp table:", temp)
				endIndex := temp.Bound[1] - sl.Bound[0] - start
				for temp != nil {
					startIndex := temp.Bound[0] - sl.Bound[0] - start
					// before table startIndex
					if startIndex > index {
						subRuneDatas := lineRunes[index:startIndex]
						for _, runeData := range subRuneDatas {
							r := runeData.(*RuneData)
							buf.Write(r.Byte)
						}
						index += len(subRuneDatas)
					}
					// ansi insert
					buf.Write(temp.Data)
					// assign sub table
					temp = temp.Sub
				}
				// add rest
				subRuneDatas := lineRunes[index:endIndex]
				for _, runeData := range subRuneDatas {
					r := runeData.(*RuneData)
					buf.Write(r.Byte)
				}
				index += len(subRuneDatas)
				// add end escape
				buf.WriteString(ESCAPE_SEQUENCE_END)
			}
			// add rest
			if index <= len(lineRunes)-1 {
				// buf.Write(lineRunes[index:])
				subRuneDatas := lineRunes[index:]
				for _, runeData := range subRuneDatas {
					r := runeData.(*RuneData)
					buf.Write(r.Byte)
				}
			}
		}
		// line break
		if lineIndex < len(lines)-1 {
			buf.WriteString(LINE_SPLIT)
		}
	}
	fmt.Println(buf.String())
	fmt.Println("done")
}
