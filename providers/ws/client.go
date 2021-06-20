package ws

import (
	"log"
	"net/url"
	"time"

	"github.com/delgus/dialogflowbot/providers/common"
	"github.com/delgus/dialogflowbot/wsapp/chat"
	"github.com/gorilla/websocket"
)

const (
	maxConnectAttempts = 3
	wsChatID           = 1
)

type Client struct {
	conn  *websocket.Conn
	wsURL *url.URL
	myID  int
}

func NewClient(websocketURL *url.URL) *Client {
	return &Client{wsURL: websocketURL}
}

func (c *Client) connect() {
	var attempts int
	for {
		conn, _, err := websocket.DefaultDialer.Dial(c.wsURL.String(), nil)
		if err != nil {
			attempts++

			if attempts > maxConnectAttempts {
				log.Fatal("can not connect to ws", err)
			}

			time.Sleep(10 * time.Second)
			continue
		}

		c.conn = conn

		break
	}
}

func (c *Client) GetMessages() <-chan common.Message {
	ch := make(chan common.Message)

	go func() {
		var msg chat.Message

		c.connect()

		for {
			if err := c.conn.ReadJSON(&msg); err != nil {
				c.conn.Close()
				log.Println(err)
				c.connect()
				continue
			}

			if msg.Type == chat.ConnectMessageType {
				c.myID = msg.UserID
			}

			// not answer on own message
			if msg.UserID == c.myID {
				continue
			}

			ch <- common.Message{ChatID: wsChatID, Content: msg.Body}
		}
	}()

	return ch
}

func (c *Client) SendMessage(msg common.Message) error {
	return c.conn.WriteJSON(&chat.Message{Body: msg.Content})
}
