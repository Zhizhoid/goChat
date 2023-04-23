package main

type Action struct {
	Action     string `json:"action"`
	ObjectName string `json:"object"`
}

type DefinedAction interface {
	Process(*Database) Response
}

type JsonObject interface {
	GetCreateAction() (DefinedAction, error)
	GetUpdateAction() (DefinedAction, error)
	GetReadAction() (DefinedAction, error)
	GetDeleteAction() (DefinedAction, error)
	GetLoginAction() (DefinedAction, error)
}
