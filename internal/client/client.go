package client

import (
	"bufio"
	"cafego/internal/database"
	"cafego/internal/interfaces"
	"time"

	"github.com/charmbracelet/log"

	"cafego/internal/objects"
	"cafego/internal/types/requests"
	"cafego/internal/types/responses"
	"net"
	"strings"
)

// Client
type Client struct {
	conn          net.Conn
	Writer        *bufio.Writer           // Buffered write connection to the client
	Reader        *bufio.Reader           // Buffered read  connection to the client
	DB            *database.CafeDB        // Connection to the database
	Location      interfaces.CafeLocation // Players current location
	Player        *objects.Player         // Player object
	ClientManager interfaces.ClientManager

	TimeoutStamp time.Time

	RequestQueue  chan *requests.Request
	ResponseQueue chan responses.Response
}

func New(conn net.Conn, dbc *database.CafeDB, cm interfaces.ClientManager) *Client {
	return &Client{
		conn:          conn,
		Reader:        bufio.NewReader(conn),
		Writer:        bufio.NewWriter(conn),
		DB:            dbc,
		ClientManager: cm,

		TimeoutStamp: time.Now(),

		RequestQueue:  make(chan *requests.Request, 255),
		ResponseQueue: make(chan responses.Response, 255),
	}
}

func (c *Client) ID() int {
	return c.Player.ID
}

func (c *Client) Start() {
	go c.receiveRequests()
	go c.sendResponses()
}

func (c *Client) Disconnect() error {
	if c.Player != nil {
		id := c.Player.ID
		c.ClientManager.DisconnectClient(id)
	}
	return nil
}

func (c *Client) SendSystemResponse(args ...string) {
	resp := responses.NewSystemResponse(args...)
	log.Logf(log.Level(-3), "%s", resp.Wrap())
	c.ResponseQueue <- resp
}

func (c *Client) SendExtensionResponse(args ...string) {
	resp := responses.NewExtensionResponse(args...)
	log.Logf(log.Level(-3), "%s", resp.Wrap())
	c.ResponseQueue <- resp
}

func (c *Client) sendResponses() {
	defer close(c.ResponseQueue)
	for resp := range c.ResponseQueue {
		if resp == nil {
			return
		}
		c.Writer.Write([]byte(resp.Wrap()))
		c.Writer.Flush()
	}
}

func (c *Client) receiveRequests() {
	defer close(c.RequestQueue)
	for {

		message, err := c.Reader.ReadString('\x00')
		if err != nil {
			return
		}
		log.Logf(log.Level(-5), "%s", message)

		// Parse request
		req, err := requests.ParseRequest(strings.Trim(message, "\x00"))
		if err != nil {
			log.Error("Failed to parse request: %v", err)
			continue
		}

		c.RequestQueue <- req
	}
}
