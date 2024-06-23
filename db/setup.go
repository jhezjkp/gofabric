package db

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/storage/memory"
	"github.com/otiai10/copy"
	"github.com/xssdoctor/gofabric/chat"
)

// Setup is a function that sets up the configuration for the program
func Setup() error {
	config, err := GetConfiguration()
	if err != nil {
		if err.Error() == "could not get configuration from database" {
			err = InitialRun()
			if err != nil {
				return err
			}
			return nil
		
	} else {
		return errors.New("there is an error with your configuration. Delete the database and run the setup command again")
	}
}
	var openAi, anthropic, ollama, model, groq, google string
		e := Entry{}
		fmt.Println("Enter your OpenAI API key: (Leave blank if you don't have one)")
		fmt.Scanln(&openAi)
		if openAi == "" {
			openAi = config.Openai_api_key
		
		}
		e.Openai_api_key = strings.TrimRight(openAi, "\n")
		fmt.Println("Enter your Anthropic API key: (Leave blank if you don't have one)")
		fmt.Scanln(&anthropic)
		if anthropic == "" {
			anthropic = config.Anthropic_api_key
		}
		e.Anthropic_api_key = strings.TrimRight(anthropic, "\n")
		fmt.Println("Enter your Google API key: (Leave blank if you don't have one)")
		fmt.Scanln(&google)
		if google == "" {
			google = config.Google_api_key
		}
		e.Google_api_key = strings.TrimRight(google, "\n")
		fmt.Println("Enter your Groq API key: (Leave blank if you don't have one)")
		fmt.Scanln(&groq)
		if groq == "" {
			groq = config.Groq_api_key
		}
		e.Groq_api_key = strings.TrimRight(groq, "\n")
		fmt.Println("Enter your Ollama URL: (leave blank if you don't have one or if you want the default of localhost:11434)")
		fmt.Scanln(&ollama)
		if ollama == "" {
			ollama = config.Ollama_url
		}
		e.Ollama_url = strings.TrimRight(ollama, "\n")
		if e.Ollama_url == "" {
			e.Ollama_url = "http://127.0.0.1:11434"
		
		}
		fmt.Println()
		fmt.Println()
		chatInstance := chat.Chat{
			OpenAIApiKey: e.Openai_api_key,
			AnthropicApiKey: e.Anthropic_api_key,
			OllamaUrl: e.Ollama_url,
			GroqApiKey: e.Groq_api_key,
			GoogleApiKey: e.Google_api_key,
		}
		models, _ := chat.ListAllModels(chatInstance)
		for key, value := range models {
			fmt.Println(key)
			fmt.Println()
			for _, model := range value {
				fmt.Println(model)
			}
			fmt.Println()
		
		}
		fmt.Println("Enter your default model: Choose from the available options")
		fmt.Scanln(&model)
		if model == "" {
			model = config.Default_model
		}
		e.Default_model = strings.TrimRight(model, "\n")
		err = e.UpdateConfiguration()
		if err != nil {
			fmt.Println(err)
			return errors.New("could not update configuration")
		}
		return nil
}

// PopulateDB downloads patterns from the internet and populates the patterns folder
func PopulateDB() error {
	fmt.Println("Downloading patterns and Populating ~/.fabric/patterns..")
	fmt.Println()
	err := gitCloneAndCopy()
	if err != nil {
		return err
	}
	err = GetPatterns() // runs getpatters. this is defined below
	if err != nil {
		return err
	}

	return nil

}

// copies custom patterns to the updated patterns directory 
func PersistPatterns() error {
	home_folder, err := os.UserHomeDir()
	if err != nil {
		return err
	
	}
	patterns_folder := filepath.Join(home_folder + "/.config/fabric/patterns")
	currentPatterns, err := os.ReadDir(patterns_folder)
	if err != nil {
		return err
	}
	newpatterns_folder := filepath.Join(os.TempDir(),"patterns")
	newPatterns, err := os.ReadDir(newpatterns_folder)
	if err != nil {
		return err
	
	}
	for _, currentPattern := range currentPatterns {
		for _, newPattern := range newPatterns {
			if currentPattern.Name() == newPattern.Name() {
				break
			}
			copy.Copy(patterns_folder + "/" + newPattern.Name(), newpatterns_folder + "/" + newPattern.Name())
		}
	}
	return nil
}

// copies the new patterns into the config directory
func GetPatterns() (error) {

	home_dir, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	config_dir := home_dir + "/.config/fabric"
	patterns_dir := os.TempDir() + "/patterns"
	err = PersistPatterns()
	if err != nil {
		return err
	
	}
	copy.Copy(patterns_dir, config_dir + "/patterns") // copies the patterns to the config directory
	err = os.RemoveAll(patterns_dir) // removes the fabric directory
	if err != nil {
		return err
	}
	return nil

}

