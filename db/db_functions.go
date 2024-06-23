package db

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/xssdoctor/gofabric/utils"
)

// Sets up the ~/.config/fabric directory and the ~/.config/fabric/.env file
func createTables() {
	usrHome, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	
	}
	fileName:= filepath.Join(usrHome, ".config/fabric/.env")
	// check if file exist
	if _, err := os.Stat(fileName); os.IsNotExist(err) {
		// create file
		fileContents := []byte(`CLAUDE_API_KEY=
OPENAI_API_KEY=
GROQ_API_KEY=
GOOGLE_API_KEY=
OLLAMA_URL=
DEFAULT_MODEL=`)
		os.WriteFile(fileName, fileContents, 0644)
	}
}

// finds all patterns in the patterns directory and enters the id, name, and pattern into a slice of Entry structs. it returns these entries or an error
func ListAllPatterns() ([]Entry, error) {
	homedir, err := os.UserHomeDir()
	if err != nil {
		return []Entry{}, errors.New("could not get user home directory")
	}
	patterns_dir := filepath.Join(homedir, ".config", "fabric", "patterns")
	patternNames, err := os.ReadDir(patterns_dir)
	var entries []Entry
	if err != nil {
		return []Entry{}, errors.New("could not read patterns directory")
	}
	for _, pattern := range patternNames {
		e := Entry{}
		e.Name = pattern.Name()
		pattern_path := filepath.Join(patterns_dir, pattern.Name())
		systemmd := filepath.Join(pattern_path, "system.md")
		// read the contents of the system.md file
		systemmd_contents, err := os.ReadFile(systemmd)
		if err != nil {
			return []Entry{}, errors.New("could not read system.md file")
		
		}
		e.Pattern = string(systemmd_contents)
		entries = append(entries, e)
	}
	return entries, nil
}

// finds all the contexts in the context directory and enters the id, name, and context into a slice of Entry structs. it returns these entries or an error
func ListAllContexts() ([]Entry, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return []Entry{}, err
	}
	context_dir := filepath.Join(homeDir, ".config", "fabric", "contexts")
	contextNames, err := os.ReadDir(context_dir)
	if err != nil {
		return []Entry{}, err
	}
	var entries []Entry
	for _, context := range contextNames {
		e := Entry{}
		e.Name = context.Name()
		context_path := filepath.Join(context_dir, context.Name())
		context_contents, err := os.ReadFile(context_path)
		if err != nil {
			return []Entry{}, err
		}
		e.Context = string(context_contents)
		entries = append(entries, e)
	}
	return entries, nil
}

// finds all sessions in the sessions directory and enters the id, name, and session into a slice of Entry structs. it returns these entries or an error
func ListAllSessions() ([]Entry, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return []Entry{}, err
	}
	sessions_dir := filepath.Join(homeDir, ".config", "fabric", "sessions")
	sessionsNames, err := os.ReadDir(sessions_dir)
	if err != nil {
		return []Entry{}, err
	}
	var entries []Entry
	for _, session := range sessionsNames {
		e := Entry{}
		e.Name = session.Name()
		session_path := filepath.Join(sessions_dir, session.Name())
		session_contents, err := os.ReadFile(session_path)
		if err != nil {
			return []Entry{}, err
		}
		e.Session = string(session_contents)
		entries = append(entries, e)
	}
	return entries, nil
}

