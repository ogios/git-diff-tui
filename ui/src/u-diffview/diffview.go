package udiffview

import (
	"bytes"
	"log"
	"path"
	"strings"

	"github.com/alecthomas/chroma/v2/lexers"
	"github.com/alecthomas/chroma/v2/quick"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/ogios/cropviewport"
	"github.com/ogios/cropviewport/process"
	"github.com/ogios/go-diffcontext"
	"github.com/ogios/merge-repo/api"
	"github.com/ogios/merge-repo/data"
	"github.com/sergi/go-diff/diffmatchpatch"
)

type DiffViewModel struct {
	cache *api.ContentCacher[*DiffContentData]
	v     tea.Model
}
type DiffContentData api.ContentData

var (
	errorMsg = lipgloss.NewStyle().
			Bold(true).
			Inline(true).
			Background(lipgloss.Color("#ff2f5a")).
			Foreground(lipgloss.Color("#000000"))
	diffInsertBg = lipgloss.NewStyle().
			Background(lipgloss.Color("#1aff66")).
			Foreground(lipgloss.Color("#33664d"))
	diffDeleteBg = lipgloss.NewStyle().
			Background(lipgloss.Color("#ff1a4d")).
			Foreground(lipgloss.Color("#990d18"))
	viewBlockStyle lipgloss.Style
)

func NewDiffViewModel(block [2]int) tea.Model {
	viewBlockStyle = lipgloss.NewStyle().Width(block[0]).Height(block[1])
	view := &DiffViewModel{
		v: cropviewport.NewCropViewportModel(),
		cache: api.NewContentCacher(func(p string) *DiffContentData {
			finalContent, err := diffContent(p, p)
			if err != nil {
				finalContent = errorMsg.Render(err.Error())
			}
			at, sl := process.ProcessContent(finalContent)
			return &DiffContentData{
				Table: at,
				Lines: sl,
			}
		}),
	}
	view.v.(*cropviewport.CropViewportModel).SetBlock(0, 0, block[0], block[1])
	return view
}

func (v *DiffViewModel) ViewFile(p string) {
	content := v.cache.Get(p)
	cv := v.v.(*cropviewport.CropViewportModel)
	cv.SetContentGivenData(content.Table, content.Lines)
	cv.BackToTop()
	cv.BackToLeft()
}

func (v *DiffViewModel) Init() tea.Cmd {
	return nil
}

func (v *DiffViewModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	vp, cmd := v.v.Update(msg)
	v.v = vp
	// cmds = append(cmds, viewport.Sync(v.v), cmd)
	cmds = append(cmds, cmd)
	return v, tea.Batch(cmds...)
}

func (v *DiffViewModel) View() string {
	return viewBlockStyle.Render(v.v.View())
}

func diffContent(p1, p2 string) (string, error) {
	isLayout2 := strings.Contains(p1, "layout2")
	code1, err := data.GetTempDiffFile(p1)
	if err != nil {
		return "", err
	}
	code2, err := data.GetDiffSrcFile(p1)
	if err != nil {
		return "", err
	}
	dmp := diffmatchpatch.New()
	diffs := dmp.DiffMain(string(code1), string(code2), true)
	diffs = dmp.DiffCleanupSemantic(diffs)
	diffs = dmp.DiffCleanupEfficiency(diffs)
	dc := diffcontext.New()
	dc.AddDiffs(diffs)

	c1, err := highlight(string(code1), p1)
	if err != nil {
		return "", err
	}
	linesC1 := strings.Split(c1, "\n")

	c2, err := highlight(string(code2), p2)
	if err != nil {
		return "", err
	}
	linesC2 := strings.Split(c2, "\n")
	if isLayout2 {
		log.Println(p1, "c1", len(strings.Split(string(code1), "\n")), len(linesC1))
		log.Println(p2, "c2", len(strings.Split(string(code2), "\n")), len(linesC2))
		log.Println("dclines:", len(dc.Lines))
	}
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
			be := []byte(diffDeleteBg.Render(linesC1[i1]))
			af := []byte(diffInsertBg.Render(linesC2[i2]))

			// be := []byte(linesC1[i1])
			// af := []byte(linesC2[i2])
			log.Println(string(be))
			log.Println(string(af))
			switch dl.State {
			case diffcontext.DiffChanged:
				dl.Before, dl.After = be, af
				i1++
				i2++
			case diffmatchpatch.DiffInsert:
				dl.After = af
				i2++
			case diffmatchpatch.DiffDelete:
				dl.Before = be
				i1++
			}
		}
	}
	return dc.GetMixed(), nil
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
