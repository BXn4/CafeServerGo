package interfaces

type ManagedItem interface {
	// Returns the id
	ID() int
}

type ClientManager interface {
	// Add managed item
	AddClient(item ManagedItem)

	// Get managed item by id
	GetClient(id int) (ManagedItem, error)

	// Remove managed item by id
	DisconnectClient(id int)
}

type CafeManager interface {
	// Add managed item
	AddCafe(item ManagedItem)

	// Get managed item by id
	GetCafe(id int) (ManagedItem, error)

	// Remove managed item by id
	DisconnectCafe(id int)
}
