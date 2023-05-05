package database

import (
	"fmt"
)

type Room struct{}

// Create action TOCHECK
type RoomCreate struct {
	Data struct {
		Token string `json:"token"`
		Name  string `json:"name"`
	} `json:"data"`
}

func (room *Room) GetCreateAction() (DefinedAction, error) {
	return &RoomCreate{}, nil
}

func (action *RoomCreate) Process(db *Database) Response {
	if action.Data.Name == "" {
		return roomResponse("create", false, "Creating room failed, room name cannon be empty")
	}

	err := db.AddRoom(action.Data.Name, action.Data.Token)
	if err != nil {
		return roomResponse("create", false, err.Error())
	}

	return roomResponse("create", true, fmt.Sprintf("Room %v successfully created", action.Data.Name))
}

// Update action TODO
type RoomUpdate struct {
	Data struct {
		ID   uint64 `json:"id"`
		Name string `json:"name"`
	} `json:"data"`
}

func (room *Room) GetUpdateAction() (DefinedAction, error) {
	return &RoomUpdate{}, nil
}

func (action *RoomUpdate) Process(db *Database) Response {
	fmt.Printf("Updating room: %v %v\n", action.Data.ID, action.Data.Name)

	if action.Data.Name == "" {
		return roomResponse("update", false, "Updating room failed, room name cannon be empty")
	}

	err := db.UpdateRoom(action.Data.ID, action.Data.Name)
	if err != nil {
		return roomResponse("update", false, err.Error())
	}

	return userResponse("update", true, fmt.Sprintf("Room %v successfully created", action.Data.Name))
}

// Delete action TODO
type RoomDelete struct {
	Data struct {
		ID uint64 `json:"id"`
	} `json:"data"`
}

func (room *Room) GetDeleteAction() (DefinedAction, error) {
	return &RoomDelete{}, nil
}

func (action *RoomDelete) Process(db *Database) Response {
	fmt.Printf("Deleting room %v\n", action.Data.ID)

	err := db.DeleteRoom(action.Data.ID)
	if err != nil {
		return roomResponse("delete", false, err.Error())
	}
	return roomResponse("delete", true, fmt.Sprintf("Room successfully deleted"))
}

// Login action
type RoomLogin struct {
	Data struct {
		Token    string `json:"token"`
		RoomName string `json:"roomName"`
	} `json:"data"`
}

func (m *Room) GetLoginAction() (DefinedAction, error) {
	return &RoomLogin{}, nil
}

func (action *RoomLogin) Process(db *Database) Response {
	err := db.LoginRoom(action.Data.RoomName, action.Data.Token)
	if err != nil {
		return roomResponse("login", false, err.Error())
	}

	return roomResponse("login", true, "Successfully joined the room")
}

// Read action
type RoomRead struct {
	Data struct {
		ID uint64 `json:"id"`
	} `json:"data"`
}

func (action *Room) GetReadAction() (DefinedAction, error) {
	return &RoomRead{}, nil
}

func (action *RoomRead) Process(db *Database) Response {
	return roomResponse("read", false, "Unimplemented")
}

// OTHER
func roomResponse(action string, success bool, status string) Response {
	return Response{
		Action:     action,
		ObjectName: "room",
		Success:    success,
		Status:     status,
	}
}
