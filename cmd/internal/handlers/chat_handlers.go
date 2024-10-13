package handlers

import (
	"fmt"

	"github.com/gofiber/contrib/websocket"
	"github.com/google/uuid"
	"github.com/raghav1030/go-ws-chat/cmd/internal/chat"
)

func RegisterHandler(c *websocket.Conn) {
	chat.Wg.Add(2)

	client := chat.NewClient(uuid.New().String(), c.Params("nick"), c, chat.Manager)
	chat.Manager.SubscribeClientChan <- client

	var registerNotification = &chat.Message{
		OriginId:   "Manager",
		OriginName: "Manager",
		Content:    fmt.Sprintf("***  %s (%s) joined to this room ***", client.Name, client.Id),
		Broadcast:  true,
	}

	chat.Manager.BroadcastNotificationChan <- registerNotification

	go client.ReadMessage()
	go client.WriteMessages()

	chat.Wg.Wait()
}
