package udiffview

import (
	"bytes"
	"fmt"
	"path"
	"strings"

	"github.com/alecthomas/chroma/v2/lexers"
	"github.com/alecthomas/chroma/v2/quick"
	"github.com/muesli/termenv"
	process "github.com/ogios/ansisgr-process"
	"github.com/ogios/cropviewport"
	"github.com/ogios/go-diffcontext"
	"github.com/ogios/merge-repo/data"
	"github.com/sergi/go-diff/diffmatchpatch"
)

func diffContent(p1, p2 string) (*process.ANSITableList, []*cropviewport.SubLine, error) {
	code1, err := data.GetTempDiffFile(p1)
	if err != nil {
		return nil, nil, err
	}
	code2, err := data.GetDiffSrcFile(p1)
	if err != nil {
		return nil, nil, err
	}
	dmp := diffmatchpatch.New()
	diffs := dmp.DiffMain(string(code1), string(code2), true)
	diffs = dmp.DiffCleanupSemantic(diffs)
	diffs = dmp.DiffCleanupEfficiency(diffs)
	dc := diffcontext.New()
	dc.AddDiffs(diffs)
	_, records := dc.GetMixedLinesAndStateRecord()

	err = highlightCodes(dc, code1, code2, p1, p2)
	if err != nil {
		return nil, nil, err
	}
	at, sl := highlightDiffLines(dc.GetMixed(), records)
	return at, sl, err
}

var (
	redBG   = []byte(fmt.Sprintf("%s%sm", termenv.CSI, termenv.RGBColor("#991a1a").Sequence(true)))
	greenBG = []byte(fmt.Sprintf("%s%sm", termenv.CSI, termenv.RGBColor("#008033").Sequence(true)))
)

func highlightDiffLines(content string, records [][3]int) (*process.ANSITableList, []*cropviewport.SubLine) {
	at, sl := cropviewport.ProcessContent(content)
	for _, v := range records {
		var color []byte
		switch v[0] {
		case int(diffmatchpatch.DiffDelete):
			color = redBG
		case int(diffmatchpatch.DiffInsert):
			color = greenBG
		}
		at.SetStyle(color, v[1], v[2])
	}
	return at, sl
}

func highlightCodes(dc *diffcontext.DiffConstractor, code1, code2 []byte, p1, p2 string) error {
	c1, err := highlight(string(code1), p1)
	if err != nil {
		return err
	}
	linesC1 := strings.Split(c1, "\n")
	// linesC1 := strings.Split(string(code1), "\n")

	c2, err := highlight(string(code2), p2)
	if err != nil {
		return err
	}
	linesC2 := strings.Split(c2, "\n")
	// linesC2 := strings.Split(string(code2), "\n")
	i1 := 0
	i2 := 0
	for _, dl := range dc.Lines {
		switch dl.State {
		case diffmatchpatch.DiffEqual:
			be := []byte(linesC1[i1])
			dl.Before, dl.After = be, be
			i1++
			i2++
		default:
			switch dl.State {
			case diffcontext.DiffChanged:
				be := []byte(linesC1[i1])
				af := []byte(linesC2[i2])
				dl.Before, dl.After = be, af
				i1++
				i2++
			case diffmatchpatch.DiffInsert:
				af := []byte(linesC2[i2])
				dl.After = af
				i2++
			case diffmatchpatch.DiffDelete:
				be := []byte(linesC1[i1])
				dl.Before = be
				i1++
			}
		}
	}
	return nil
}

func highlight(content, p string) (string, error) {
	buf := new(bytes.Buffer)
	lex := lexers.Match(path.Base(p))
	lang := "plaintext"
	if lex != nil {
		lang = lex.Config().Name
	}
	err := quick.Highlight(buf, string(content), lang, "terminal16m", "catppuccin-mocha")
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}
