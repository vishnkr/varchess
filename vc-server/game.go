package main

import (
	//"github.com/gorilla/websocket"
)


type Game struct{
	p1 *Client
	p2 *Client
	fen string
	moveList []string
	result string
}