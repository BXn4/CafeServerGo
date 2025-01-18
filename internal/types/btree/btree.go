package btree

type Tree struct {
	Left  *Tree
	Value int
	Right *Tree
}
