package notice

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"sync"
	"time"

	zaplog "github.com/dokidokikoi/go-common/log/zap"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
)

type BroadcastParam struct {
	Topic string
	Msg   []byte
}

// Hub maintains the set of active clients and broadcasts messages to the
// clients.
type Hub struct {
	// Registered clients.
	// topic -> uid -> *Client
	clients map[string]map[string]*Client

	// Inbound messages from the clients.
	broadcast chan BroadcastParam

	// Register requests from the clients.
	register chan *Client

	// Unregister requests from clients.
	unregister chan *Client
}

func NewHub() *Hub {
	return &Hub{
		broadcast:  make(chan BroadcastParam),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[string]map[string]*Client),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			m, ok := h.clients[client.Topic]
			if !ok {
				m = make(map[string]*Client)
				h.clients[client.Topic] = m
			}
			m[client.Uid] = client
		case client := <-h.unregister:
			if m, ok := h.clients[client.Topic]; ok {
				// 防止多次关闭通道
				client.closeOnce.Do(func() {
					delete(m, client.Uid)
					close(client.send)
				})
			}
		case message := <-h.broadcast:
			for _, client := range h.clients[message.Topic] {
				select {
				case client.send <- message.Msg:
				default:
					h.unregister <- client
				}
			}
		}
	}
}

func (h *Hub) SendMsg(topic, uid string, data NoticeResponse) error {
	msg, err := json.Marshal(data)
	if err != nil {
		return err
	}
	if m, ok := h.clients[topic]; ok {
		if client, ok := m[uid]; ok {
			select {
			case client.send <- msg:
				return nil
			default:
				h.unregister <- client
			}
		}
	}
	return errors.New("send message failed")
}

func (h *Hub) SendBroadcast(topic string, data NoticeResponse) error {
	msg, err := json.Marshal(data)
	if err != nil {
		return err
	}
	for _, client := range h.clients[topic] {
		select {
		case client.send <- msg:
		default:
			h.unregister <- client
		}
	}
	return errors.New("send message failed")
}

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

// Client is a middleman between the websocket connection and the hub.
type Client struct {
	hub *Hub

	Topic string
	Uid   string

	// The websocket connection.
	conn *websocket.Conn

	// Buffered channel of outbound messages.
	send chan []byte

	closeOnce sync.Once
}

// readPump pumps messages from the websocket connection to the hub.
//
// The application runs readPump in a per-connection goroutine. The application
// ensures that there is at most one reader on a connection by executing all
// reads from this goroutine.
func (c *Client) readPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()
	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error {
		c.conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})
	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				zaplog.L().Error("read websocket error", zap.Error(err))
			}
			break
		}
		message = bytes.TrimSpace(bytes.Replace(message, newline, space, -1))
		if string(message) == "ping" {
			c.send <- []byte("pong")
		}
		// c.hub.broadcast <- message
	}
}

// writePump pumps messages from the hub to the websocket connection.
//
// A goroutine running writePump is started for each connection. The
// application ensures that there is at most one writer to a connection by
// executing all writes from this goroutine.
func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()
	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The hub closed the channel.
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			// Add queued chat messages to the current websocket message.
			n := len(c.send)
			for range n {
				w.Write(newline)
				w.Write(<-c.send)
			}

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}
