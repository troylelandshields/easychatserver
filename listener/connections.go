package listener

import (
	"fmt"
	"log"
	"net"

	"github.com/troylelandshields/easychatserver/chatroom"
)

func ListenForChatConnections(chatRoom *chatroom.ChatRoom) {
	ln, err := net.Listen("tcp", ":1234")
	if err != nil {
		log.Fatal(err)
	}
	defer ln.Close()

	fmt.Println("Listening for connections to chatroom:", chatRoom.ChatRoomName)

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Fatal(err)
		}

		go chatRoom.AddNewChatMember(conn)
	}
}
