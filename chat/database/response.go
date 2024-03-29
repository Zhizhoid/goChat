package database

type Response struct {
	Action       string  `json:"action"`
	ObjectName   string  `json:"object"`
	Success      bool    `json:"success"`
	Status       string  `json:"status"`
	Message      Message `json:"message"`
	Token        string  `json:"token"`
	ReadResponse `json:"read"`
}

type UserReadResponse struct {
	Name     string   `json:"name"`
	Username string   `json:"username"`
	Rooms    []string `json:"rooms"`
}

type ReadResponse struct {
	UserReadResponse `json:"user"`
}
