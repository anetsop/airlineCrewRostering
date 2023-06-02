package graph

import (
	"go-airline-crew-rostering/airline"
)

// Container for the functions related to the graph
type GraphRepo interface {
	Initialization() interface{}
	AddEdge()
	Populate()
}

// Edge of the graph
type Edge struct {
	Id       int
	SourceId int       // Id of source node
	GoalId   int       // Id of goal node
	Position []float64 // list of agents' positions
}

// Node of the graph
type Node struct {
	Edges map[int]*Edge // Map of edges connected to the node
}

// Graph struct
type Graph struct {
	Nodes         map[int]*Node // map of graph nodes
	Edges         []*Edge       // list of graph edges
	NumberOfNodes int
	NumberOfEdges int
	Agents        int // number of agents (used for the length of the Edges' list of positions)
}

func (edge *Edge) Initialization(id int, sourceId int, goalId int, agents int) interface{} {
	// Initialize a new edge
	edge.Id = id
	edge.SourceId = sourceId
	edge.GoalId = goalId
	edge.Position = make([]float64, agents)
	for i := range edge.Position {
		edge.Position[i] = 1
	}
	return edge
}

func (node *Node) Initialization() interface{} {
	// Initialize a new node with an empty map of edges
	node.Edges = make(map[int]*Edge)
	return node
}

func (graph *Graph) Initialization(agents int) interface{} {
	// Initialize a new, empty graph
	graph.Nodes = make(map[int]*Node)
	graph.Edges = []*Edge{}
	graph.Agents = agents
	graph.NumberOfEdges = 0
	graph.NumberOfNodes = 0
	return graph
}

func (graph *Graph) AddEdge(sourceId int, goalId int) {
	// Add an edge between the nodes "sourceId" and "goalId"
	edge := new(Edge)
	edge.Initialization(graph.NumberOfEdges, sourceId, goalId, graph.Agents)
	graph.Nodes[sourceId].Edges[goalId] = edge
	graph.Nodes[goalId].Edges[sourceId] = edge
	graph.Edges = append(graph.Edges, edge)
	graph.NumberOfEdges++
}

func (graph *Graph) Populate(pairsArray []*airline.Pair) {
	// Use a list of pairings to populate the graph
	// with nodes based on the pairings' id
	for _, pair := range pairsArray {
		node := new(Node)
		node.Initialization()
		graph.Nodes[pair.Id] = node
		graph.NumberOfNodes++
	}
}
