package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"sort"
	"strings"
)

// We store the graph as an adjacency map: each computer (string) maps to
// a set (map[string]bool) of all computers to which it is directly connected.
type Graph map[string]map[string]bool

func main() {
	file, err := os.Open("input.txt")
	if err != nil {
		log.Fatalln("can't open file", err)
	}
	defer file.Close()
    scanner := bufio.NewScanner(file)
    edges := make([]string, 0)
    for scanner.Scan() {
        line := strings.TrimSpace(scanner.Text())
        if line == "" {
            continue
        }
        edges = append(edges, line)
    }
    if err := scanner.Err(); err != nil {
        panic(err)
    }

    // Build adjacency structure
    adjacency := buildGraph(edges)

    // Find the largest clique via Bron–Kerbosch (with pivoting)
    largestClique := findMaximumClique(adjacency)

    // Sort it alphabetically and print as password
    sort.Strings(largestClique)
    password := strings.Join(largestClique, ",")
    fmt.Println(password)
}

// buildGraph takes a slice of edges (e.g. "ka-de") and constructs
// an undirected adjacency map.
func buildGraph(edges []string) Graph {
    g := make(Graph)
    for _, e := range edges {
        parts := strings.Split(e, "-")
        if len(parts) != 2 {
            // Just skip any malformed lines
            continue
        }
        a, b := parts[0], parts[1]

        if g[a] == nil {
            g[a] = make(map[string]bool)
        }
        if g[b] == nil {
            g[b] = make(map[string]bool)
        }
        g[a][b] = true
        g[b][a] = true
    }
    return g
}

// findMaximumClique runs Bron–Kerbosch with pivoting on the entire graph,
// returning the largest clique (maximum clique).
func findMaximumClique(g Graph) []string {
    // Convert graph keys to a slice so we can initialize our sets properly
    allVertices := make([]string, 0, len(g))
    for v := range g {
        allVertices = append(allVertices, v)
    }

    // We'll keep track of the best (largest) clique found so far.
    var bestClique []string

    // Bron–Kerbosch recursion function
    var bk func(r, p, x map[string]bool)
    bk = func(r, p, x map[string]bool) {
        if len(p) == 0 && len(x) == 0 {
            // We found a maximal clique. Check if it's the largest.
            if len(r) > len(bestClique) {
                bestClique = setToSlice(r)
            }
            return
        }

        // Pick a pivot (any vertex in P ∪ X) to reduce branching
        pivot := ""
        for u := range union(p, x) {
            pivot = u
            break
        }

        // Explore vertices in P that are not neighbors of the pivot
        nonNeighborsOfPivot := difference(p, g[pivot])
        for v := range nonNeighborsOfPivot {
            // Recurse with R ∪ {v}, P ∩ N(v), X ∩ N(v)
            bk(
                union(r, newSet(v)),
                intersection(p, g[v]),
                intersection(x, g[v]),
            )

            // Move v from P to X
            p = difference(p, newSet(v))
            x = union(x, newSet(v))
        }
    }

    // Convert all vertices into a set
    pSet := make(map[string]bool)
    for _, v := range allVertices {
        pSet[v] = true
    }
    // R and X start empty
    rSet := make(map[string]bool)
    xSet := make(map[string]bool)

    // Run the Bron–Kerbosch with pivoting
    bk(rSet, pSet, xSet)

    return bestClique
}

// ----------------------- //
//     Helper functions    //
// ----------------------- //

// setToSlice converts a set (map[string]bool) to a slice.
func setToSlice(s map[string]bool) []string {
    res := make([]string, 0, len(s))
    for v := range s {
        res = append(res, v)
    }
    return res
}

// newSet creates a set from a list of strings.
func newSet(vals ...string) map[string]bool {
    m := make(map[string]bool)
    for _, v := range vals {
        m[v] = true
    }
    return m
}

// union returns the set union of two sets A, B.
func union(a, b map[string]bool) map[string]bool {
    res := make(map[string]bool)
    for v := range a {
        res[v] = true
    }
    for v := range b {
        res[v] = true
    }
    return res
}

// intersection returns the set intersection of two sets A, B.
func intersection(a, b map[string]bool) map[string]bool {
    res := make(map[string]bool)
    for v := range a {
        if b[v] {
            res[v] = true
        }
    }
    return res
}

// difference returns the set difference A \ B.
func difference(a, b map[string]bool) map[string]bool {
    res := make(map[string]bool)
    for v := range a {
        if !b[v] {
            res[v] = true
        }
    }
    return res
}
