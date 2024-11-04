package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/rs/zerolog/log"
)

// type RequestPayload struct {
// 	Action string      `json:"action"`
// 	Auth   AuthPayload `json:"auth,omitempty"`
// 	Log    LogEntry    `json:"log,omitempty"`
// }

// type AuthPayload struct {
// 	Email    string `json:"email"`
// 	Password string `json:"password"`
// }

type LogEntry struct {
	Name string `json:"name"`
	Data string `json:"data"`
}

// Mailer api test handler
func (s *Server) Mailer(w http.ResponseWriter, r *http.Request) {
	payload := jsonResponse{
		Error:   false,
		Massage: "Hello from Mailer!",
	}
	writeLog(w, LogEntry{
		Name: "Test send log from mailer",
		Data: fmt.Sprintf("Mailer was sending a test log message at: %s", time.Now().String()),
	}, s.Conf.LogService)
	_ = writeJSON(w, http.StatusAccepted, payload)
}

func writeLog(w http.ResponseWriter, logs LogEntry, logService string) {
	log.Debug().Msg("post log into logger service")
	jsonData, _ := json.MarshalIndent(logs, "", "\t")
	logURL := fmt.Sprintf("%s/log", logService)
	log.Debug().Msgf("logURL: %s", logURL)
	request, err := http.NewRequest("POST", logURL, bytes.NewBuffer(jsonData))
	if err != nil {
		errorJSON(w, err)
		return
	}

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		errorJSON(w, err)
		return
	}
	defer response.Body.Close()

	// if response.StatusCode == http.StatusUnauthorized {
	// 	errorJSON(w, errors.New("invalid credentials"))
	// 	return
	// }
	if response.StatusCode != http.StatusAccepted {
		errorJSON(w, errors.New("error calling logger service"))
		return
	}

	var jsonFromService jsonResponse
	maxBytes := 1048576 // 1 Mb

	response.Body = http.MaxBytesReader(w, response.Body, int64(maxBytes))
	dec := json.NewDecoder(response.Body)
	err = dec.Decode(&jsonFromService)
	log.Debug().Msgf("jsonFromService: %+v", jsonFromService)
	if err != nil {
		errorJSON(w, errors.New(err.Error()))
		log.Error().Err(err)
		return
	}

	if jsonFromService.Error {
		errorJSON(w, errors.New(jsonFromService.Massage))
		return
	}

	payload := &jsonResponse{
		Error:   false,
		Massage: "logged!",
		Data:    jsonFromService.Data,
	}

	writeJSON(w, http.StatusAccepted, payload)
}