// checks if a pattern already exists in the directory
func DoesPatternExistAlready(name string) (bool, error) {
	entry := Entry{
		Name: name,
	}
	_, err := entry.GetPatternByName()
	if err != nil {
		return false, err
	}
	return true, nil
}


func gitCloneAndCopy() error {
	// Clones the given repository, creating the remote, the local branches
	// and fetching the objects, everything in memory:
	r, err := git.Clone(memory.NewStorage(), nil, &git.CloneOptions{
		URL: "https://github.com/danielmiessler/fabric.git",
	})
	if err != nil {
		fmt.Println(err)
		return err
	}

	// ... retrieves the branch pointed by HEAD
	ref, err := r.Head()
	if err != nil {
		fmt.Println(err)
		return err
	}

	// ... retrieves the commit history for /patterns folder
	cIter, err := r.Log(&git.LogOptions{
		From: ref.Hash(),
		PathFilter: func(path string) bool {
			return path == "patterns" || strings.HasPrefix(path, "patterns/")
		},
	})
	if err != nil {
		fmt.Println(err)
		return err
	}

	var changes []DirectoryChange
	// ... iterates over the commits
	err = cIter.ForEach(func(c *object.Commit) error {
		// Get the files changed in this commit by comparing with its parents
		parentIter := c.Parents()
		err = parentIter.ForEach(func(parent *object.Commit) error {
			patch, err := parent.Patch(c)
			if err != nil {
				fmt.Println(err)
				return err
			}

			for _, fileStat := range patch.Stats() {
				if strings.HasPrefix(fileStat.Name, "patterns/") {
					dir := filepath.Dir(fileStat.Name)
					changes = append(changes, DirectoryChange{Dir: dir, Timestamp: c.Committer.When})
				}
			}
			return nil
		})
		if err != nil {
			fmt.Println(err)
			return err
		}
		return nil
	})
	if err != nil {
		fmt.Println(err)
		return err
	}

	// Sort changes by timestamp
	sort.Slice(changes, func(i, j int) bool {
		return changes[i].Timestamp.Before(changes[j].Timestamp)
	})

	makeUniqueList(changes)

	commit, err := r.CommitObject(ref.Hash())
	if err != nil {
		fmt.Println(err)
		return err
	}
	tree, err := commit.Tree()
	if err != nil {
		fmt.Println(err)
		return err
	}
	err = tree.Files().ForEach(func(f *object.File) error {
		if strings.HasPrefix(f.Name, "patterns/") {
			// Create the local file path
			localPath := filepath.Join(os.TempDir(), f.Name)

			// Create the directories if they don't exist
			err := os.MkdirAll(filepath.Dir(localPath), os.ModePerm)
			if err != nil {
				fmt.Println(err)
				return err
			}

			// Write the file to the local filesystem
			blob, err := r.BlobObject(f.Hash)
			if err != nil {
				fmt.Println(err)
				return err
			}
			return writeBlobToFile(blob, localPath)
		}
		return nil
	})
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

func writeBlobToFile(blob *object.Blob, path string) error {
	reader, err := blob.Reader()
	if err != nil {
		return err
	}
	defer reader.Close()

	// Create the file
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	// Copy the contents of the blob to the file
	_, err = io.Copy(file, reader)
	if err != nil {
		return err
	}
	return nil
}
type DirectoryChange struct {
	Dir       string
	Timestamp time.Time
}

func makeUniqueList(changes []DirectoryChange) {
	uniqueItems := make(map[string]bool)
	for _, change := range changes {
		if strings.TrimSpace(change.Dir) != "" && !strings.Contains(change.Dir, "=>") {
			pattern := strings.ReplaceAll(change.Dir, "patterns/", "")
			pattern = strings.TrimSpace(pattern)
			uniqueItems[pattern] = true
		}
	}

	finalList := make([]string, 0, len(uniqueItems))
	for _, change := range changes {
		pattern := strings.ReplaceAll(change.Dir, "patterns/", "")
		pattern = strings.TrimSpace(pattern)
		if _, exists := uniqueItems[pattern]; exists {
			finalList = append(finalList, pattern)
			delete(uniqueItems, pattern) // Remove to avoid duplicates in the final list
		}
	}

	joined := strings.Join(finalList, "\n")
	home_dir, _ := os.UserHomeDir()
	file := filepath.Join(home_dir, ".config/fabric/unique_patterns.txt")
	os.WriteFile(file, []byte(joined), 0644)
}

// unzips the patterns zip file that is downloaded


