package main

type Response struct {
	Action     string  `json:"action"`
	ObjectName string  `json:"object"`
	Success    bool    `json:"success"`
	Status     string  `json:"status"`
	Message    Message `json:"message"`
	ID         uint64  `json:"id"`
}
