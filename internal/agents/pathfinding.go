package agents

import (
	"cafego/internal/models/cafe"
	"cafego/internal/models/simple"
	"container/heap"
	_ "fmt"
	"math"
)

type CafePoint struct {
	cafe *cafe.Cafe
	x    int
	y    int
}

func NewCafePoint(pos simple.Position, cafe *cafe.Cafe) *CafePoint {
	return &CafePoint{
		x:    pos.X,
		y:    pos.Y,
		cafe: cafe,
	}
}

func NewCafePointXY(x, y int, cafe *cafe.Cafe) *CafePoint {
	return &CafePoint{
		x:    x,
		y:    y,
		cafe: cafe,
	}
}

func (p *CafePoint) EqualsPosition(pos simple.Position) bool {
	return pos.X == p.x && pos.Y == p.y
}

func (p1 *CafePoint) Equals(p2 CafePoint) bool {
	return p1.x == p2.x && p1.y == p2.y
}

func (p CafePoint) Pos() simple.Position {
	return simple.NewPosition(p.x, p.y)
}

func (p CafePoint) inBounds() bool {
	return !(p.x <= 0 || p.y <= 0 || p.x >= p.cafe.GetSize() || p.y >= p.cafe.GetSize())
}

func (p *CafePoint) Neighbors() []*CafePoint {
	// Check if object
	for _, obj := range p.cafe.GetObjects() {
		if p.EqualsPosition(obj.GetPos()) {
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
			x:    p.x + direction[0],
			y:    p.y + direction[1],
			cafe: p.cafe,
		}
		if np.inBounds() {
			neighbors = append(neighbors, np)
		}
	}

	// Check if there are no objects at that position
	var finalNeighbors []*CafePoint
	for _, neighbor := range neighbors {
		empty := true
		for _, obj := range p.cafe.GetObjects() {

			if obj.IsCounter() || obj.IsChair() {
				continue
			}

			if neighbor.EqualsPosition(obj.GetPos()) {
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
	closeset := map[simple.Position]*node{start.Pos(): startNode}
	openset := &priorityQueue{startNode}
	heap.Init(openset)

	// While openset not empty
	for openset.Len() > 0 {
		current := heap.Pop(openset).(*node)

		// Reached the goal
		if current.pather.Pos() == end.Pos() {
			p := []*CafePoint{}
			curr := current
			for curr != nil {
				p = append(p, curr.pather)
				curr = curr.parent
			}
			return p, current.cost, true
		}

		for _, neighbor := range current.pather.Neighbors() {
			neighborNode, exists := closeset[neighbor.Pos()]

			// If neigbour not in the closeset
			if !exists {
				neighborNode = &node{pather: neighbor}
				closeset[neighbor.Pos()] = neighborNode
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
