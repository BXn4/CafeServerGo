package daos

type WaiterDAO struct {
	ID       int    `json:"waiter_id"`
	Name     string `json:"name"`
	Priority int    `json:"priority"`
	Avatar   string `json:"avatar"`
}
