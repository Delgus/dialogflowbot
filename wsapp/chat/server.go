package chat

import (
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

// Server -  chat server.
type Server struct {
	upgrader  websocket.Upgrader
	clients   map[*Client]bool
	broadcast chan *Message
	maxID     int
	sync.Mutex
}

// NewServer create new chat server.
func NewServer() *Server {
	return &Server{
		upgrader: websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
		},
		clients:   make(map[*Client]bool),
		broadcast: make(chan *Message),
	}
}

// Add added new client to connection
func (s *Server) Add(c *Client) {
	s.Lock()
	s.clients[c] = true
	s.Unlock()
}

// Del delete client connection
func (s *Server) Del(c *Client) {
	s.Lock()
	delete(s.clients, c)
	s.Unlock()
}

func Run() {
	server := NewServer()

	http.Handle("/entry", server)

	go server.Listen()

	http.Handle("/", http.FileServer(http.Dir("wsapp/web")))
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	conn, err := s.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	client := NewClient(s.newClientID(), conn, s)

	s.Add(client)

	go client.listenWrite()
	go client.listenRead()

	client.send <- &Message{Type: ConnectMessageType, Body: "connected", UserID: client.id}
}

// Listen and serve.
// It serves client connection and broadcast request.
func (s *Server) Listen() {
	for msg := range s.broadcast {
		s.Lock()
		for c := range s.clients {
			c.send <- msg
		}
		s.Unlock()
	}
}

func (s *Server) newClientID() int {
	s.Lock()
	defer s.Unlock()

	s.maxID++
	return s.maxID
}
