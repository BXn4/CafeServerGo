package client

import (
	"bufio"
	"cafego/internal/database"
	"cafego/internal/interfaces"
	"cafego/internal/models/player"
	"time"

	"github.com/charmbracelet/log"

	"cafego/internal/types/requests"
	"cafego/internal/types/responses"
	"net"
	"strings"
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

func (c *Client) Start() {
	go c.receiveRequests()
	go c.sendResponses()
	go c.autoSave()
	go c.listenToPin()
}

func (c *Client) Disconnect() error {
	if c.Player != nil {
		// id := c.Player.ID

		if c.Player.GetIsTutorialCompleted() {
			c.Player.SetLastLogin(time.Now().UTC())
			c.DB.UpdateLastLogin(c.Player.GetID(), c.Player.GetLastLogin())
			c.DB.SavePlayer(c.Player)
			if c.Location != nil {
				if c.Location.Cafe() != nil {
					c.DB.SaveCafe(c.Location.Cafe())
				}
			}
		}
	}

	c.ClientManager.DisconnectClient(c.ClientID)
	c.Conn.Close()
	log.Infof("Client disconnected: %s", c.GetIP())

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

func (c *Client) GetIP() string {
	return strings.Split(c.Conn.RemoteAddr().String(), ":")[0]
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

func (c *Client) autoSave() {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	for {
		if c == nil {
			return
		}
		select {
		case <-ticker.C:
			// Check if player exists
			if c.Player == nil {
				continue
			}

			// Save player data
			if c.Player.GetIsTutorialCompleted() {
				err := c.DB.SavePlayer(c.Player)
				if err != nil {
					log.Errorf("Failed to auto-save player data: %v", err)
				} else {
					log.Debugf("Auto-saved player %v data", c.Player.GetID())
				}

				// Save cafe
				cafe, err := c.DB.GetCafeByPlayerID(c.Player.GetID())
				if err != nil {
					log.Errorf("Failed to get cafe for auto-save: %v", err)
					continue
				}

				err = c.DB.SaveCafe(cafe)
				if err != nil {
					log.Error("Failed to auto-save cafe data: %v", err)
				} else {
					log.Debug("Auto-saved cafe %v data", cafe.GetID())
				}
			}
		}
	}
}

func (c *Client) listenToPin() {
	const timeout = 1 * time.Minute // client sends pin command in every 1 minutes
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		if c.Player == nil {
			return
		}
		if time.Since(c.TimeoutStamp) > timeout {
			log.Warnf("Client %v timed out", c.Player.GetID())
			c.Disconnect()
			return
		}
	}
}
