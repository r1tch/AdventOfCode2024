package main

import (
	"bufio"
	"log"
	"os"
	"sort"
	"strings"
)

// this is dead slow, but honest work :)

type computers struct {
	connections map[string]map[string]bool
	computers   map[string]bool
}

func (c *computers) addConnection(a string, b string) {
	if _, exists := c.connections[a]; !exists {
		c.connections[a] = map[string]bool{}
	}
	c.connections[a][b] = true

	if _, exists := c.connections[b]; !exists {
		c.connections[b] = map[string]bool{}
	}
	c.connections[b][a] = true
}

func (c *computers) addComputer(a string) {
	c.computers[a] = true
}

func readConnections(filename string) computers {
	file, err := os.Open(filename)
	if err != nil {
		log.Fatalln("can't open file", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	computers := computers{make(map[string]map[string]bool), make(map[string]bool)}
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Split(line, "-")
		computers.addConnection(parts[0], parts[1])
		computers.addComputer(parts[0])
		computers.addComputer(parts[1])
	}

	return computers
}

type clusters struct {
	clusters map[string]bool
}

func makeClusterKey(members []string) string {
	sort.Strings(members)
	return strings.Join(members, ",")
}

func (clusters *clusters) addConnection(computers computers, existing string, newComer string) {
	newClusters := []string{makeClusterKey([]string{existing, newComer})}
	// check all clusters...
	for cluster, _ := range clusters.clusters {
		clusterParts := strings.Split(cluster, ",")
		toAdd := true
		// if all parts of the cluster are connected to the newComer, add the newComer to the cluster
		for _, part := range clusterParts {
			if _, exists := computers.connections[part][newComer]; !exists {
				toAdd = false
				break
			}
		}
		if toAdd {
			newClusters = append(newClusters, makeClusterKey(append(clusterParts, newComer)))
		}
	}
	for _, newCluster := range newClusters {
		clusters.clusters[newCluster] = true
	}
}

func main() {
	computers := readConnections("input.txt")

	sum := 0
	for c1 := range computers.computers {
		for c2, _ := range computers.connections[c1] {
			for c3, _ := range computers.connections[c2] {
				if _, exists := computers.connections[c3][c1]; exists && c1 < c2 && c2 < c3 && (c1[0] == 't' || c2[0] == 't' || c3[0] == 't') {
					// log.Println(c1, c2, c3)
					sum++
				}
			}
		}
	}
	log.Println("sum:", sum)

	// idea: we need to find the largest clusters
	// store the clusters of interconnected computers
	clusters := clusters{map[string]bool{}}
	i := 0
	for c1 := range computers.computers {
		log.Println(i, "/", len(computers.computers), "computing connections for", c1, "clusters found so far:", len(clusters.clusters))
		for connected := range computers.connections[c1] {
			clusters.addConnection(computers, c1, connected)
		}
		i++
	}

	maxLen := 0
	maxCluster := ""
	for cluster, _ := range clusters.clusters {
		if len(cluster) > maxLen {
			maxLen = len(cluster)
			maxCluster = cluster
		}
	}

	log.Println(maxCluster)

}
