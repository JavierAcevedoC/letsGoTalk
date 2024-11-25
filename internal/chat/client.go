package chat

import (
	"bytes"
	"log"
	"time"

	"github.com/gorilla/websocket"
)

const (
	writeWait  = 10 * time.Second
	pongWait   = 60 * time.Second
	pingPeriod = (pongWait * 9) / 10
)

type Client struct {
	Hub  *Hub
	Conn *websocket.Conn
	Send chan []byte
}

func (client *Client) ReadPump() {
	defer func() {
		client.Hub.Unregister <- client
		client.Conn.Close()
	}()
	client.Conn.SetReadDeadline(time.Now().Add(pongWait))
	client.Conn.SetPongHandler(func(string) error {
		client.Conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})
	for {
		_, rawMessage, err := client.Conn.ReadMessage()
		if err != nil {
			log.Println("read error:", err)
			break
		}
		rawMessage = bytes.TrimSpace(rawMessage)

		// Parsear el mensaje recibido
		msg, err := ParseMessage(rawMessage)
		if err != nil {
			log.Println("error parsing message:", err)
			continue
		}

		log.Printf("New message from %s: %s", msg.Username, msg.Content)

		// Opcional: Reenviar a todos en la misma sala (si usas salas)
		// Aquí simplificamos y difundimos el mensaje a todos los clientes.
		client.Hub.Broadcast <- rawMessage
	}
}

func (client *Client) WritePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		client.Conn.Close()
	}()
	for {
		select {
		case message, ok := <-client.Send:
			client.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				client.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			w, err := client.Conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)
			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			client.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := client.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}
