package objects

import (
  "encoding/json"
  "time"
  "strconv"
  "strings"
)

type CafeObjectKind int
//TODO: Write down types
const (
  STOVE CafeObjectKind = 0
  COUNTER = 1
  CHAIR = 2
  TABLE = 3
  VENDING = 4
  OTHER = 5
)

type CafeObjectRotation int
const (
  Up CafeObjectRotation = iota
  Left
  Down
  Right
)

type CafeObject struct {
  Kind CafeObjectKind         `json:"id"`
  Pos  []int                  `json:"pos"`
  Rotation CafeObjectRotation `json:"rotation"`

  DishID     int  `json:"dish_id,omitempty"`
	DishStatus int  `json:"dish_status,omitempty"`
	DishAmount int  `json:"dish_amount,omitempty"`


	FancyIng          bool      `json:"fancy_ing,omitempty"`
	StartedAt         *time.Time `json:"started_at,omitempty"`
	FinishesAt        *time.Time `json:"finishes_at,omitempty"`

}

func NewCafeObjectFromString(s string) (*CafeObject,error) {
  
  var cafeObj CafeObject

  if err := json.Unmarshal([]byte(s), &cafeObj); err != nil {
    return nil, err
  }

  return &cafeObj, nil
}

func (c *CafeObject) IsDoor() bool {
  return 201 <= c.Kind && c.Kind <= 207
}

func (c *CafeObject) IsStove() bool {
  return 252 <= c.Kind && c.Kind <= 259
}

func (c *CafeObject) IsCounter() bool {
  return 301 <= c.Kind && c.Kind <= 312
}

func (c *CafeObject) IsChair() bool {
  return 601 <= c.Kind && c.Kind <= 624
}

func (c *CafeObject) IsTable() bool {
  return 401 <= c.Kind && c.Kind <= 423
}

func (c *CafeObject) IsVending() bool {
  return c.Kind == 1701 
}

func (c *CafeObject) String() string{
  args := []string{
    strconv.Itoa(c.Pos[0]),
    strconv.Itoa(c.Pos[1]),
    strconv.Itoa(int(c.Kind)),
    strconv.Itoa(int(c.Rotation)),
  }

  // if STOVE
  if c.IsStove() {
      
    args = append(args, strconv.Itoa(c.DishID))
    println("DISH_ID:", c.DishID)
    if !(c.DishID == -1 || c.DishID == -2) {
      
      fancyIngStr := "0"
      if c.FancyIng {
        fancyIngStr = "1"
      }
      args = append(args, fancyIngStr)

      if c.StartedAt != nil {
        println("IDO CICA")
        currentTime := time.Now() 
        passedTime := currentTime.Second() - c.StartedAt.Second()
        remainingTime := c.FinishesAt.Second() - currentTime.Second()
        
        args = append(args, strconv.Itoa(passedTime))
        args = append(args, strconv.Itoa(remainingTime))
      }else{
        println("NEM IDO CICA")
        args = append(args,"-1")
        args = append(args,"-1")
      }
    }
  // if COUNTER
  } else if c.IsCounter() {
    args = append(args, strconv.Itoa(c.DishID))
    args = append(args, strconv.Itoa(c.DishAmount))
  // if CHAIR
  } else if c.IsChair() {
    args = append(args, strconv.Itoa(c.DishID))
    args = append(args, strconv.Itoa(c.DishStatus))
  // if VENDING
  } else if c.IsVending() {
    args = append(args, strconv.Itoa(c.DishID))
    args = append(args, strconv.Itoa(c.DishAmount))
  }

  return strings.Join(args, "+")
}

/*
    args = [str(self.pos[0]), str(self.pos[1]), str(self.id), str(self.rotation)]

        if self.type == ObjectType.STOVE:
            args.append(str(self.dish_id))

            if self.dish_id not in (-1, -2):
                args.append(str(int(self.fancy_ing)))
                if not self.started_at:
                    args.append('-1')
                    args.append('-1')
                else:
                    current_time = datetime.now(timezone.utc)
                    args.append(str(round((current_time - self.started_at).total_seconds())))
                    args.append(str(round((self.finishes_at - current_time).total_seconds())))

        elif self.type == ObjectType.COUNTER:
            args.append(str(self.dish_id))
            args.append(str(self.dish_amount))

        elif self.type == ObjectType.CHAIR:
            args.append(str(self.dish_id))
            args.append(str(self.dish_status))

        elif self.type == ObjectType.VENDING:
            args.append(str(self.fast_food_id))
            args.append(str(self.fast_food_amount))

        return '+'.join(args)


*/

