package commands

import (
	"cafego/internal/client"
	"cafego/internal/managers"
	"cafego/internal/models/object"
	"cafego/internal/types/requests"
	"cafego/internal/types/responses"
	"cafego/internal/utils"
	"fmt"
	"strconv"
)

func init() {
	RegisterCommand(requests.C2S_EDITOR_ROTATE_OBJECT,
		CommandConfig{
			Name:       "RotateObject",
			Identifier: responses.S2C_EDITOR_ROTATE_OBJECT,
			MinArgs:    5,
			MaxArgs:    5,
		},
		RotateObjectValidator,
		RotateObject,
	)
}

// ero - C2S_EDITOR_ROTATE_OBJECT
func RotateObject(req *requests.Request, c *client.Client, gm *managers.GameManager) error {
	objX, _ := strconv.Atoi(req.Args[2])
	objY, _ := strconv.Atoi(req.Args[3])
	objRotation, _ := strconv.Atoi(req.Args[4])

	obj := c.Location.Cafe().GetObjectByPosXY(objX, objY)
	obj.SetRotation(object.CafeObjectRotation(objRotation))

	c.DB.UpdateObjects(c.Location.Cafe().ID, c.Location.Cafe().Objects.StringForDB())

	c.SendExtensionResponse("ero", "-1", "0", strconv.Itoa(objX), strconv.Itoa(objY), strconv.Itoa(objRotation))
	return nil
}

func RotateObjectValidator(req *requests.Request, c *client.Client, gm *managers.GameManager, cm CommandConfig) (string, ErrorCodes) {
	if len(req.Args) < cm.MinArgs {
		return fmt.Sprintf("Not enough args. NEEDED/GOT: %d/%d", cm.MinArgs, len(req.Args)), MIN_ARGS
	}

	if cm.MinArgs > 0 {
		if len(req.Args) > cm.MaxArgs {
			return fmt.Sprintf("Too much args. NEEDED/GOT: %d/%d", cm.MaxArgs, len(req.Args)), MAX_ARGS
		}
	}

	// Dont allow players to modify the packet and sending us ERO while not in editor.
	if c.Location.IsRunning() {
		return "The location is running", ERROR_EDITOR_ONLY_IN_EDITOR
	}

	if c.Location.Cafe().GetPlayerID() != c.Player.ID {
		return "Not the owner!", NOT_DECLARED
	}

	items, err := utils.MultiAtoi(req.Args[2:]...)
	if err != nil {
		return "Cant convert string to int!", CONVERT_ERROR
	}
	objX, objY, objRotation := items[0], items[1], items[2]

	if objRotation < 0 || objRotation > 3 {
		return "Invalid rotation!", ERROR_EDITOR_WATCH_OUT
	}

	if obj := c.Location.Cafe().GetObjectByPosXY(objX, objY); obj == nil {
		return "No object found at the pos!", ERROR_EDITOR_POSITION_NOT_VALID
	}

	return "Command ran without any errors.", NO_ERROR
}
