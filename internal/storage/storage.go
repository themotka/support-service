package storage

import "errors"

var (
	ErrUserNotFound         = errors.New("user not found")
	ErrPlatformNotFound     = errors.New("platform not found")
	ErrConversationNotFound = errors.New("conversation not found")
)

type Storage interface {
	CreateMessage(text string, conversationID int) (int, error)
	CreateUser(userID, platformID int) error
	CreateConversation(userID, managerID int) (int, error)
	UserExists(userID int) (bool, error)
	PlatformID(platform string) (int, error)
	ConversationID(userID int) (int, error)
}

type NewClaim struct {
	ExternalID string
	Platform   string
	Text       string
}
