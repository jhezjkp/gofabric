package cli

import (
	"fmt"
	"strconv"

	"github.com/xssdoctor/gofabric/db"
	"github.com/xssdoctor/gofabric/flags"
	"github.com/xssdoctor/gofabric/interactive"
)

// Controls the cli. It takes in the flags and runs the appropriate functions
func Cli() (string, error){
	Flags, err := flags.Init() // initializes flags
	if err != nil {
		return "", err
	}
	if Flags.Setup { // if the setup flag is set, run the setup function
		err := db.Setup()
		if err != nil {
			return "", err
		}
		return "", nil
	
	}
	if Flags.UpdatePatterns {
		err := db.PopulateDB() // if the update patterns flag is set, run the update patterns function
		if err != nil {
			return "", err
		}
		return "", nil
	
	}
	if Flags.LatestPatterns != "0" {
		parsedToInt, err := strconv.Atoi(Flags.LatestPatterns)
		if err != nil {
			return "", err
		}
		err = latestPatterns(parsedToInt) // if the latest patterns flag is set, run the latest patterns function
		if err != nil {
			return "", err
		}
		return "", nil

	}
	if Flags.AddContext { // if the add context flag is set, run the add context function
		err := ContextAdd()
		if err != nil {
			return "", err
		}
		return "", nil
	
	}
	if Flags.ListPatterns { // if the list patterns flag is set, run the list all patterns function
		err = listAllPatterns()
		if err != nil {
			return "", err
		}
		return "", nil
	}
	if Flags.ListAllModels { // if the list all models flag is set, run the list all models function
		err = listAllModels()
		if err != nil {
			return "", err
		}
		return "", nil
	
	}
	if Flags.ListAllContexts { // if the list all contexts flag is set, run the list all contexts function
		err = listAllContexts()
		if err != nil {
			return "", err
		}
		return "", nil
	}
	if Flags.ListAllSessions { // if the list all sessions flag is set, run the list all sessions function
		err = ListSessions()
		if err != nil {
			return "", err
		}
		return "", nil
	}
	if Flags.Interactive {
		interactive.Interactive()
	} // if the interactive flag is set, run the interactive function
	message, err := initiateChat(Flags) // if none of the above flags are set, run the initiate chat function
	if err != nil {
		return "", err
	}
	if !Flags.Stream {
		fmt.Println(message)
	}
	if Flags.Copy {
		err := copyToClipboard(message) // if the copy flag is set, copy the message to the clipboard
		if err != nil {
			return "", err
		}
	}
	if Flags.Output != "" {
		createOutputFile(message, Flags.Output) // if the output flag is set, create an output file
	
	}
	return message, nil
}