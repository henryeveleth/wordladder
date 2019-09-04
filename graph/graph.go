package graph

import (
	"fmt"
	"io/ioutil"

	"github.com/vmihailenco/msgpack"
)

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
