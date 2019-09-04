package respond

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type ReponseWords struct {
	Count int      `json:"count"`
	Words []string `json:"words"`
}

type ResponsePathSuccess struct {
	Length int      `json:"length"`
	Path   []string `json:"path"`
}

type ResponseNeighborsSuccess struct {
	Length    int      `json:"length"`
	Neighbors []string `json:"neighbors"`
}

type ResponseReachablesSuccess struct {
	Count      int      `json:"count"`
	Reachables []string `json:"reachables"`
	Percent    float32  `json:"percent_of_graph"`
}

type ResponseError struct {
	Message string `json:"message"`
}

func respondNotFound(w http.ResponseWriter, r *http.Request, word string) {
	w.WriteHeader(http.StatusNotFound)
	message := fmt.Sprintf("Starting word <%s> not found in dictionary.", word)

	resp := ResponseError{
		Message: message,
	}

	js, err := json.Marshal(resp)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Write(js)
}

func respondBadRequest(w http.ResponseWriter, r *http.Request, message string) {
	w.WriteHeader(http.StatusBadRequest)

	resp := ResponseError{
		Message: message,
	}

	js, err := json.Marshal(resp)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Write(js)
}

func respondValid(w http.ResponseWriter, r *http.Request, resp interface{}) {
	w.WriteHeader(http.StatusOK)

	js, err := json.Marshal(resp)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Write(js)
}
