package agents

import (
	"cafego/internal/interfaces"
	"container/heap"
	_ "fmt"
	"math"
)

type cafeKey struct {
	x int
	y int
}

type CafePoint struct {
	l interfaces.CafeLocation
	x int
	y int
}

func NewCafePoint(pos []int, l interfaces.CafeLocation) *CafePoint {
	return &CafePoint{
		x: pos[0],
		y: pos[1],
		l: l,
	}
}
func (p CafePoint) Key() cafeKey {
	return cafeKey{x: p.x, y: p.y}
}

func (p CafePoint) inBounds() bool {
	return !(p.x <= 0 || p.y <= 0 || p.x >= p.l.Cafe().Size || p.y >= p.l.Cafe().Size)
}

func (p *CafePoint) Neighbors() []*CafePoint {
	// Check if object
	for _, object := range p.l.Cafe().Objects {
		if object.Pos[0] == p.x && object.Pos[1] == p.y {
			return []*CafePoint{}
		}
	}

	// Gather neigbors
	var neighbors []*CafePoint

	directions := [][]int{
		{0, 1},
		{0, -1},
		{-1, 0},
		{1, 0},
	}

	// Adds neighbours that are inside bounds
	for _, direction := range directions {
		np := &CafePoint{
			x: p.x + direction[0],
			y: p.y + direction[1],
			l: p.l,
		}
		if np.inBounds() {
			neighbors = append(neighbors, np)
		}
	}

	// Check if there are no objects at that position
	var finalNeighbors []*CafePoint
	for _, neighbor := range neighbors {
		//println("neighbor: ", neighbor.x, neighbor.y)
		empty := true
		for _, object := range p.l.Cafe().Objects {

			if object.IsCounter() || object.IsChair() {
				continue
			}

			if object.Pos[0] == neighbor.x && object.Pos[1] == neighbor.y {
				empty = false
				break
			}
		}
		if empty {
			finalNeighbors = append(finalNeighbors, neighbor)
		}
	}
	return finalNeighbors
}

// Manhattan distance
func (p1 *CafePoint) EstimateCost(p2 *CafePoint) int {
	dx := math.Abs(float64(p1.x - p2.x))
	dy := math.Abs(float64(p1.y - p2.y))
	return int(dx + dy)
}

// node is a wrapper to store A* data for a Pather node.
type node struct {
	pather *CafePoint
	cost   int
	rank   int
	parent *node
	index  int
}

// Implementation of the a* algorithm
func Path(start, end *CafePoint) (path []*CafePoint, distance int, found bool) {
	// Init tables
	startNode := &node{pather: start}
	closeset := map[cafeKey]*node{start.Key(): startNode}
	openset := &priorityQueue{startNode}
	heap.Init(openset)

	// While openset not empty
	for openset.Len() > 0 {
		current := heap.Pop(openset).(*node)

		// Reached the goal
		if current.pather.Key() == end.Key() {
			p := []*CafePoint{}
			curr := current
			for curr != nil {
				p = append(p, curr.pather)
				curr = curr.parent
			}
			return p, current.cost, true
		}

		for _, neighbor := range current.pather.Neighbors() {
			neighborNode, exists := closeset[neighbor.Key()]

			// If neigbour not in the closeset
			if !exists {
				neighborNode = &node{pather: neighbor}
				closeset[neighbor.Key()] = neighborNode
			}
			//
			cost := current.cost + 1
			if !exists || cost < neighborNode.cost {
				neighborNode.cost = cost
				neighborNode.rank = cost + neighbor.EstimateCost(end)
				neighborNode.parent = current
				if exists {
					heap.Fix(openset, neighborNode.index)
				} else {
					heap.Push(openset, neighborNode)
				}
			}

		}
	}
	return
}
