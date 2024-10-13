package chat

import (
	"encoding/json"
	"fmt"
	"sync"

	"github.com/gofiber/contrib/websocket"
)

type Message struct {
	OriginId        string `json:"originId"`
	DestinationId   string `json:"destinationId"`
	OriginName      string `json:"originName"`
	DestinationName string `json:"destinationName"`
	Content         string `json:"content"`
	Broadcast       bool   `json:"broadcast"`
}

type Client struct {
	Id                 string
	Name               string
	WebsocketConn      *websocket.Conn
	ReceiveMessageChan chan *Message
	Manager            *ChatManager
}

func NewClient(id string, name string, ws *websocket.Conn, manager *ChatManager) *Client {
	return &Client{
		Id:                 id,
		Name:               name,
		WebsocketConn:      ws,
		ReceiveMessageChan: make(chan *Message),
		Manager:            manager,
	}
}

var Wg = sync.WaitGroup{}

func (c *Client) WriteMessages() {
	defer func() {
		Wg.Done()
		c.Manager.UnSubscribeClientChan <- c
		_ = c.WebsocketConn.Close()

		var unregisterNotification = &Message{
			OriginId:   "Manager",
			OriginName: "Manager",
			Content:    fmt.Sprintf("*** %s (%s) left the room ***", c.Name, c.Id),
			Broadcast:  true,
		}

		c.Manager.BroadcastNotificationChan <- unregisterNotification
	}()

	for {
		_, msg, err := c.WebsocketConn.ReadMessage()

		if err != nil {
			fmt.Println(err)
			break
		}

		chatMessage := Message{}

		json.Unmarshal(msg, &chatMessage)
		chatMessage.OriginId = c.Id
		chatMessage.OriginName = c.Name
		fmt.Println("Message to be sent: ", chatMessage)
		c.Manager.SendMessageChan <- &chatMessage
	}
}

func (c *Client) ReadMessage() {
	defer func() {
		Wg.Done()
		_ = c.WebsocketConn.Close()
	}()

	for {
		select {
		case messageReceived := <-c.ReceiveMessageChan:
			data, _ := json.Marshal(messageReceived)
			fmt.Println("Message received: ", messageReceived)
			c.WebsocketConn.WriteMessage(websocket.TextMessage, data)
		}
	}
}