// finds all configurations in the .env file and enters the id, name, and configuration into a slice of Entry structs. it returns these entries or an error
func GetConfiguration() (Entry, error) {
	var claudeKey, openaiKey, groqKey, googleKey, ollamaUrl string
	usrHome, err := os.UserHomeDir()
	if err != nil {
		return Entry{}, err
	}
	fileName:= filepath.Join(usrHome, ".config/fabric/.env")
	// check if file exist
	if _, err := os.Stat(fileName); os.IsNotExist(err) {
		createTables()
	}
	claudepattern := `CLAUDE_API_KEY=(.*)\n`
	openaiPattern := `OPENAI_API_KEY=(.*)\n`
	groqPattern := `GROQ_API_KEY=(.*)\n`
	googlePattern := `GOOGLE_API_KEY=(.*)\n`
	ollamaUrlPattern := `OLLAMA_URL=(.*)\n`
	defaultModel := `DEFAULT_MODEL=(.*)\n`
	claudeKey, err = utils.FindRegex(claudepattern, fileName)
	if err != nil {
		return Entry{}, err
	
	}
	openaiKey, err = utils.FindRegex(openaiPattern, fileName)
	if err != nil {
		return Entry{}, err
	}
	groqKey, err = utils.FindRegex(groqPattern, fileName)
	if err != nil {
		return Entry{}, err
	}
	googleKey, err = utils.FindRegex(googlePattern, fileName)
	if err != nil {
		return Entry{}, err
	}
	ollamaUrl, err = utils.FindRegex(ollamaUrlPattern, fileName)
	if err != nil {
		return Entry{}, err
	}
	defaultModel, err = utils.FindRegex(defaultModel, fileName)
	if err != nil {
		return Entry{}, err
	}
	en := Entry{}
	en.Openai_api_key = openaiKey
	en.Anthropic_api_key = claudeKey
	en.Ollama_url = ollamaUrl
	en.Default_model = defaultModel
	en.Groq_api_key = groqKey
	en.Google_api_key = googleKey
	return en, nil
}

// inserts a new context into the context directory
func (e *Entry) InsertContext() error {
	homedir, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	context_dir := filepath.Join(homedir, ".config", "fabric", "contexts")
	if _, err := os.Stat(context_dir); os.IsNotExist(err) {
		os.MkdirAll(context_dir, os.ModePerm)
	}
	context_path := filepath.Join(context_dir, e.Name)
	err = os.WriteFile(context_path, []byte(e.Context), 0644)
	if err != nil {
		return err
	}
	return nil
}

// identical to insertcontext. Will likely be removed in the future
func (e *Entry) UpdateContext() error {
	homedir, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	context_dir := filepath.Join(homedir, ".config", "fabric", "contexts")
	if _, err := os.Stat(context_dir); os.IsNotExist(err) {
		os.MkdirAll(context_dir, os.ModePerm)
	}
	context_path := filepath.Join(context_dir, e.Name)
	err = os.WriteFile(context_path, []byte(e.Context), 0644)
	if err != nil {
		return err
	}
	return nil
}

// inserts a new session into the session directory
func (e *Entry) InsertSession() error {
	homedir, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	sessions_dir := filepath.Join(homedir, ".config", "fabric", "sessions")
	if _, err := os.Stat(sessions_dir); os.IsNotExist(err) {
		os.MkdirAll(sessions_dir, os.ModePerm)
	}
	sessions_path := filepath.Join(sessions_dir, e.Name)
	err = os.WriteFile(sessions_path, []byte(e.Session), 0644)
	if err != nil {
		return err
	}
	return nil
}

// inserts or replaces the configuration in the .env file
func (e *Entry) InsertConfiguration() error {
	err := utils.InsertIntoConfiguration("OPENAI_API_KEY", e.Openai_api_key, createTables)
	if err != nil {
		return err
	}
	err = utils.InsertIntoConfiguration("CLAUDE_API_KEY", e.Anthropic_api_key, createTables)
	if err != nil {
		return err
	}
	err = utils.InsertIntoConfiguration("OLLAMA_URL", e.Ollama_url, createTables)
	if err != nil {
		return err
	}
	err = utils.InsertIntoConfiguration("GOOGLE_API_KEY", e.Google_api_key, createTables)
	if err != nil {
		return err
	}
	err = utils.InsertIntoConfiguration("DEFAULT_MODEL", e.Default_model, createTables)
	if err != nil {
		return err
	}
	err = utils.InsertIntoConfiguration("GROQ_API_KEY", e.Groq_api_key, createTables)
	if err != nil {
		return err
	}
	return nil
}

func (e *Entry) InsertOpenaiApiKey() error {
	err := utils.InsertIntoConfiguration("OPENAI_API_KEY", e.Openai_api_key, createTables)
	if err != nil {
		return err
	}
	return nil
	
}

