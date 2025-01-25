package client

import (
	"bufio"
	"cafego/internal/database"
	"cafego/internal/interfaces"

	"github.com/charmbracelet/log"

	"cafego/internal/objects"
	"cafego/internal/types/requests"
	"cafego/internal/types/responses"
	"errors"
	"fmt"
	"io"
	"net"
	"strings"
)

// Client
type Client struct {
	Writer        *bufio.Writer           // Buffered write connection to the client
	Reader        *bufio.Reader           // Buffered read  connection to the client
	DB            *database.CafeDB        // Connection to the database
	Location      interfaces.CafeLocation // Players current location
	Player        *objects.Player         // Player object
	ClientManager interfaces.ClientManager
}

func New(conn net.Conn, dbc *database.CafeDB, cm interfaces.ClientManager) *Client {
	return &Client{
		Reader:        bufio.NewReader(conn),
		Writer:        bufio.NewWriter(conn),
		DB:            dbc,
		ClientManager: cm,
	}
}

func (c *Client) ID() int {
	return c.Player.ID
}

func (c *Client) Alive() bool {
	peekedData, err := c.Reader.Peek(1)
	return err != nil && err != io.EOF || len(peekedData) > 0
}

func (c *Client) Disconnect() {
	if c.Player != nil {
		id := c.Player.ID
		c.ClientManager.DisconnectClient(id)
	}
}

func (c *Client) NextRequest() (*requests.Request, error) {

	// Read message
	message, err := c.Reader.ReadString('\x00')
	if err != nil {
		c.Disconnect()
		return nil, errors.New(fmt.Sprintf("Error while reading msq: %s", err.Error()))
	}
	log.Logf(log.Level(-5), "%s", message)

	// Parse request
	req, err := requests.ParseRequest(strings.Trim(message, "\x00"))
	if err != nil {
		c.Disconnect()
		return nil, errors.New(fmt.Sprintf("Error parsing request: %s", err.Error()))
	}

	return req, nil
}

func (c *Client) SendSystemResponse(args ...string) {
	msg := responses.WrapSystemResponse(args...)
	log.Logf(log.Level(-3), "%s", msg)
	c.Writer.Write([]byte(msg))
	c.Writer.Flush()
}

func (c *Client) SendExtensionResponse(args ...string) {
	msg := responses.WrapExtensionResponse(args...)
	log.Logf(log.Level(-3), "%s", msg)
	c.Writer.Write([]byte(msg))
	c.Writer.Flush()
}
