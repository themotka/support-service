package postgres

import (
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/lib/pq"
	"support-bot/internal/storage"
)

type Storage struct {
	db *sql.DB
}

func NewStorage(conn string) (*Storage, error) {
	const op = "storage.postgres.NewStorage"
	db, err := sql.Open("postgres", conn)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return &Storage{db: db}, nil
}

func (s *Storage) CreateMessage(text string, conversationID int) (messageID int, err error) {
	const op = "storage.postgres.CreateMessage"

	query := `INSERT INTO messages (content, conversation_id, is_manager) VALUES ($1, $2, $3) RETURNING message_id`
	err = s.db.QueryRow(query, text, conversationID, true).Scan(&messageID)
	if err != nil {
		return -1, fmt.Errorf("%s: %w", op, err)
	}

	return messageID, nil
}

func (s *Storage) CreateUser(userID int, platformID int) error {
	const op = "storage.postgres.CreateUser"

	query := "INSERT INTO users (user_id, platform_id) VALUES ($1, $2)"
	_, err := s.db.Exec(query, userID, platformID)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func (s *Storage) CreateConversation(userID, managerID int) (conversationID int, err error) {
	const op = "storage.postgres.CreateConversation"

	query := "INSERT INTO conversations (user_id, manager_id) VALUES ($1, $2) RETURNING conversation_id"
	err = s.db.QueryRow(query, userID, managerID).Scan(&conversationID)
	if err != nil {
		return -1, fmt.Errorf("%s: %w", op, err)
	}
	return conversationID, nil
}

func (s *Storage) UserExists(userID int) (exists bool, err error) {
	const op = "storage.postgres.UserExists"

	query := "SELECT EXISTS (SELECT 1 FROM users WHERE user_id=$1)"
	err = s.db.QueryRow(query, userID).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("%s: %w", op, err)
	} else if !exists {
		return false, fmt.Errorf("%s: %w", op, storage.ErrUserNotFound)
	}

	return exists, nil
}

func (s *Storage) PlatformID(platform string) (platformID int, err error) {
	const op = "storage.postgres.PlatformID"

	query := "SELECT platform_id FROM platforms WHERE name=$1"
	err = s.db.QueryRow(query, platform).Scan(&platformID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return -1, fmt.Errorf("%s: %w", op, storage.ErrPlatformNotFound)
		}
		return -1, fmt.Errorf("%s: %w", op, err)
	}
	return platformID, nil
}

func (s *Storage) ConversationID(userID int) (conversationID int, err error) {
	const op = "storage.postgres.ConversationID"

	query := "SELECT conversation_id FROM conversations WHERE user_id=$1"
	err = s.db.QueryRow(query, userID).Scan(&conversationID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return -1, fmt.Errorf("%s: %w", op, storage.ErrPlatformNotFound)
		}
		return -1, fmt.Errorf("%s: %w", op, err)
	}
	return conversationID, nil
}
