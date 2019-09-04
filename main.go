package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

var allGraphs = make(map[int]map[string][]string)
var minLength = 1
var maxLength = 8

type MostConnected struct {
	Word        string `json:"word"`
	Connections int    `json:"number_of_connections"`
}

type Singletons struct {
	Words []string `json:"words"`
	Count int      `json:"count"`
}

type Stats struct {
	Nodes int           `json:"nodes"`
	Edges int           `json:"edges"`
	Most  MostConnected `json:"most_connected"`

	Singles Singletons `json:"singletons"`
}

func pathHandler(w http.ResponseWriter, r *http.Request, shortest bool) {
	log.Print("STARTING PATH")

	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)

	start := strings.ToLower(vars["start"])
	end := strings.ToLower(vars["end"])
	wordLength := len(start)

	if len(end) == wordLength {
		graph := loadGraph(wordLength)

		if _, ok := graph[start]; !ok {
			respondNotFound(w, r, start)
			return
		}

		if _, ok := graph[end]; !ok {
			respondNotFound(w, r, end)
			return
		}

		var path []string
		if shortest {
			path = findShortestPath(start, end, graph)
		} else {
			path = findLongestPath(start, end, graph)
		}

		resp := ResponsePathSuccess{
			Length: len(path),
			Path:   path,
		}

		respondValid(w, r, resp)
		return
	} else {
		respondBadRequest(w, r, "Please provide two equal length words.")
		return
	}
}

func neighborsHandler(w http.ResponseWriter, r *http.Request) {
	log.Print("STARTING NEIGHBORS")

	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)

	word := strings.ToLower(vars["word"])
	wordLength := len(word)

	if wordLength >= minLength && wordLength <= maxLength {
		graph := allGraphs[wordLength]

		if neighbors, ok := graph[word]; !ok {
			respondNotFound(w, r, word)
			return
		} else {
			resp := ResponseNeighborsSuccess{
				Length:    len(neighbors),
				Neighbors: neighbors,
			}
			respondValid(w, r, resp)
			return
		}
	} else {
		msg := fmt.Sprintf("Please provide a word of length >= %d and <= %d.", minLength, maxLength)
		respondBadRequest(w, r, msg)
		return
	}
}

func reachableHandler(w http.ResponseWriter, r *http.Request) {
	log.Print("STARTING NEIGHBORS")

	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)

	word := strings.ToLower(vars["word"])
	wordLength := len(word)

	if wordLength >= minLength && wordLength <= maxLength {
		graph := allGraphs[wordLength]
		reachableWords := reachable(word, graph)
		percent := float32(len(reachableWords)) / float32(len(graph))

		resp := ResponseReachablesSuccess{
			Count:      len(reachableWords),
			Percent:    percent,
			Reachables: reachableWords,
		}
		respondValid(w, r, resp)
		return
	} else {
		msg := fmt.Sprintf("Please provide a word of length >= %d and <= %d.", minLength, maxLength)
		respondBadRequest(w, r, msg)
		return
	}
}

func statsHandler(w http.ResponseWriter, r *http.Request) {
	log.Print("STARTING STATS")

	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)

	wordLength, err := strconv.Atoi(vars["length"])
	if err != nil {
		respondBadRequest(w, r, "Please provide an integer word length.")
	}

	if wordLength >= minLength && wordLength <= maxLength {
		graph := loadGraph(wordLength)
		stats := getStats(graph)

		respondValid(w, r, stats)
		return
	} else {
		msg := fmt.Sprintf("Please provide a word of length >= %d and <= %d.", minLength, maxLength)
		respondBadRequest(w, r, msg)
		return
	}
}

func wordsHandler(w http.ResponseWriter, r *http.Request) {
	log.Print("STARTING WORDS")

	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)

	wordLength, err := strconv.Atoi(vars["length"])
	if err != nil {
		respondBadRequest(w, r, "Please provide an integer word length.")
	}

	if wordLength >= minLength && wordLength <= maxLength {
		graph := loadGraph(wordLength)

		words := make([]string, len(graph))

		i := 0
		for k := range graph {
			words[i] = k
			i++
		}

		wordsReponse := WordsReponse{
			Count: len(words),
			Words: words,
		}

		respondValid(w, r, wordsReponse)
		return
	} else {
		msg := fmt.Sprintf("Please provide a word of length >= %d and <= %d.", minLength, maxLength)
		respondBadRequest(w, r, msg)
		return
	}
}

func main() {
	log.SetFormatter(&log.JSONFormatter{})
	log.SetReportCaller(true)

	router := mux.NewRouter()
	router.StrictSlash(true)

	log.Print("INITIALIZING GRAPHS")
	loadAllGraphs()
	log.Print("DONE INITIALIZING GRAPHS")

	json.Marshal(router)
	router.HandleFunc("/path/{start}/{end}", func(w http.ResponseWriter, r *http.Request) {
		pathHandler(w, r, true)
	}).Methods("GET")

	router.HandleFunc("/longpath/{start}/{end}", func(w http.ResponseWriter, r *http.Request) {
		pathHandler(w, r, false)
	}).Methods("GET")

	router.HandleFunc("/neighbors/{word}", neighborsHandler).Methods("GET")
	router.HandleFunc("/reachables/{word}", reachableHandler).Methods("GET")
	router.HandleFunc("/stats/{length}", statsHandler).Methods("GET")
	router.HandleFunc("/words/{length}", wordsHandler).Methods("GET")
	http.Handle("/", router)

	log.Fatalln(http.ListenAndServe(":8080", nil))
}
