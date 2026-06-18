package client

import (
	"bufio"
	"cafego/internal/database"
	"cafego/internal/interfaces"
	"cafego/internal/models/player"
	"io"
	"strings"
	"time"

	"cafego/internal/types/requests"
	"cafego/internal/types/responses"
	"net"

	"github.com/charmbracelet/log"
)

// Client
type Client struct {
	ClientID      int
	Conn          net.Conn
	Writer        *bufio.Writer           // Buffered write connection to the client
	Reader        *bufio.Reader           // Buffered read  connection to the client
	DB            *database.CafeDB        // Connection to the database
	Location      interfaces.CafeLocation // Players current location
	Player        *player.Player          // Player object
	ClientManager interfaces.ClientManager

	TimeoutStamp time.Time

	isDisconnecting bool

	RequestQueue  chan *requests.Request
	ResponseQueue chan responses.Response
}

func New(conn net.Conn, dbc *database.CafeDB, cm interfaces.ClientManager) *Client {
	return &Client{
		ClientID:      0,
		Conn:          conn,
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
	return c.ClientID
}

func (c *Client) SetClientID(id int) {
	c.ClientID = id
}

func (c *Client) GetIP() string {
	return strings.Split(c.Conn.RemoteAddr().String(), ":")[0]
}

func (c *Client) Start() {
	go c.receiveRequests()
	go c.sendResponses()
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

func (c *Client) receiveRequests() {
	defer close(c.RequestQueue)
	for {
		message, err := c.Reader.ReadString('\x00')
		if err != nil {
			if err == io.EOF {
				log.Infof("Client disconnected (EOF): %s", c.GetIP())
			} else if netErr, ok := err.(net.Error); ok {
				log.Infof("Network error from %s: %v", c.GetIP(), netErr)
			} else if strings.Contains(err.Error(), "use of closed network connection") {
				log.Infof("Connection closed: %s", c.GetIP())
			} else {
				log.Errorf("Read error from %s: %v", c.GetIP(), err)
			}
			c.Disconnect()
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

func (c *Client) Disconnect() error {
	c.isDisconnecting = true

	log.Infof("Client being disconnected: %s", c.GetIP())

	c.ClientManager.DisconnectClient(c.ClientID)

	c.Conn.Close()

	return nil
}

func (c *Client) GetIsDisconnecting() bool {
	return c.isDisconnecting
}
