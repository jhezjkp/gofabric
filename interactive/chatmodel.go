package interactive

import (
	"fmt"

	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/xssdoctor/gofabric/chat"
)

type chatModel struct {
    userInput    textarea.Model
    outputView   viewport.Model
    senderStyle  lipgloss.Style
    err          errMsg
    chat          *chat.Chat
    responses    string
    quitting     bool
}

type errMsg error



func InitialChatModel(chat *chat.Chat) chatModel {
	ta := textarea.New()
	ta.Placeholder = "Send a message..."
	ta.Focus()
	ta.Prompt = " "
	ta.SetWidth(30)
	ta.SetHeight(3)
	ta.FocusedStyle.CursorLine = lipgloss.NewStyle()
	ta.ShowLineNumbers = false
	vp := viewport.New(30, 5)
	vp.SetContent("Keybindings:\n\nCtrl+s: Send message or choose a pattern/model\nCtrl+p: Choose a pattern\nCtrl+b: Choose a model\n\nCtrl+f: interact with llm\nCtrl+c: quit")
    return chatModel{
        userInput:    ta,
        outputView:   vp,
        senderStyle:  lipgloss.NewStyle().Foreground(lipgloss.Color("#FF00FF")),
        err:          nil,
        chat:          chat,
        responses:    "",
        quitting:     false,
    }
}

func (m chatModel) Init() tea.Cmd {
	return textarea.Blink
}

func (m *chatModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    var (
        tiCmd tea.Cmd
        vpCmd tea.Cmd
    )

    m.userInput, tiCmd = m.userInput.Update(msg)
    m.outputView, vpCmd = m.outputView.Update(msg)

    switch msg := msg.(type) {
    case tea.KeyMsg:
        switch msg.Type {
        case tea.KeyCtrlS:
            m.chat.Message = m.userInput.Value()
            m.chat.Stream = true
            m.chat.ResponseChan = make(chan string)
            go func() {
                m.chat.SendMessageToModel()
            }()
            m.userInput.Reset()
            m.responses = ""
            return m, tea.Batch(tiCmd, vpCmd, m.waitForResponse())
        case tea.KeyCtrlC:
            m.quitting = true
            return m, tea.Quit
        }
    case string:
        m.responses += msg
        m.outputView.SetContent(m.responses)
        m.outputView.GotoBottom()
        return m, tea.Batch(tiCmd, vpCmd, m.waitForResponse())
    }

    return m, tea.Batch(tiCmd, vpCmd)
}

func (m *chatModel) waitForResponse() tea.Cmd {
    return func() tea.Msg {
        response, ok := <-m.chat.ResponseChan
        if !ok {
            return nil
        }
        return response
    }
}


func (m *chatModel) View() string {
    s := fmt.Sprintf(
        "%s\n\n%s",
        m.outputView.View(),
        m.userInput.View(),
    ) + "\n\n"
    if m.quitting {
        s += "\n"
    }
    return s
}