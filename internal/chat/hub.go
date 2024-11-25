package chat

import "sync"

// Hub gestiona las conexiones y difunde los mensajes.
type Hub struct {
	Rooms      map[string]map[*Client]bool
	Broadcast  chan *Message
	Register   chan *Client
	Unregister chan *Client
	mu         sync.Mutex
}

// NewHub crea una nueva instancia de Hub.
func NewHub() *Hub {
	return &Hub{
		Rooms:      make(map[string]map[*Client]bool),
		Broadcast:  make(chan *Message),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
	}
}

// Run inicia el bucle principal del hub.
func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Register:

			h.mu.Lock()
			if _, ok := h.Rooms[client.Room]; !ok {
				h.Rooms[client.Room] = make(map[*Client]bool)
			}
			h.Rooms[client.Room][client] = true
			h.mu.Unlock()

		case client := <-h.Unregister:

			h.mu.Lock()
			if clients, ok := h.Rooms[client.Room]; ok {
				if _, exists := clients[client]; exists {
					delete(clients, client)
					close(client.Send)
					if len(clients) == 0 {
						delete(h.Rooms, client.Room)
					}
				}
			}
			h.mu.Unlock()

		case message := <-h.Broadcast:
			h.mu.Lock()
			if clients, ok := h.Rooms[message.Room]; ok {
				for client := range clients {
					select {
					case client.Send <- message.Serialize():
					default:
						close(client.Send)
						delete(clients, client)
					}
				}
			}

			h.mu.Unlock()
		}
	}
}
