package graph_test

import (
	"testing"

	"go-airline-crew-rostering/graph"
)

func TestGraph(t *testing.T) {
	// simple test for the graph package
	pairGraph := new(graph.Graph)
	pairGraph.Initialization(10)
	node := new(graph.Node)
	node.Initialization()
	node2 := new(graph.Node)
	node2.Initialization()
	nodeId := 1
	nodeId2 := 2
	pairGraph.Nodes[nodeId] = node
	pairGraph.Nodes[nodeId2] = node2
	edge := new(graph.Edge)
	edge.Initialization(0, nodeId, nodeId2, 10)
	pairGraph.AddEdge(nodeId, nodeId2)
	edge2 := pairGraph.Edges[0]
	edge3 := pairGraph.Nodes[nodeId].Edges[nodeId2]
	edge4 := pairGraph.Nodes[nodeId2].Edges[nodeId]
	edge2.Position[1] = 5
	t.Log(edge2.Position[1], edge3.Position[1], edge4.Position[1])
	edge3.Position[1]++
	t.Log(edge2.Position[1], edge3.Position[1], edge4.Position[1])
	t.Log(edge.Position[1])
}
