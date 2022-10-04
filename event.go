package main

type Event struct {
	Event string  `json:"event"`
	Data  *Answer `json:"data"`
}
