package database

import (
	"fmt"
)

type User struct{}

// Create action
type UserCreate struct {
	Data struct {
		Username string `json:"username"`
		Password string `json:"password"`
		Name     string `json:"name"`
	} `json:"data"`
}

func (user *User) GetCreateAction() (DefinedAction, error) {
	return &UserCreate{}, nil
}

func (action *UserCreate) Process(db *Database) Response {
	if action.Data.Username == "" {
		return userResponse("create", false, "Username cannon be empty")
	}

	if action.Data.Password == "" {
		return userResponse("create", false, "Password cannon be empty")
	}

	if action.Data.Name == "" {
		return userResponse("create", false, "Name cannon be empty")
	}

	err := db.AddUser(action.Data.Username, action.Data.Password, action.Data.Name)
	if err != nil {
		return userResponse("create", false, err.Error())
	}
	return userResponse("create", true, fmt.Sprintf("User %v successfully created", action.Data.Username))
}

// Update action
type UserUpdate struct {
	Data struct {
		Token    string `json:"token"`
		Username string `json:"username"`
		Password string `json:"password"`
		Name     string `json:"name"`
	} `json:"data"`
}

func (user *User) GetUpdateAction() (DefinedAction, error) {
	return &UserUpdate{}, nil
}

func (action *UserUpdate) Process(db *Database) Response {
	err := db.UpdateUser(action.Data.Token, action.Data.Username, action.Data.Password, action.Data.Name)
	if err != nil {
		return userResponse("update", false, err.Error())
	}
	return userResponse("update", true, fmt.Sprintf("User %v successfully updated", action.Data.Username))
}

// Delete action
type UserDelete struct {
	Data struct {
		Token string `json:"token"`
	} `json:"data"`
}

func (user *User) GetDeleteAction() (DefinedAction, error) {
	return &UserDelete{}, nil
}

func (action *UserDelete) Process(db *Database) Response {
	err := db.DeleteUser(action.Data.Token)
	if err != nil {
		return userResponse("delete", false, err.Error())
	}
	return userResponse("delete", true, fmt.Sprintf("User successfully deleted"))
}

// Login action
type UserLogin struct {
	Data struct {
		Username string `json:"username"`
		Password string `json:"password"`
	} `json:"data"`
}

func (m *User) GetLoginAction() (DefinedAction, error) {
	return &UserLogin{}, nil
}

func (action *UserLogin) Process(db *Database) Response {
	token, err := db.LoginUser(action.Data.Username, action.Data.Password)
	if err != nil {
		return userResponse("login", false, err.Error())
	}

	return Response{
		Action:     "login",
		ObjectName: "user",
		Success:    true,
		Status:     "Successfully logged in",
		Token:      token,
	}
}

// Read action
type UserRead struct {
	Data struct {
		Token       string `json:"token"`
		GetRoomList bool   `json:"getRoomList"`
	} `json:"data"`
}

func (action *User) GetReadAction() (DefinedAction, error) {
	return &UserRead{}, nil
}

func (action *UserRead) Process(db *Database) Response {
	response := Response{
		Action:     "read",
		ObjectName: "user",
		Success:    true,
		Status:     "Successfully read user",
		ReadResponse: ReadResponse{
			UserReadResponse: UserReadResponse{},
		},
	}

	username, name, err := db.ReadUser(action.Data.Token)
	if err != nil {
		return userResponse("read", false, err.Error())
	}

	response.ReadResponse.UserReadResponse.Name = name

	if action.Data.GetRoomList {
		rooms, err := db.ReadUserRooms(username)
		if err != nil {
			return userResponse("read", false, err.Error())
		}

		response.ReadResponse.UserReadResponse.Rooms = rooms
	}

	return response
}

// OTHER
func userResponse(action string, success bool, status string) Response {
	return Response{
		Action:     action,
		ObjectName: "user",
		Success:    success,
		Status:     status,
	}
}
