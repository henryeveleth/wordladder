package solve

import lane "gopkg.in/oleiade/lane.v1"

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

func reachable(start string, graph map[string][]string) []string {
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

	reachableWords := make([]string, 0, len(visited))
	for word := range visited {
		reachableWords = append(reachableWords, word)
	}

	return reachableWords
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
