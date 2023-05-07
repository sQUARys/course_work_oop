package main

import (
	"fmt"
	"math"
)

// Г-42-2
// Найти максимальный поток в сети с помощью алгоритма Форда-Фалкерсона.
// Задачу решить в общем виде.
// В качестве контрольного примера использовать задание соответствующего варианта.

// TODO : добавить еще граф для примера

type Edge struct {
	// две вершины ребра
	u, v int
	// вес ребра
	weight float64
}

type Graph struct {
	edges []Edge
	// ключ - вершина, значение-массив всех ребер этой вершины
	adjList map[int][]Edge
	// количество вершин
	vertexCount int
}

func NewGraph(vertexCount int) *Graph {
	return &Graph{
		adjList:     make(map[int][]Edge),
		vertexCount: vertexCount,
	}
}

// добавляем ребро
func (g *Graph) AddEdge(u, v int, weight float64) {
	g.edges = append(g.edges, Edge{u, v, weight})
	g.adjList[u] = append(g.adjList[u], Edge{u, v, weight})
	g.adjList[v] = append(g.adjList[v], Edge{v, u, 0})
}

// поиск максимального потока
func (g *Graph) MaxFlow(source, sink int) float64 {
	var flow float64
	for {
		parent := make([]int, g.vertexCount)
		for i := range parent {
			parent[i] = -1
		}
		parent[source] = -2
		queue := []int{source}
		for len(queue) > 0 {
			u := queue[0]
			queue = queue[1:]
			for _, e := range g.adjList[u] {
				if parent[e.v] == -1 && e.weight > 0 {
					parent[e.v] = u
					queue = append(queue, e.v)
				}
			}
		}
		if parent[sink] == -1 {
			break
		}
		pathFlow := math.Inf(1)
		for v := sink; v != source; v = parent[v] {
			u := parent[v]
			for _, e := range g.adjList[u] {
				if e.v == v {
					pathFlow = math.Min(pathFlow, e.weight)
					break
				}
			}
		}
		for v := sink; v != source; v = parent[v] {
			u := parent[v]
			for i := range g.adjList[u] {
				if g.adjList[u][i].v == v {
					g.adjList[u][i].weight -= pathFlow
					break
				}
			}
			found := false
			for i := range g.adjList[v] {
				if g.adjList[v][i].v == u {
					g.adjList[v][i].weight += pathFlow
					found = true
					break
				}
			}
			if !found {
				g.adjList[v] = append(g.adjList[v], Edge{v, u, pathFlow})
			}
		}
		flow += pathFlow
	}
	return flow
}

func main() {
	// доп граф для проверки
	// graph := NewGraph(6)
	// graph.AddEdge(0, 1, 16)
	// graph.AddEdge(0, 2, 13)
	// graph.AddEdge(1, 3, 12)
	// graph.AddEdge(2, 1, 4)
	// graph.AddEdge(2, 4, 14)
	// graph.AddEdge(1, 2, 10)
	// graph.AddEdge(3, 5, 20)
	// graph.AddEdge(4, 3, 7)
	// graph.AddEdge(4, 5, 4)
	// source := 0
	// sink := 5
	// maxFlow := graph.MaxFlow(source, sink)
	// fmt.Printf("Максимальный поток между вершиной %d и %d равен %g\n", source, sink, maxFlow)

	// граф из контрольного примера
	graph := NewGraph(15)
	graph.AddEdge(1, 2, 9)
	graph.AddEdge(1, 3, 8)
	graph.AddEdge(1, 4, 9)
	graph.AddEdge(1, 5, 6)

	graph.AddEdge(2, 6, 6)
	graph.AddEdge(2, 3, 5)

	graph.AddEdge(3, 6, 2)
	graph.AddEdge(3, 7, 5)

	graph.AddEdge(4, 3, 4)
	graph.AddEdge(4, 8, 6)
	graph.AddEdge(4, 5, 3)

	graph.AddEdge(5, 9, 4)

	graph.AddEdge(6, 10, 7)

	graph.AddEdge(7, 6, 5)
	graph.AddEdge(7, 10, 3)
	graph.AddEdge(7, 11, 5)
	graph.AddEdge(7, 8, 5)

	graph.AddEdge(8, 3, 2)
	graph.AddEdge(8, 12, 7)
	graph.AddEdge(8, 9, 6)

	graph.AddEdge(9, 4, 2)
	graph.AddEdge(9, 13, 5)

	graph.AddEdge(10, 14, 9)

	graph.AddEdge(11, 10, 3)
	graph.AddEdge(11, 14, 7)

	graph.AddEdge(12, 7, 3)
	graph.AddEdge(12, 11, 1)
	graph.AddEdge(12, 14, 6)
	graph.AddEdge(12, 13, 1)

	graph.AddEdge(13, 8, 3)
	graph.AddEdge(13, 14, 9)

	source := 1
	sink := 14
	maxFlow := graph.MaxFlow(source, sink)
	fmt.Printf("Максимальный поток между вершиной %d и %d равен %g\n", source, sink, maxFlow)
}
