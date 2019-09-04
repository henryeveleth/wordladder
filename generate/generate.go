package generate

import (
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/vmihailenco/msgpack"
)

func oneEditAway(w1 string, w2 string) bool {
	edits := 0

	for i := range w1 {
		if w1[i] != w2[i] {
			edits++
		}

		if edits >= 2 {
			return false
		}
	}

	return edits == 1
}

func generateGraph(words []string) map[string][]string {
	m := make(map[string][]string)

	/* Initialize with all words, to include those without
	 * any neighbors. */
	for _, w1 := range words {
		m[w1] = []string{w1}
	}

	for _, w1 := range words {
		for _, w2 := range words {
			if w1 != w2 && oneEditAway(w1, w2) {
				m[w1] = append(m[w1], w2)
			}
		}
	}

	return m
}

func lowerWords(words []string) []string {
	var newWords []string

	for _, e := range words {
		newWords = append(newWords, strings.ToLower(e))
	}

	return newWords
}

func generate(wordsFile string) {
	content, _ := ioutil.ReadFile(wordsFile)
	words := strings.Split(string(content), "\n")
	words = words[:len(words)-1]
	wordLength := len(words[0])
	outFile := fmt.Sprintf("../wordladder%d.msgpack", wordLength)

	words = lowerWords(words)
	m := generateGraph(words)
	packed, _ := msgpack.Marshal(&m)
	ioutil.WriteFile(outFile, packed, 0666)
}
