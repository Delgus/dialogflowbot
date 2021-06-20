package chat

import (
	"log"
	"time"

	"github.com/gorilla/websocket"
)

// Message struct for message
type Message struct {
	UserID int    `json:"user_id"`
	Body   string `json:"body"`
	Type   string `json:"type"`
}

const (
	// message types
	ConnectMessageType = "connect"
	TextMessageType    = "text"

	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 30 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10
)

// Client - chat client.
type Client struct {
	id     int
	ws     *websocket.Conn
	server *Server
	send   chan *Message
}

// NewClient create new chat client.
func NewClient(id int, ws *websocket.Conn, server *Server) *Client {
	if ws == nil {
		panic("ws cannot be nil")
	}

	if server == nil {
		panic("server cannot be nil")
	}

	return &Client{
		id:     id,
		ws:     ws,
		server: server,
		send:   make(chan *Message),
	}
}

// Listen read request via channel
func (c *Client) listenWrite() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.ws.Close()
	}()
	for {
		select {
		case msg := <-c.send:

			err := c.ws.SetWriteDeadline(time.Now().Add(writeWait))
			if err != nil {
				log.Println(err)
			}

			err = c.ws.WriteJSON(msg)
			if err != nil {
				c.server.Del(c)
				log.Println(err)
				return
			}

		case <-ticker.C:
			err := c.ws.SetWriteDeadline(time.Now().Add(writeWait))
			if err != nil {
				log.Println(err)
			}

			if err := c.ws.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
				return
			}
		}
	}
}

// Listen read request via channel
func (c *Client) listenRead() {
	// thx https://stackoverflow.com/questions/37696527/go-gorilla-websockets-on-ping-pong-fail-user-disconnct-call-function
	err := c.ws.SetReadDeadline(time.Now().Add(pongWait))
	if err != nil {
		log.Println(err)
	}

	c.ws.SetPongHandler(func(string) error {
		if err := c.ws.SetReadDeadline(time.Now().Add(pongWait)); err != nil {
			log.Println(err)
		}
		return nil
	})

	for {
		var msg Message
		err := c.ws.ReadJSON(&msg)

		if err != nil {
			c.server.Del(c)
			c.ws.Close()
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway) {
				log.Println(err)
			}
			break
		}

		msg.UserID = c.id
		c.server.broadcast <- &msg
	}
}
