package chatroom

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"sync"
	"time"
)

type ChatRoom struct {
	ChatRoomName string

	chatMembers     map[net.Conn]chatParticipant
	chatClientsLock sync.RWMutex
}

func NewChatRoom(chatRoomName string) *ChatRoom {
	return &ChatRoom{
		ChatRoomName: chatRoomName,

		chatMembers: make(map[net.Conn]chatParticipant),
	}
}

type chatParticipant struct {
	Name string

	receivingMessages *json.Encoder
	sendingMessages   *json.Decoder
}

type chatMessage struct {
	From string
	Body string
	Time time.Time
}

func (room *ChatRoom) AddNewChatMember(conn net.Conn) {
	defer conn.Close()

	dec := json.NewDecoder(conn)

	var chatMember chatParticipant
	err := dec.Decode(&chatMember.Name)
	if err != nil {
		fmt.Println("Error decoding new chat member, they could not be added: ", err)
		return
	}

	chatMember.receivingMessages = json.NewEncoder(conn)
	chatMember.sendingMessages = dec

	room.chatClientsLock.Lock()
	room.chatMembers[conn] = chatMember
	room.chatClientsLock.Unlock()

	log.Printf("[CONNECTED] %s", chatMember.Name)

	// send initial message welcoming the member to the chatroom
	chatMember.receivingMessages.Encode(chatMessage{
		From: room.ChatRoomName,
		Body: "Welcome to the chatroom!",
		Time: time.Now(),
	})

	for {
		var msg chatMessage
		err := chatMember.sendingMessages.Decode(&msg)
		if err != nil {
			break
		}

		// in case rogue or buggy clients want to impersonate someone else
		msg.From = chatMember.Name

		log.Printf("%s: %s", msg.From, msg.Body)

		room.chatClientsLock.RLock()
		for c, otherChatMember := range room.chatMembers {
			if c == conn {
				continue
			}
			go otherChatMember.receivingMessages.Encode(msg)
		}
		room.chatClientsLock.RUnlock()
	}

	room.chatClientsLock.Lock()
	delete(room.chatMembers, conn)
	room.chatClientsLock.Unlock()

	log.Printf("[DISCONNECTED] %s", chatMember.Name)
}
