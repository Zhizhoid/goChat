package main

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
	fmt.Printf("Creating user: %s\n", action.Data.Name)

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
		SessionID uint64 `json:"sessionId"`
		Username  string `json:"username"`
		Password  string `json:"password"`
		Name      string `json:"name"`
	} `json:"data"`
}

func (user *User) GetUpdateAction() (DefinedAction, error) {
	return &UserUpdate{}, nil
}

func (action *UserUpdate) Process(db *Database) Response {
	id, err := db.sm.GetUserIDFromSession(action.Data.SessionID)
	if err != nil {
		return userResponse("update", false, err.Error())
	}

	fmt.Printf("Updating user: %v %v %v %v\n", id, action.Data.Username, action.Data.Password, action.Data.Name)

	err = db.UpdateUser(id, action.Data.Username, action.Data.Password, action.Data.Name)
	if err != nil {
		return userResponse("update", false, err.Error())
	}
	return userResponse("update", true, fmt.Sprintf("User %v successfully updated", action.Data.Username))
}

// Delete action
type UserDelete struct {
	Data struct {
		SessionID uint64 `json:"sessionId"`
	} `json:"data"`
}

func (user *User) GetDeleteAction() (DefinedAction, error) {
	return &UserDelete{}, nil
}

func (action *UserDelete) Process(db *Database) Response {
	id, err := db.sm.GetUserIDFromSession(action.Data.SessionID)
	if err != nil {
		return userResponse("delete", false, err.Error())
	}

	fmt.Printf("Deleting user %v\n", id)

	err = db.DeleteUser(id)
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
	sessionId, err := db.LoginUser(action.Data.Username, action.Data.Password)
	if err != nil {
		return userResponse("login", false, err.Error())
	}

	return Response{
		Action:     "login",
		ObjectName: "user",
		Success:    true,
		SessionID:  sessionId,
	}
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
