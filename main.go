package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"github.com/vmihailenco/msgpack"
	lane "gopkg.in/oleiade/lane.v1"
)

var allGraphs = make(map[int]map[string][]string)
var minLength = 1
var maxLength = 8

type WordsReponse struct {
	Count int      `json:"count"`
	Words []string `json:"words"`
}

type ResponsePathSuccess struct {
	Length int      `json:"length"`
	Path   []string `json:"path"`
}

type ResponseNeighborsSuccess struct {
	Neighbors []string `json:"neighbors"`
}

type ResponseError struct {
	Message string `json:"message"`
}

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

func findShortestPath(start string, end string, graph map[string][]string) []string {
	dq := lane.NewDeque()
	visited := make(map[string]bool)

	dq.Append([]string{start})

	for {
		if dq.Empty() {
			break
		} else {
			path := dq.Shift()

			if p, ok := path.([]string); ok {
				word := p[len(p)-1]
				visited[word] = true

				if word == end {
					return p
				}

				for _, adj := range graph[word] {
					if _, found := visited[adj]; !found {
						newPath := make([]string, len(p))
						copy(newPath, p)
						newPath = append(newPath, adj)
						dq.Append(newPath)
					}
				}
			}
		}
	}

	return nil
}

func findLongestPath(start string, end string, graph map[string][]string) []string {
	dq := lane.NewDeque()
	visited := make(map[string]bool)

	dq.Append([]string{start})

	for {
		if dq.Empty() {
			break
		} else {
			path := dq.Pop()

			if p, ok := path.([]string); ok {
				word := p[len(p)-1]
				visited[word] = true

				if word == end {
					return p
				}

				for _, adj := range graph[word] {
					if _, found := visited[adj]; !found {
						newPath := make([]string, len(p))
						copy(newPath, p)
						newPath = append(newPath, adj)
						dq.Append(newPath)
					}
				}
			}
		}
	}

	return nil
}

func loadGraph(wordLength int) map[string][]string {
	wordsFile := fmt.Sprintf("wordladder%d.msgpack", wordLength)
	data, _ := ioutil.ReadFile(wordsFile)

	graph := make(map[string][]string)
	msgpack.Unmarshal(data, &graph)

	return graph
}

func loadAllGraphs() {
	for i := minLength; i <= maxLength; i++ {
		allGraphs[i] = loadGraph(i)
	}
}

func getStats(graph map[string][]string) Stats {
	nodes := 0
	edges := 0

	var mostConnectedWord string
	var mostConnectedWordCount int
	var singletons []string

	for k, v := range graph {
		nodes++
		connections := len(v)
		edges += connections

		if connections > mostConnectedWordCount {
			mostConnectedWordCount = connections
			mostConnectedWord = k
		}

		if connections == 1 {
			singletons = append(singletons, k)
		}
	}

	mostConnected := MostConnected{
		Word:        mostConnectedWord,
		Connections: mostConnectedWordCount,
	}

	singletonsObj := Singletons{
		Count: len(singletons),
		Words: singletons,
	}

	result := Stats{
		Nodes:   nodes,
		Edges:   edges,
		Most:    mostConnected,
		Singles: singletonsObj,
	}
	return result
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
	router.HandleFunc("/stats/{length}", statsHandler).Methods("GET")
	router.HandleFunc("/words/{length}", wordsHandler).Methods("GET")
	http.Handle("/", router)

	log.Fatalln(http.ListenAndServe(":8080", nil))
}
