package main

import (
	"github.com/troylelandshields/easychatserver/chatroom"
	"github.com/troylelandshields/easychatserver/listener"
)

func main() {
	// create a new chat room
	chatRoom := chatroom.NewChatRoom("Women Who Go!")

	// listen for connections to the chat room
	listener.ListenForChatConnections(chatRoom)
}
