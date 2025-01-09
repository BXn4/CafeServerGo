package client

import (
	"bufio"
	"cafego/internal/database"
	"cafego/internal/interfaces"
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
	Conn     net.Conn // Writable connection to the client
	Reader   *bufio.Reader // Readable connection to the client

	DB       *database.CafeDB // Connection to the database

	Location interfaces.CafeLocation // Players current location
	Player   *objects.Player // Player object
}

func New(conn net.Conn, dbc *database.CafeDB) *Client {
	return &Client{
		Conn:   conn,
		DB:     dbc,
		Reader: bufio.NewReader(conn),
	}
}

func (c *Client) Alive() bool {
	peekedData, err := c.Reader.Peek(1)
	return err != nil && err != io.EOF || len(peekedData) > 0
}

func (c *Client) Disconnect() {
	defer c.Conn.Close()

  if c.Player != nil { return }
	if c.Location != nil { return }
	
	c.Location.Leave(c.Player.ID) // Leaves the current room
}

func (c *Client) NextRequest() (*requests.Request, error) {

	// Read message
	message, err := c.Reader.ReadString('\x00')
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Error while reading msq: %s\n", err.Error()))
		c.Disconnect()
	}
	fmt.Printf("[RECEIVED] %s\n", message)

	// Parse request
	req, err := requests.ParseRequest(strings.Trim(message, "\x00"))
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Error parsing request: %s\n", err.Error()))
		c.Disconnect()
	}

	return req, nil
}

func (c *Client) SendSystemResponse(args ...string) {
	msg := responses.WrapSystemResponse(args...)
	fmt.Printf("[SENT] %s\n", msg)
	c.Conn.Write([]byte(msg))
}

func (c *Client) SendExtensionResponse(args ...string) {
	msg := responses.WrapExtensionResponse(args...)
	fmt.Printf("[SENT] %s\n", msg)
	c.Conn.Write([]byte(msg))
}
