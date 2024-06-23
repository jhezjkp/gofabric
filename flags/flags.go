package flags

import (
	"bufio"
	"errors"
	"os"

	"github.com/jessevdk/go-flags"
)

// create flags struct. the users flags go into this, this will be passed to the chat struct in cli
type Flags struct {
    Pattern          string  `short:"p" long:"pattern" description:"Choose a pattern" default:""`
    Context          string  `short:"C" long:"context" description:"Choose a context" default:""`
    Session          string  `long:"session" description:"Choose a session" default:""`
    Setup            bool    `short:"S" long:"setup" description:"Run setup"`
    Temperature      float64 `short:"t" long:"temperature" description:"Set temperature" default:"0.7"`
    TopP             float64 `short:"T" long:"topp" description:"Set top P" default:"0.9"`
    Stream           bool    `short:"s" long:"stream" description:"Stream"`
    PresencePenalty  float64 `short:"P" long:"presencepenalty" description:"Set presence penalty" default:"0.0"`
    FrequencyPenalty float64 `short:"F" long:"frequencypenalty" description:"Set frequency penalty" default:"0.0"`
    ListPatterns     bool    `short:"l" long:"listpatterns" description:"List all patterns"`
    ListAllModels    bool    `short:"L" long:"listmodels" description:"List all available models"`
    ListAllContexts  bool    `short:"x" long:"listcontexts" description:"List all contexts"`
    ListAllSessions  bool    `short:"X" long:"listsessions" description:"List all sessions"`
    UpdatePatterns   bool    `short:"U" long:"updatepatterns" description:"Update patterns"`
    AddContext       bool `short:"A" long:"addcontext" description:"Add a context"`
    Message          string  `hidden:"true" description:"Message to send to chat"`
    Copy             bool    `short:"c" long:"copy" description:"Copy to clipboard"`
    Model            string  `short:"m" long:"model" description:"Choose model"`
    Url              string  `short:"u" long:"url" description:"Choose ollama url" default:"http://127.0.0.1:11434"`
    Output           string  `short:"o" long:"output" description:"Output to file" default:""`
    Interactive     bool    `short:"i" long:"interactive" description:"Interactive mode"`
    LatestPatterns string    `short:"n" long:"latest" description:"Number of latest patterns to list" default:"0"`
}

// Initialize flags. returns a Flags struct and an error
func Init() (Flags, error) {
    var o = Flags{}
    var message string

    parser := flags.NewParser(&o, flags.Default)
    args, err := parser.Parse()
    if err != nil {
        return Flags{}, err
    }

    info, _ := os.Stdin.Stat()
    hasStdin := (info.Mode() & os.ModeCharDevice) == 0

    // takes input from stdin if it exists, otherwise takes input from args (the last argument)
    if hasStdin {
        message, err = readStdin()
        if err != nil {
            return Flags{}, errors.New("error: could not read from stdin")
        }
    } else if len(args) > 0 {
        message = args[len(args)-1]
    } else {
        message = ""
    }

    o.Message = message
    return o, nil
}

// readStdin reads from stdin and returns the input as a string or an error
func readStdin() (string, error) {
    var input string
    scanner := bufio.NewScanner(os.Stdin)
    for scanner.Scan() {
        input += scanner.Text() + "\n"
    }
    if err := scanner.Err(); err != nil {
        return "", errors.New("error: could not read from stdin")
    }
    return input, nil
}