func (e *Entry) InsertAnthropicApiKey() error {
	err := utils.InsertIntoConfiguration("CLAUDE_API_KEY", e.Anthropic_api_key, createTables)
	if err != nil {
		return err
	}
	return nil
}

func (e *Entry) InsertOllamaUrl() error {
	err := utils.InsertIntoConfiguration("OLLAMA_URL", e.Ollama_url, createTables)
	if err != nil {
		return err
	}
	return nil
}

func (e *Entry) InsertGoogleApiKey() error {
	err := utils.InsertIntoConfiguration("GOOGLE_API_KEY", e.Google_api_key, createTables)
	if err != nil {
		return err
	}
	return nil
}

func (e *Entry) InsertDefaultModel() error {
	err := utils.InsertIntoConfiguration("DEFAULT_MODEL", e.Default_model, createTables)
	if err != nil {
		return err
	}
	return nil
}

func (e *Entry) InsertGroqApiKey() error {
	err := utils.InsertIntoConfiguration("GROQ_API_KEY", e.Groq_api_key, createTables)
	if err != nil {
		return err
	}
	return nil
}

func (e *Entry) UpdateConfiguration() error {
	err := utils.InsertIntoConfiguration("OPENAI_API_KEY", e.Openai_api_key, createTables)
	if err != nil {
		return err
	}
	err = utils.InsertIntoConfiguration("CLAUDE_API_KEY", e.Anthropic_api_key, createTables)
	if err != nil {
		return err
	}
	err = utils.InsertIntoConfiguration("OLLAMA_URL", e.Ollama_url, createTables)
	if err != nil {
		return err
	}
	err = utils.InsertIntoConfiguration("GOOGLE_API_KEY", e.Google_api_key, createTables)
	if err != nil {
		return err
	}
	err = utils.InsertIntoConfiguration("DEFAULT_MODEL", e.Default_model, createTables)
	if err != nil {
		return err
	}
	err = utils.InsertIntoConfiguration("GROQ_API_KEY", e.Groq_api_key, createTables)
	if err != nil {
		return err
	}
	return nil
}

// finds a session by name and returns the session as an entry or an error
func (e *Entry) GetSessionByName() (Entry, error) {
	returnEntry := Entry{}
	homedir, err := os.UserHomeDir()
	if err != nil {
		return Entry{}, err
	}
	sessions_dir := filepath.Join(homedir, ".config", "fabric", "sessions")
	if _, err := os.Stat(sessions_dir); os.IsNotExist(err) {
		os.MkdirAll(sessions_dir, os.ModePerm)
	}
	sessions_file := filepath.Join(sessions_dir, e.Name)
	session, err := os.ReadFile(sessions_file)
	if err != nil {
		return Entry{}, err
	}
	returnEntry.Name = e.Name
	returnEntry.Session = string(session)
	return returnEntry, nil
}

// finds a context by name and returns the context as an entry or an error
func (e *Entry) GetContextByName() (Entry, error) {
	returnEntry := Entry{}
	homedir, err := os.UserHomeDir()
	if err != nil {
		return Entry{}, err
	}
	context_dir := filepath.Join(homedir, ".config", "fabric", "contexts")
	if _, err := os.Stat(context_dir); os.IsNotExist(err) {
		os.MkdirAll(context_dir, os.ModePerm)
	}
	context_file := filepath.Join(context_dir, e.Name)
	context, err := os.ReadFile(context_file)
	if err != nil {
		return Entry{}, err
	}
	returnEntry.Name = e.Name
	returnEntry.Context = string(context)
	return returnEntry, nil
}

// finds a pattern by name and returns the pattern as an entry or an error
func (e *Entry) GetPatternByName() (Entry, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return Entry{}, err
	}
	pattern_dir := filepath.Join(homeDir, ".config/fabric/patterns")
	pattern_path := filepath.Join(pattern_dir, e.Name, "system.md")
	pattern, err := os.ReadFile(pattern_path)
	if err != nil {
		return Entry{}, err
	}
    var en Entry
    en.Name = e.Name
	en.Pattern = string(pattern)
    return en, nil
}