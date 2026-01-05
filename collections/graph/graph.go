package graph

import "slices"

// Graph представляет неориентированный граф, используя список смежности.
// T - это тип вершин, который должен быть сравнимым.
type Graph[T comparable] struct {
	adjacencyList map[T][]T
}

// New creates and returns a new empty Graph.
func New[T comparable]() *Graph[T] {
	return &Graph[T]{
		adjacencyList: make(map[T][]T),
	}
}

// AddVertex adds a vertex to the graph.
// If the vertex already exists, it does nothing.
func (g *Graph[T]) AddVertex(vertex T) {
	if _, exists := g.adjacencyList[vertex]; !exists {
		g.adjacencyList[vertex] = make([]T, 0)
	}
}

// RemoveVertex removes a vertex and all its edges from the graph.
// If the vertex does not exist, it does nothing.
func (g *Graph[T]) RemoveVertex(vertex T) {
	if _, exists := g.adjacencyList[vertex]; !exists {
		return
	}

	// Remove the vertex from the adjacency list of all its neighbors
	for _, neighbor := range g.adjacencyList[vertex] {
		g.removeEdgeFromList(neighbor, vertex)
	}

	// Delete the vertex's own entry
	delete(g.adjacencyList, vertex)
}

// AddEdge adds an edge between two vertices.
// If either vertex does not exist, it returns false.
// If the edge already exists, it does nothing.
func (g *Graph[T]) AddEdge(vertex1, vertex2 T) bool {
	if _, exists := g.adjacencyList[vertex1]; !exists {
		return false
	}
	if _, exists := g.adjacencyList[vertex2]; !exists {
		return false
	}

	// Check if edge already exists
	if slices.Contains(g.adjacencyList[vertex1], vertex2) {
		return true // Edge already exists
	}

	// Add the edge (undirected graph: add to both vertices' lists)
	g.adjacencyList[vertex1] = append(g.adjacencyList[vertex1], vertex2)
	g.adjacencyList[vertex2] = append(g.adjacencyList[vertex2], vertex1)
	return true
}

// RemoveEdge removes an edge between two vertices.
// If either vertex does not exist or the edge does not exist, it does nothing.
func (g *Graph[T]) RemoveEdge(vertex1, vertex2 T) {
	if _, exists := g.adjacencyList[vertex1]; !exists {
		return
	}
	if _, exists := g.adjacencyList[vertex2]; !exists {
		return
	}

	g.removeEdgeFromList(vertex1, vertex2)
	g.removeEdgeFromList(vertex2, vertex1)
}

// removeEdgeFromList is a helper function to remove an edge from a vertex's adjacency list.
func (g *Graph[T]) removeEdgeFromList(vertex, neighbor T) {
	neighbors := g.adjacencyList[vertex]
	newNeighbors := make([]T, 0, len(neighbors))
	for _, n := range neighbors {
		if n != neighbor {
			newNeighbors = append(newNeighbors, n)
		}
	}
	g.adjacencyList[vertex] = newNeighbors
}

// HasVertex returns true if the vertex exists in the graph, false otherwise.
func (g *Graph[T]) HasVertex(vertex T) bool {
	_, exists := g.adjacencyList[vertex]
	return exists
}

// HasEdge returns true if there is an edge between the two vertices, false otherwise.
func (g *Graph[T]) HasEdge(vertex1, vertex2 T) bool {
	if _, exists := g.adjacencyList[vertex1]; !exists {
		return false
	}
	if _, exists := g.adjacencyList[vertex2]; !exists {
		return false
	}

	return slices.Contains(g.adjacencyList[vertex1], vertex2)
}

// GetNeighbors returns a slice of neighbors for the given vertex.
// If the vertex does not exist, it returns nil.
func (g *Graph[T]) GetNeighbors(vertex T) []T {
	if neighbors, exists := g.adjacencyList[vertex]; exists {
		// Return a copy to prevent external modification
		result := make([]T, len(neighbors))
		copy(result, neighbors)
		return result
	}
	return nil
}

// Vertices returns a slice of all vertices in the graph.
func (g *Graph[T]) Vertices() []T {
	vertices := make([]T, 0, len(g.adjacencyList))
	for vertex := range g.adjacencyList {
		vertices = append(vertices, vertex)
	}
	return vertices
}

// Edges returns a slice of all edges in the graph.
// Each edge is represented as a pair of vertices [vertex1, vertex2].
// In an undirected graph, each edge will appear twice: [A,B] and [B,A].
func (g *Graph[T]) Edges() [][2]T {
	edges := make([][2]T, 0)
	for vertex, neighbors := range g.adjacencyList {
		for _, neighbor := range neighbors {
			edges = append(edges, [2]T{vertex, neighbor})
		}
	}
	return edges
}

// Size returns the number of vertices in the graph.
func (g *Graph[T]) Size() int {
	return len(g.adjacencyList)
}
