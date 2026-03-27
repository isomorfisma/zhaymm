package dag

import (
	"fmt"
)

// Graph represents table relations
type Graph struct {
	nodes map[string][]string // Key: Table name, Value: Tables that Key depends on
}

func NewGraph() *Graph {
	return &Graph{
		nodes: make(map[string][]string),
	}
}

// Registers table and its dependencies
func (g *Graph) AddNode(tableName string, dependsOn []string) {
	g.nodes[tableName] = dependsOn
}

// Topological Sort (DFS)
// Returns right order of tables, or error of Circular Dependency
func (g *Graph) Sort() ([]string, error) {
	visited := make(map[string]bool) 
	stack := make(map[string]bool)   
	var order []string              

	// DFS Recursive Function
	var dfs func(node string) error
	dfs = func(node string) error {
		// Circular Dependency Filter
		if stack[node] {
			return fmt.Errorf("fatal: terdeteksi circular dependency pada tabel '%s'", node)
		}
		// Skip if already visited
		if visited[node] {
			return nil
		}

		// Marked as being processed
		stack[node] = true

		// Search all tables needed by this table
		for _, dep := range g.nodes[node] {
			if err := dfs(dep); err != nil {
				return err
			}
		}

		// Finish
		stack[node] = false
		visited[node] = true
		order = append(order, node)

		return nil
	}

	// Run DFS for all nodes
	for node := range g.nodes {
		if !visited[node] {
			if err := dfs(node); err != nil {
				return nil, err
			}
		}
	}

	return order, nil
}