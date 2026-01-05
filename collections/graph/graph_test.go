package graph

import (
	"testing"
)

func TestNewGraph(t *testing.T) {
	g := New[string]()
	if g == nil {
		t.Error("New() should not return nil")
	}
	if g.Size() != 0 {
		t.Errorf("New graph size should be 0, got %d", g.Size())
	}
	if len(g.Vertices()) != 0 {
		t.Error("New graph should have no vertices")
	}
}

func TestAddVertex(t *testing.T) {
	g := New[string]()
	vertex := "A"

	g.AddVertex(vertex)
	if g.Size() != 1 {
		t.Errorf("Graph size should be 1 after adding a vertex, got %d", g.Size())
	}
	if !g.HasVertex(vertex) {
		t.Error("Graph should have the added vertex")
	}

	// Adding the same vertex again should not change the size
	g.AddVertex(vertex)
	if g.Size() != 1 {
		t.Errorf("Graph size should still be 1 after adding the same vertex again, got %d", g.Size())
	}
}

func TestRemoveVertex(t *testing.T) {
	g := New[string]()
	vertexA := "A"
	vertexB := "B"
	vertexC := "C"

	g.AddVertex(vertexA)
	g.AddVertex(vertexB)
	g.AddVertex(vertexC)
	g.AddEdge(vertexA, vertexB)
	g.AddEdge(vertexA, vertexC)

	// Ensure the vertex and its edges exist
	if !g.HasVertex(vertexA) {
		t.Error("Graph should have vertex A")
	}
	if !g.HasEdge(vertexA, vertexB) || !g.HasEdge(vertexA, vertexC) {
		t.Error("Graph should have edges from A to B and C")
	}
	if !g.HasEdge(vertexB, vertexA) || !g.HasEdge(vertexC, vertexA) {
		t.Error("Graph should have reciprocal edges in undirected graph")
	}

	g.RemoveVertex(vertexA)

	if g.HasVertex(vertexA) {
		t.Error("Graph should not have vertex A after removal")
	}
	if g.Size() != 2 {
		t.Errorf("Graph size should be 2 after removing a vertex, got %d", g.Size())
	}
	if g.HasEdge(vertexA, vertexB) || g.HasEdge(vertexA, vertexC) {
		t.Error("Graph should not have edges from A after removal")
	}
	if g.HasEdge(vertexB, vertexA) || g.HasEdge(vertexC, vertexA) {
		t.Error("Graph should not have reciprocal edges from A after removal")
	}

	// Ensure other vertices are unaffected
	if !g.HasVertex(vertexB) || !g.HasVertex(vertexC) {
		t.Error("Other vertices should remain after removing A")
	}
}

func TestAddEdge(t *testing.T) {
	g := New[string]()
	vertexA := "A"
	vertexB := "B"

	// Adding an edge between non-existent vertices should fail
	if g.AddEdge(vertexA, vertexB) {
		t.Error("AddEdge should return false if vertices do not exist")
	}

	g.AddVertex(vertexA)
	g.AddVertex(vertexB)

	// Adding an edge between existing vertices should succeed
	if !g.AddEdge(vertexA, vertexB) {
		t.Error("AddEdge should return true for valid vertices")
	}
	if !g.HasEdge(vertexA, vertexB) || !g.HasEdge(vertexB, vertexA) {
		t.Error("Graph should have the added edge in both directions for undirected graph")
	}

	// Adding the same edge again should not change anything
	if !g.AddEdge(vertexA, vertexB) { // This should still return true as the edge exists
		t.Error("AddEdge should return true if edge already exists")
	}
	if len(g.GetNeighbors(vertexA)) != 1 || len(g.GetNeighbors(vertexB)) != 1 {
		t.Error("Adding same edge twice should not duplicate neighbors")
	}
}

func TestRemoveEdge(t *testing.T) {
	g := New[string]()
	vertexA := "A"
	vertexB := "B"
	vertexC := "C"

	g.AddVertex(vertexA)
	g.AddVertex(vertexB)
	g.AddVertex(vertexC)
	g.AddEdge(vertexA, vertexB)
	g.AddEdge(vertexA, vertexC)

	// Ensure edges exist
	if !g.HasEdge(vertexA, vertexB) || !g.HasEdge(vertexA, vertexC) {
		t.Error("Graph should have the added edges")
	}

	g.RemoveEdge(vertexA, vertexB)

	if g.HasEdge(vertexA, vertexB) || g.HasEdge(vertexB, vertexA) {
		t.Error("Graph should not have the removed edge in either direction")
	}
	if !g.HasEdge(vertexA, vertexC) {
		t.Error("Other edges should remain after removing one")
	}

	// Removing a non-existent edge should not cause issues
	g.RemoveEdge(vertexA, vertexB) // Already removed
	if g.HasEdge(vertexA, vertexB) {
		t.Error("Removing non-existent edge should not re-add it")
	}
}

func TestHasVertex(t *testing.T) {
	g := New[int]()
	vertex := 1

	if g.HasVertex(vertex) {
		t.Error("HasVertex should return false for non-existent vertex")
	}

	g.AddVertex(vertex)
	if !g.HasVertex(vertex) {
		t.Error("HasVertex should return true for added vertex")
	}
}

