package web

import "fmt"

// ConnectionHub xD
type ConnectionHub struct {
	connections map[*WSConn]bool
	broadcast   chan *WSMessage
	unregister  chan *WSConn
	register    chan *WSConn
}

// Hub xD
var Hub = ConnectionHub{
	connections: make(map[*WSConn]bool),
	broadcast:   make(chan *WSMessage),
	unregister:  make(chan *WSConn),
	register:    make(chan *WSConn),
}

func (h *ConnectionHub) run() {
	for {
		select {
		case conn := <-h.register:
			fmt.Printf("REGISTERING %#v\n", conn)
			h.connections[conn] = true
		case conn := <-h.unregister:
			if _, ok := h.connections[conn]; ok {
				delete(h.connections, conn)
				conn.disconnect()
			}
		case wsMessage := <-h.broadcast:
			message := wsMessage.Payload.ToJSON()
			for conn := range h.connections {
				// Figure out if this connection should even be sent this message
				if wsMessage.LevelRequired > 0 && (conn.user == nil /* || conn.user.Level < wsMessage.LevelRequired*/) {
					// The user did not fulfill the message Level Requirement
					fmt.Printf("Not sending %#v to %#v\n", wsMessage, conn)
					continue
				}

				if wsMessage.MessageType != MessageTypeAll && conn.messageType != MessageTypeAll && wsMessage.MessageType != conn.messageType {
					// Invalid message type
					fmt.Printf("Not sending %#v to %#v cuz message types differ\n", wsMessage, conn)
					continue
				}
				select {
				case conn.send <- message:
				default:
					// Not sure what this is for
					close(conn.send)
					delete(h.connections, conn)
				}
			}
		}
	}
}

// Broadcast some data to all connections
func (h *ConnectionHub) Broadcast(data *WSMessage) {
	h.broadcast <- data
}
