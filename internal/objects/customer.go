package objects

/*

class WaiterTask:
    CLEAN = 0
    SERVE = 1

*/

type CustomerAction int

const (
	CUSTOMER_INSERT               CustomerAction = 0
	CUSTOMER_WALK_TO_CHAIR                       = 1
	CUSTOMER_SIT_DOWN                            = 2
	CUSTOMER_EAT                                 = 3
	CUSTOMER_LEAVE                               = 4
	CUSTOMER_FAST_FOOD                           = 8
	CUSTOMER_GOTO_VENDING_MACHINE                = 9
	CUSTOMER_LEAVE_COMPLETE                      = 41
)

type Customer struct {
	ID         int
	Avatar     Avatar
	Pos        []int
	Task       string // TODO: 'Task' here
	Dish       int
	Action     CustomerAction
	IsThirsty int
  AssignedWaiter int
}