func TestHasEdge(t *testing.T) {
	g := New[int]()
	vertex1 := 1
	vertex2 := 2

	if g.HasEdge(vertex1, vertex2) {
		t.Error("HasEdge should return false for non-existent edge")
	}

	g.AddVertex(vertex1)
	g.AddVertex(vertex2)
	g.AddEdge(vertex1, vertex2)

	if !g.HasEdge(vertex1, vertex2) || !g.HasEdge(vertex2, vertex1) {
		t.Error("HasEdge should return true for added edge in undirected graph")
	}

	g.RemoveEdge(vertex1, vertex2)
	if g.HasEdge(vertex1, vertex2) {
		t.Error("HasEdge should return false for removed edge")
	}
}

func TestGetNeighbors(t *testing.T) {
	g := New[string]()
	vertexA := "A"
	vertexB := "B"
	vertexC := "C"

	// Get neighbors of non-existent vertex should return nil
	if neighbors := g.GetNeighbors(vertexA); neighbors != nil {
		t.Error("GetNeighbors should return nil for non-existent vertex")
	}

	g.AddVertex(vertexA)
	g.AddVertex(vertexB)
	g.AddVertex(vertexC)
	g.AddEdge(vertexA, vertexB)
	g.AddEdge(vertexA, vertexC)

	neighbors := g.GetNeighbors(vertexA)
	if len(neighbors) != 2 {
		t.Errorf("Vertex A should have 2 neighbors, got %d", len(neighbors))
	}

	// Check if neighbors are correct (order may vary)
	hasB := false
	hasC := false
	for _, n := range neighbors {
		if n == vertexB {
			hasB = true
		}
		if n == vertexC {
			hasC = true
		}
	}
	if !hasB || !hasC {
		t.Error("Vertex A should have B and C as neighbors")
	}

	// Ensure modifying the returned slice doesn't affect the graph
	neighbors[0] = "DUMMY"
	neighbors2 := g.GetNeighbors(vertexA)
	if len(neighbors2) != 2 {
		t.Errorf("Modifying returned neighbors slice should not affect graph, got %d neighbors", len(neighbors2))
	}
}

func TestVertices(t *testing.T) {
	g := New[string]()
	vertexA := "A"
	vertexB := "B"
	vertexC := "C"

	vertices := g.Vertices()
	if len(vertices) != 0 {
		t.Error("Empty graph should return empty vertices slice")
	}

	g.AddVertex(vertexA)
	g.AddVertex(vertexB)
	g.AddVertex(vertexC)

	vertices = g.Vertices()
	if len(vertices) != 3 {
		t.Errorf("Graph should return 3 vertices, got %d", len(vertices))
	}

	// Check if all vertices are present (order may vary)
	hasA := false
	hasB := false
	hasC := false
	for _, v := range vertices {
		if v == vertexA {
			hasA = true
		}
		if v == vertexB {
			hasB = true
		}
		if v == vertexC {
			hasC = true
		}
	}
	if !hasA || !hasB || !hasC {
		t.Error("Vertices() should return all added vertices")
	}
}

func TestEdges(t *testing.T) {
	g := New[string]()
	vertexA := "A"
	vertexB := "B"
	vertexC := "C"

	edges := g.Edges()
	if len(edges) != 0 {
		t.Error("Empty graph should return empty edges slice")
	}

	g.AddVertex(vertexA)
	g.AddVertex(vertexB)
	g.AddVertex(vertexC)
	g.AddEdge(vertexA, vertexB)
	g.AddEdge(vertexA, vertexC)

	edges = g.Edges()
	// In an undirected graph, each edge appears twice (A-B and B-A, A-C and C-A)
	if len(edges) != 4 {
		t.Errorf("Graph should return 4 edges (2 in undirected graph), got %d", len(edges))
	}

	// Check if edges are correct
	hasAB := false
	hasBA := false
	hasAC := false
	hasCA := false
	for _, e := range edges {
		if e[0] == vertexA && e[1] == vertexB {
			hasAB = true
		}
		if e[0] == vertexB && e[1] == vertexA {
			hasBA = true
		}
		if e[0] == vertexA && e[1] == vertexC {
			hasAC = true
		}
		if e[0] == vertexC && e[1] == vertexA {
			hasCA = true
		}
	}
	if !hasAB || !hasBA || !hasAC || !hasCA {
		t.Error("Edges() should return all added edges in both directions for undirected graph")
	}
}

func TestSize(t *testing.T) {
	g := New[int]()

	if g.Size() != 0 {
		t.Errorf("New graph size should be 0, got %d", g.Size())
	}

	g.AddVertex(1)
	if g.Size() != 1 {
		t.Errorf("Graph size should be 1 after adding a vertex, got %d", g.Size())
	}

	g.AddVertex(2)
	if g.Size() != 2 {
		t.Errorf("Graph size should be 2 after adding another vertex, got %d", g.Size())
	}

	g.RemoveVertex(1)
	if g.Size() != 1 {
		t.Errorf("Graph size should be 1 after removing a vertex, got %d", g.Size())
	}

	g.RemoveVertex(1) // Remove same vertex again
	if g.Size() != 1 {
		t.Errorf("Graph size should remain 1 after removing non-existent vertex, got %d", g.Size())
	}

	g.RemoveVertex(2)
	if g.Size() != 0 {
		t.Errorf("Graph size should be 0 after removing all vertices, got %d", g.Size())
	}
}
