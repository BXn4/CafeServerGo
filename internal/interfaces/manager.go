package interfaces

type ManagedItem interface {
  // Returns the id
  ID() int
}

type Manager interface {
  // Add managed item
  Add(itme *ManagedItem)

  // Get managed item by id
  Get(id int) (*ManagedItem, error)

  // Remove managed item by id
  Remove(id int)
}
