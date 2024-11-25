package chat

import (
	"encoding/json"
	"fmt"
)

// Message representa un mensaje enviado entre los usuarios.
type Message struct {
	Username string `json:"username"`
	Content  string `json:"content"`
	Room     string `json:"room,omitempty"` // Opcional: útil para salas de chat
}

// ParseMessage convierte un JSON a un objeto Message.
func ParseMessage(jsonData []byte) (*Message, error) {
	var msg Message
	err := json.Unmarshal(jsonData, &msg)
	if err != nil {
		return nil, fmt.Errorf("error parsing message: %v", err)
	}
	return &msg, nil
}

// Serialize convierte un objeto Message a JSON.
func (m *Message) Serialize() []byte {
	jsonData, _ := json.Marshal(m)
	return jsonData
}
