package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
	"github.com/ogios/merge-repo/ui/comp"
)

var (
	copyColor                = lipgloss.Color("#00bd86")
	currentLineStyle         = lipgloss.NewStyle().Foreground(lipgloss.Color("#000000"))
	currentLineSelectedStyle = currentLineStyle.Copy().Background(copyColor)
	currentLineNormalStyle   = currentLineStyle.Copy().Background(lipgloss.Color("#ffffff"))
	selectedStyle            = lipgloss.NewStyle().Foreground(copyColor)
)

func main() {
	// s := "一二三"
	// a := s + selectedStyle.Render(s) + s
	// a = s + currentLineNormalStyle.Render(a)
	// a = currentLineSelectedStyle.Render(a)
	// fmt.Println(a, len(a))
	// os.WriteFile("./test.log", []byte(a), 0766)

	// ansi, shit := comp.GetANSIs(a)
	// fmt.Println(ansi)
	// fmt.Println(json.Marshal(ansi))
	// fmt.Println(shit)
	// b, _ := json.Marshal(ansi)
	// os.WriteFile("test1.log", b, 0766)
	// os.WriteFile("test1.log", []byte(shit), 0766)

	// boundsTest()
	clipviewTest()
}

func clipviewTest() {
	b := `0123456789`
	s := b + currentLineSelectedStyle.Render(b) + currentLineSelectedStyle.Render(b) + b + "\n"
	start := time.Now()
	var a strings.Builder
	a.Grow(len(s) * 4)
	a.WriteString(s)
	a.WriteString(s)
	a.WriteString(s)
	a.WriteString(s[:len(s)-1])
	fmt.Println("cost:", time.Now().UnixMicro()-start.UnixMicro())

	start = time.Now()
	comp.ClipView(a.String(), 2, 2, 10, 10)
	fmt.Println("cost:", time.Now().UnixMicro()-start.UnixMicro())
}

func boundsTest() {
	a := []*BoundTest{
		{
			Bound: [2]int{1, 10},
		},
		{
			Bound: [2]int{15, 20},
		},
		{
			Bound: [2]int{20, 25},
		},
		{
			Bound: [2]int{30, 40},
		},
	}
	start := time.Now()
	res := search(a, 9)
	fmt.Println(res)
	res = search(a, 20)
	fmt.Println(res)
	// for i := 0; i < 60; i++ {
	// 	res := search(a, i)
	// 	fmt.Println(res)
	// }
	fmt.Println("cost:", time.Now().UnixMicro()-start.UnixMicro())
}

type BoundTest struct {
	Bound [2]int
}

func search(list []*BoundTest, pos int) []int {
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
		fmt.Println("for round:", i)
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
