package chat

import (
	"database/sql"
	"fmt"

	_ "modernc.org/sqlite" // Driver para SQLite
)

type Storage struct {
	DB *sql.DB
}

func NewStorage(databasePath string) (*Storage, error) {
	db, err := sql.Open("sqlite", databasePath)
	if err != nil {
		return nil, fmt.Errorf("error opening database: %v", err)
	}

	query := `
	CREATE TABLE IF NOT EXISTS messages (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		username TEXT NOT NULL,
		content TEXT NOT NULL,
		room TEXT NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);`
	if _, err := db.Exec(query); err != nil {
		return nil, fmt.Errorf("error creating table: %v", err)
	}

	return &Storage{DB: db}, nil
}

func (s *Storage) SaveMessage(msg *Message) error {
	query := `
	INSERT INTO messages (username, content, room) 
	VALUES (?, ?, ?);`
	_, err := s.DB.Exec(query, msg.Username, msg.Content, msg.Room)
	return err
}

func (s *Storage) GetMessages(room string, limit int) ([]Message, error) {
	query := `
	SELECT username, content, room, created_at 
	FROM messages 
	WHERE room = ? 
	ORDER BY created_at DESC 
	LIMIT ?;`
	rows, err := s.DB.Query(query, room, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []Message
	for rows.Next() {
		var msg Message
		var createdAt string
		if err := rows.Scan(&msg.Username, &msg.Content, &msg.Room, &createdAt); err != nil {
			return nil, err
		}
		messages = append(messages, msg)
	}
	return messages, nil
}
