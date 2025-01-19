package btree

// TODO: Put item collection in binary tree
type Tree struct {
	Left  *Tree
	Value int
	Right *Tree
}
