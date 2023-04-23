package database

import (
	"errors"
	"fmt"
	"time"
)

type Message struct {
	Content
	SentAt    time.Time
	AuthorID  uint64
	ReplyToID uint64
	Edited    bool
}

// Create action
type MessageCreate struct {
	Data struct {
		Content   `json:"content"`
		SenderID  uint64 `json:"senderId"`
		RoomID    uint64 `json:"roomId"`
		ReplyToID uint64 `json:"replyToId"`
	} `json:"data"`
}

func (m *Message) GetCreateAction() (DefinedAction, error) {
	return &MessageCreate{}, nil
}

func (action *MessageCreate) Process(db *Database) Response {
	fmt.Printf("Sending message %v\n", action.Data)

	if action.Data.Content.IsEmpty() {
		return messageResponse("create", false, "Sending message failed, message cannon be empty")
	}

	err := db.SendMessage(action.Data.Content, action.Data.SenderID, action.Data.RoomID, action.Data.ReplyToID)
	if err != nil {
		return messageResponse("create", false, err.Error())
	}

	return messageResponse("create", true, "Message successfully sent")
}

// Update action
type MessageUpdate struct {
	Data struct {
		RoomID    uint64 `json:"roomId"`
		MessageID uint64 `json:"messageId"`
		Content   `json:"content"`
		ReplyToID uint64 `json:"replyToId"`
	} `json:"data"`
}

func (m *Message) GetUpdateAction() (DefinedAction, error) {
	return &MessageUpdate{}, nil
}

func (action *MessageUpdate) Process(db *Database) Response {
	fmt.Printf("Updating message %v\n", action.Data)

	if action.Data.Content.IsEmpty() {
		return messageResponse("update", false, "Updating message failed, message cannon be empty")
	}

	err := db.UpdateMessage(action.Data.RoomID, action.Data.MessageID, action.Data.Content, action.Data.ReplyToID)
	if err != nil {
		return messageResponse("update", false, err.Error())
	}

	return messageResponse("update", true, "Message successfully updated")
}

// Delete action
type MessageDelete struct {
	Data struct {
		RoomID    uint64 `json:"roomId"`
		MessageID uint64 `json:"messageId"`
	} `json:"data"`
}

func (m *Message) GetDeleteAction() (DefinedAction, error) {
	return &MessageDelete{}, nil
}

func (action *MessageDelete) Process(db *Database) Response {
	fmt.Printf("Deleting message %v\n", action.Data)

	err := db.DeleteMessage(action.Data.RoomID, action.Data.MessageID)
	if err != nil {
		return messageResponse("delete", false, err.Error())
	}

	return messageResponse("delete", true, "Message successfully deleted")
}

// Login action
func (m *Message) GetLoginAction() (DefinedAction, error) {
	return nil, errors.New("No login action")
}

// Read action
type MessageRead struct {
	Data struct {
		RoomID    uint64 `json:"roomId"`
		MessageID uint64 `json:"messageId"`
	} `json:"data"`
}

func (action *Message) GetReadAction() (DefinedAction, error) {
	return &MessageRead{}, nil
}

func (action *MessageRead) Process(db *Database) Response {
	return Response{}
}

// OTHER
func messageResponse(action string, success bool, status string) Response {
	return Response{
		Action:     action,
		ObjectName: "message",
		Success:    success,
		Status:     status,
	}
}
