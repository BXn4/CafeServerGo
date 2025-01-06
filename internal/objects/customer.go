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
	id         int
	avatar     Avatar
	pos        []int
	task       string // TODO: 'Task' here
	dish       int
	action     CustomerAction
	is_thirsty int
}
