package commands

import(
  "cafego/internal/objects"
  "cafego/internal/types/requests"
  "cafego/internal/managers"
  "cafego/internal/utils"
  "cafego/internal/client"
  "strconv"
)


func StoveDeliverInfo(req *requests.Request, c *client.Client, gm *managers.GameManager) error {

  stoveX, err := strconv.Atoi(req.Args[2])
  if err != nil {
    return err
  }
  stoveY, err := strconv.Atoi(req.Args[3])
  if err != nil {
    return err
  }
  var stove *objects.CafeObject
  for _, object := range c.Location.Cafe().Objects {
    if object.Pos[0] == stoveX && object.Pos[1] == stoveY{
      stove = object
      break
    }
  }

  // Choose counter that is empty or has the same food type
  var counter *objects.CafeObject
  for _, object := range c.Location.Cafe().Objects {
    if !object.IsCounter() { continue }
    
    if stove.DishID == object.DishID {
      counter = object
      break
    }else if object.DishID == -1 {
      counter = object
    }
  }

  // Set args
  counterX := utils.If(counter != nil, counter.Pos[0], -1)
  counterY := utils.If(counter != nil, counter.Pos[1], -1)
  status := utils.If(counter != nil, "0", "37")

  args := []string {
    "csi", "-1", "0",
    status,
    req.Args[2],
    req.Args[3],
    strconv.Itoa(counterX),
    strconv.Itoa(counterY),
  }


  c.SendExtensionResponse(args...)
  return nil
}

