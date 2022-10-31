package components

import (
	"fmt"
	"github.com/rivo/tview"
	"strconv"
)

func NewSlideShow() *SlideShow {
	s := &SlideShow{}
	s.Flex = tview.NewFlex()
	s.Flex.SetDirection(tview.FlexRow)

	s.info = tview.NewTextView().
		SetDynamicColors(true).
		SetRegions(true).
		SetWrap(false).
		SetTextAlign(tview.AlignCenter).
		SetHighlightedFunc(func(added []string, _ []string, _ []string) {
			s.pages.SwitchToPage(added[0])
		})

	s.pages = tview.NewPages()
	s.pages.SetBorderPadding(0, 0, 1, 1)

	s.AddItem(s.info, 2, 0, false)
	s.AddItem(s.pages, 0, 1, true)

	return s
}

type SlideShow struct {
	*tview.Flex

	nextIdx int

	info  *tview.TextView
	pages *tview.Pages
}

func (s *SlideShow) Add(title string, page tview.Primitive) {
	idx := s.nextIdx
	s.nextIdx++

	first := idx == 0
	textIdx := fmt.Sprintf("%d", idx)

	_, _ = fmt.Fprintf(s.info, `%d ["%d"][yellow]%s[white][""]  `, idx+1, idx, title)
	s.pages.AddPage(textIdx, page, true, first)

	if first {
		s.info.Highlight(textIdx)
	}
}

func (s *SlideShow) Next() {
	slide, _ := strconv.Atoi(s.info.GetHighlights()[0])
	slide = (slide + 1) % s.nextIdx
	s.info.Highlight(strconv.Itoa(slide)).ScrollToHighlight()
}

func (s *SlideShow) Prev() {
	slide, _ := strconv.Atoi(s.info.GetHighlights()[0])
	slide = (slide - 1 + s.nextIdx) % s.nextIdx
	s.info.Highlight(strconv.Itoa(slide)).ScrollToHighlight()
}
