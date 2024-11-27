package main

import (
	"encoding/json"
	"letsGOTalk/internal/chat"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func serveWs(hub *chat.Hub, w http.ResponseWriter, r *http.Request) {

	room := r.URL.Query().Get("room") // Obtener el nombre de la sala desde la query
	if room == "" {
		http.Error(w, "Room name is required", http.StatusBadRequest)
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("upgrade error:", err)
		return
	}
	client := &chat.Client{
		Hub:  hub,
		Conn: conn,
		Send: make(chan []byte, 256),
		Room: room,
	}
	hub.Register <- client

	go client.WritePump()
	go client.ReadPump()
}

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
}

func handleHistory(storage *chat.Storage, w http.ResponseWriter, r *http.Request) {
	enableCors(&w)

	room := r.URL.Query().Get("room")
	if room == "" {
		http.Error(w, "Room name is required", http.StatusBadRequest)
		return
	}

	messages, err := storage.GetMessages(room, 50) // Devuelve los últimos 50 mensajes
	if err != nil {
		http.Error(w, "Failed to retrieve messages", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(messages)
}

func main() {
	stg, err := chat.NewStorage("chat.rooms")
	if err != nil {
		log.Fatalf("Error creating storage: %v", err)
	}
	hub := chat.NewHub(stg)
	go hub.Run()

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		serveWs(hub, w, r)
	})

	http.HandleFunc("/history", func(w http.ResponseWriter, r *http.Request) {
		handleHistory(stg, w, r)
	})

	log.Println("Server started on :8080")
	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe error:", err)
	}
}
