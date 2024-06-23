package interactive

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type model struct {
	lists []list.Model
	focus int
	chat *chatModel
	fullscreen bool
	patternsList bool
	modelsList bool
	}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) View() string {
	if m.fullscreen {
		m.chat.userInput.Focus()
	return lipgloss.JoinVertical(lipgloss.Left, m.chat.userInput.View(), m.chat.outputView.View())
	} else if m.patternsList && m.modelsList {
		left := m.lists[0].View()
		middle := lipgloss.JoinVertical(lipgloss.Left, m.chat.userInput.View(), m.chat.outputView.View())
		right := m.lists[1].View()

		return lipgloss.JoinHorizontal(lipgloss.Top, left, middle, right)
	} else if m.patternsList && !m.modelsList {
		left := m.lists[0].View()
		right := lipgloss.JoinVertical(lipgloss.Left, m.chat.userInput.View(), m.chat.outputView.View())

		return lipgloss.JoinHorizontal(lipgloss.Top, left, right)
	} else {
		right := m.lists[1].View()
		left := lipgloss.JoinVertical(lipgloss.Left, m.chat.userInput.View(), m.chat.outputView.View())

		return lipgloss.JoinHorizontal(lipgloss.Top, left, right)
	
	}
}

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		if m.fullscreen {
			m.chat.userInput.SetWidth(msg.Width)
			m.chat.userInput.SetHeight(msg.Height / 2)
			m.chat.outputView.Width = msg.Width
			m.chat.outputView.Height = msg.Height / 2
			return m, nil
		} else {
			if m.patternsList && m.modelsList {
				leftWidth := msg.Width / 3
				middleWidth := msg.Width / 3
				rightWidth := msg.Width - (leftWidth + middleWidth)
				m.lists[0].SetSize(leftWidth, msg.Height)
				m.lists[1].SetSize(rightWidth, msg.Height)
				m.chat.userInput.SetWidth(middleWidth)
				m.chat.userInput.SetHeight(msg.Height / 2)
				m.chat.outputView.Width = middleWidth
				m.chat.outputView.Height = msg.Height / 2
		
				// Assuming m.lists[1] is the right column that needs to be displayed
		
				return m, nil
			} else if m.patternsList && !m.modelsList {
				leftWidth := msg.Width / 3
				rightWidth := msg.Width - leftWidth

				m.lists[0].SetSize(leftWidth, msg.Height)
				m.chat.userInput.SetWidth(rightWidth)
				m.chat.userInput.SetHeight(msg.Height / 2)
				m.chat.outputView.Width = rightWidth
				m.chat.outputView.Height = msg.Height / 2

				return m, nil
			} else if m.modelsList && !m.patternsList{
				rightWidth := msg.Width / 3
				leftWidth := msg.Width - rightWidth

				m.lists[1].SetSize(rightWidth, msg.Height)
				m.chat.userInput.SetWidth(leftWidth)
				m.chat.userInput.SetHeight(msg.Height / 2)
				m.chat.outputView.Width = leftWidth
				m.chat.outputView.Height = msg.Height / 2

				return m, nil
			}
		}
	case tea.KeyMsg:
		switch msg.String() {
			case "ctrl+c":
				return m, tea.Quit
			case "ctrl+f":
				m.fullscreen = true
				m.patternsList = false
				m.modelsList = false
				m.focus = 3
				if m.fullscreen {
				m.chat.userInput.Focus()
				}
				return m, nil
			case "ctrl+p":
				m.patternsList = true
				m.fullscreen = false
				m.focus = 0
				return m, nil
			case "ctrl+b":
				m.modelsList = true
				m.fullscreen = false
				m.focus = 1
				return m, nil
			case "enter":
				if m.focus == 0 {
					m.lists[0].SetDelegate(itemDelegate{highlightedIndex: m.lists[0].Index()})
					patternName := m.lists[0].SelectedItem().(item).FilterValue()
					system_md := filepath.Join(patterns_dir, patternName, "system.md")
					fileBytes, err := os.ReadFile(system_md)
					if err != nil {
						fmt.Fprintf(os.Stderr, "could not read file: %v", err)
						os.Exit(1)
					}

					m.chat.chat.Pattern = string(fileBytes)
				} else if m.focus == 1 {
					m.lists[1].SetDelegate(itemDelegate{highlightedIndex: m.lists[1].Index()})
					m.chat.chat.Model = m.lists[1].SelectedItem().(item).FilterValue()
									}
								}
				}
	
	// Delegate to the appropriate model based on focus
	var cmd tea.Cmd
	switch m.focus {
		case 0:
			m.lists[0], cmd = m.lists[0].Update(msg)
		case 1:
			m.lists[1], cmd = m.lists[1].Update(msg)
		case 3:
			var updatedModel tea.Model
			updatedModel, cmd = m.chat.Update(msg)
			if updatedChatModel, ok := updatedModel.(*chatModel); ok {
				m.chat = updatedChatModel
			}
		}
		return m, cmd
	}