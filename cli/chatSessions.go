package cli

import (
	"encoding/json"
	"errors"

	"github.com/xssdoctor/gofabric/db"
)

// updates the session with the user input and the response from the LLM
func UpdateSession(name string, userInput string, llmResponse string) error {
	messageList := []map[string]string{{
		"Role":    "user",
		"Content": userInput,
	}, {
		"Role":    "system",
		"Content": llmResponse,
	},
}
	e := db.Entry{
		Name:    name,
	}
	sessionEntry, err := e.GetSessionByName()
	if err != nil {
		jsonString, err := json.Marshal(messageList)
		if err != nil {
			return errors.New("could not marshal message list")
		}
		e.Session = string(jsonString)
		e.InsertSession()
		return nil
	}
	var session []map[string]string
	err = json.Unmarshal([]byte(sessionEntry.Session), &session)
	if err != nil {
		return errors.New("could not unmarshal session")
	}
	session = append(session, messageList...)
	jsonString, err := json.Marshal(session)
	if err != nil {
		return errors.New("could not marshal session")
	}
	e.Session = string(jsonString)
	e.InsertSession()
	return nil
}

func getSession(name string) ([]map[string]string, error) {
	e := db.Entry{
		Name: name,
	}
	sessionEntry, err := e.GetSessionByName()
	if err != nil {
		return []map[string]string{}, errors.New("could not get session by name")
	}
	var session []map[string]string
	err = json.Unmarshal([]byte(sessionEntry.Session), &session)
	if err != nil {
		return nil, errors.New("could not unmarshal session")
	}
	return session, nil
}