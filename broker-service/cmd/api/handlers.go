package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
)

type RequestPayload struct {
	Auth AuthPayload `json:"auth"`
	Action string    `json:"action,omitempty"`
}

type AuthPayload struct {
	Email string `json:"email"`
	Password string `json:"password"`
}

func (app *Config) Broker(w http.ResponseWriter, r *http.Request) {
	payload := jsonResponse{
		Error:   false,
		Message: "Hit the broker",
	}

	app.writeJSON(w, http.StatusOK, payload)
}

func (app *Config) HandleSubmission(w http.ResponseWriter, r *http.Request) {
	var requestPayload RequestPayload

	err := app.readJSON(w, r, &requestPayload)
	if err != nil {
		app.errorJSON(w, err, http.StatusBadRequest)
		return
	}

	switch requestPayload.Action {
	case "auth":
		    app.authenticate(w, requestPayload.Auth)
	default:
		app.errorJSON(w, errors.New("Unknown action."), http.StatusBadRequest)
	}
}

func (app *Config) authenticate(w http.ResponseWriter, auth AuthPayload) {
	jsonData, err := json.MarshalIndent(auth, "", "\t")
	if err != nil {
        app.errorJSON(w, err, http.StatusInternalServerError)
        return
    }

	requst, err := http.NewRequest("POST", "http://authentication-service/authenticate", bytes.NewBuffer(jsonData))
	if err != nil {
        app.errorJSON(w, err, http.StatusInternalServerError)
        return
    }
	requst.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	response, err := client.Do(requst)
	if err != nil {
        app.errorJSON(w, err, http.StatusInternalServerError)
        return
    }

	defer response.Body.Close()

	if response.StatusCode == http.StatusUnauthorized {
		app.errorJSON(w, errors.New("Invalid credentials."), http.StatusUnauthorized)
		return
	} else if response.StatusCode != http.StatusAccepted {
		app.errorJSON(w, errors.New("Failed to authenticate."), http.StatusInternalServerError)
		return
	}

	var jsonFromService jsonResponse

	err = json.NewDecoder(response.Body).Decode(&jsonFromService)
	if err!= nil {
        app.errorJSON(w, err, http.StatusInternalServerError)
        return
    }

	payload := jsonResponse {
		Error:   false,
        Message: "Authenticated successfully.",
        Data:    jsonFromService.Data,
	}

	app.writeJSON(w, http.StatusAccepted, payload)
}

