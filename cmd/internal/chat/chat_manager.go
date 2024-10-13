package chat

type ChatManager struct {
	Clients                   []*Client
	SubscribeClientChan       chan *Client
	UnSubscribeClientChan     chan *Client
	BroadcastNotificationChan chan *Message
	SendMessageChan           chan *Message
}

func (manager *ChatManager) Start() {
	for {
		select {
		case channel := <-manager.SubscribeClientChan:
			manager.Clients = append(manager.Clients, channel)

		case channel := <-manager.UnSubscribeClientChan:
			for i, client := range manager.Clients {
				if client.Id == channel.Id {
					manager.Clients = append(manager.Clients[:i], manager.Clients[:i+1]...)
				}
			}

		case channel := <-manager.SendMessageChan:
			for _, client := range manager.Clients {
				if client.Id == channel.DestinationId {
					client.ReceiveMessageChan <- channel
				}
			}

		case channel := <-manager.BroadcastNotificationChan:
			for _, client := range manager.Clients {
				if client.Id != channel.OriginId {
					client.ReceiveMessageChan <- channel
				}
			}
		}
	}
}

var Manager = &ChatManager{
	Clients:                   make([]*Client, 0),
	SubscribeClientChan:       make(chan *Client),
	UnSubscribeClientChan:     make(chan *Client),
	BroadcastNotificationChan: make(chan *Message),
	SendMessageChan:           make(chan *Message),
}
