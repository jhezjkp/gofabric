package interactive

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/joho/godotenv"
	"github.com/xssdoctor/gofabric/chat"
)

var (
	home_dir, _  = os.UserHomeDir()
	patterns_dir = filepath.Join(home_dir, ".config", "fabric", "patterns")
	env = filepath.Join(home_dir, ".config/fabric/.env")
	itemStyle    = lipgloss.NewStyle().PaddingLeft(4)
	selectedItemStyle = lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.Color("170"))
	highlightedItemStyle = lipgloss.NewStyle().
		PaddingLeft(2).
		Foreground(lipgloss.Color("#FFFFFF")). // White text color
		Background(lipgloss.Color("#5A5AAD")). // Purple background color
		Bold(true).
		Italic(true)
)

func getPatterns() []list.Item {
	patterns, err := os.ReadDir(patterns_dir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "could not read patterns directory: %v", err)
		os.Exit(1)
	}
	finalList := make([]list.Item, 0, len(patterns))
	for _, pattern := range patterns {
		finalList = append(finalList, item(pattern.Name()))
	}
	return finalList
}

func getModels(c chat.Chat) []list.Item {
	models, _ := chat.ListAllModels(c)
	finalList := make([]list.Item, 0, len(models))
	for _, modelList := range models {
		for _, model := range modelList {
			finalList = append(finalList, item(model))
		}
	}
	return finalList
}

func Interactive() {
    patterns := getPatterns()
	godotenv.Load(env)
    openaiAPIKey := os.Getenv("OPENAI_API_KEY")
	groqApiKey := os.Getenv("GROQ_API_KEY")
	google_api_key := os.Getenv("GOOGLE_API_KEY")
	claude_api_key := os.Getenv("CLAUDE_API_KEY")
    chat := chat.Chat{
		OpenAIApiKey: openaiAPIKey,
		GroqApiKey: groqApiKey,
		GoogleApiKey: google_api_key,
		AnthropicApiKey: claude_api_key,
		ResponseChan: make(chan string),
	}
	models := getModels(chat)
    chatModel := InitialChatModel(&chat)
    l1 := list.New(patterns, itemDelegate{}, 20, 10)
	l1.Title = "Patterns"
    l2 := list.New(models, itemDelegate{}, 20, 10)
	l2.Title = "Models"
    m := model{
        lists: []list.Model{l1, l2},
        focus: 0,
        chat:  &chatModel,
		fullscreen: false,
		patternsList: true,
		modelsList: true,
    }
    if _, err := tea.NewProgram(&m, tea.WithAltScreen()).Run(); err != nil {
        fmt.Println("Error running program:", err)
    }
}
