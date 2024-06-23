package interactive

import (
	"fmt"
	"io"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

type item string

func (i item) FilterValue() string { return string(i) }

type itemDelegate struct {
	highlightedIndex int
}

func (d itemDelegate) Spacing() int { return 0 }
func (d itemDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd {
	return nil
}
func (d itemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(item)
	if !ok {
		return
	}

	str := fmt.Sprint(i)

	fn := itemStyle.Render
	if index == m.Index() {
		if index == d.highlightedIndex {
			fn = func(s ...string) string {
				return highlightedItemStyle.Render(strings.Join(s, " "))
			}
		} else {
			fn = func(s ...string) string {
				return selectedItemStyle.Render(strings.Join(s, " "))
			}
		}
	}

	fmt.Fprint(w, fn(str))
}

func (d itemDelegate) Height() int {
	return 1
}