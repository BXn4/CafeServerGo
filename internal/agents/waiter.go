package agents

import (
  "cafego/internal/interfaces"
  "cafego/internal/objects"
  "strconv"
  "strings"
  "math/rand"
)

// Spawns a waiter at the location
func SpawnWaiter(l interfaces.CafeLocation, w *objects.Waiter) {
  
  // --- Spawn waiter ----------
  println("WAITER SPAWNED: ", w.ID)
  
  // Set waiter starter position
  w.Pos = []int{
    l.Cafe().PlayerStart[0],
    l.Cafe().PlayerStart[1],
  } 
  
  // Send waiter info 
  args := []string{
    strconv.Itoa(w.ID), 
    "1", // NPC type (1: Waiter)
    strconv.Itoa(w.Priority),
    "-1", // DishID (unnecessary for waiters)
    w.Avatar.String(),
  }
  l.Broadcast("nav","-1","0", strings.Join(args, "+")) 

  // Spawn waiter 
  l.Broadcast("nac","-1","0", args[0], "0") 

  // --- Go to the nearest counter ------

  // Send waiter action
  counter := GetRandomCounter(l)
  if counter != nil {
    MoveWaiter(l,w,counter.Pos[0],counter.Pos[1], objects.MOVE_TO_COUNTER)
  }
}

// [SENT] %xt%nav%-1%0%0+1+50+-1+James+2+1002$0#1022$2#1042$6#1052$0#1061$0#1082$0%
// [SENT] %xt%nac%-1%0%0+0%
// [SENT] %xt%nac%-1%0%0+5+3+6%


// Does one iteration of the waiter tasks
func IterateWaiters(l interfaces.CafeLocation) {

  // Check task priority
  // if priority normal:
  // - Serve food
  // - Take plates
  // if priority serve food:
  // - Serve food while there are customers
  // - if there are no customers take the plates
  // if priority take plates
  // - Take plates while there are plates
  // - if there are no plates serve customers
}

func TakePlates(l interfaces.CafeLocation, waiter *objects.Waiter) error {
  // if there are dirty dishes
  // - Move to location
  // - Take the dishes
  // - Bring them back to the counter
  return nil
}

func ServeFood(l interfaces.CafeLocation, w *objects.Waiter) {
  //-------------------------------------
  // If there are customers without food
  // - Go to the counter
  // - Get food 
  // - Bring them food 
  // - Move back to the counter
  //-------------------------------------
  
  for{
    // Store customers waiting for food
    var waitingCustomers []*objects.Customer
    for _, customer := range l.Cafe().Customers {
      if customer.Action == objects.CUSTOMER_SIT_DOWN {
        waitingCustomers = append(waitingCustomers, customer)
      }
    }

    // If there are no customers return
    if len(waitingCustomers) == 0 {
      break
    }
    
    //
    for _, customer := range waitingCustomers {
      // Check if already a other waiter serving it
      if customer.AssignedWaiter == -1 { continue }

      // Go to a random counter
      // var counter 
      // for {
      //  TODO: Get to random counter
      //  if no counter return
      //  possible, MoveWaiter(l, w, x, y, objects.WAITER_MOVE_TO_COUNTER)
      //  if possible { break }
      // }

      // Get food
      // TODO: Send food pick up

      // Bring them the food (we dont need to check if we can reach it since the customer is there)
      // TODO: Go to customer

      // Move back to the counter
      // 
    }

  } 


}

// Get a random counter that has food and it is reachable from the start location
func GetRandomCounter(l interfaces.CafeLocation) *objects.CafeObject{

  var counters []*objects.CafeObject

  for _, object := range l.Cafe().Objects {

    if !object.IsCounter() { continue }

    counters = append(counters, object)

    // Check if counter with food
    if object.DishID >= 0{

      // Check if blocked
      start := &CafePoint{x: l.Cafe().PlayerStart[0], y: l.Cafe().PlayerStart[1], l: l}
      end :=  &CafePoint{x: object.Pos[0], y: object.Pos[1], l:l}
      _, _, found := Path(start, end)

      // If found path there return it
      if found {
        return object
      }
    }
  }

  // check if th
  for len(counters) != 0 {
    i := rand.Intn(len(counters))
    rc := counters[i] // random counter

    // Check if blocked
    start := &CafePoint{x: l.Cafe().PlayerStart[0], y: l.Cafe().PlayerStart[1], l: l}
    end :=  &CafePoint{x: rc.Pos[0], y: rc.Pos[1], l:l}
    _, _, found := Path(start, end)
    
    //
    if found { return rc }
    counters = append(counters[:i], counters[i+1:]...)
  }

  return nil
}

// This moves the waiter
// returns if possible and the path length
func MoveWaiter(l interfaces.CafeLocation, w *objects.Waiter, x int, y int, action objects.Action) (bool, int) {
  
  // Get length of path
  start := &CafePoint{x: w.Pos[0], y: w.Pos[1], l: l}
  end := &CafePoint{x: x, y: y, l: l}
  println("BEFORE PATH CALC")
  path, distance, found := Path(start, end)
  println("AFTER PATH CALC")
  
  if !found {
    return false, -1
  }

  println("MOVE WAITER MSG SEND")
  println("Distance: ", distance)
  // Send move msg
  waiterPos := path[1]
  w.Pos[0] = waiterPos.x
  w.Pos[1] = waiterPos.y
  args := []string{
    strconv.Itoa(w.ID),
    strconv.Itoa(int(action)),
    strconv.Itoa(x),
    strconv.Itoa(y),
  }
  l.Broadcast("nac","-1","0", strings.Join(args, "+")) 

  return found, distance
}



