package utils

import (
	"encoding/json"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo"
)

type Message struct {
	Username string `json:"user"`
	Body     string `json:"body"`
}

type Chat struct {
	Upgrader    *websocket.Upgrader
	Broadcaster Broadcaster
}

func NewChat() *Chat {
	return &Chat{
		Upgrader:    &websocket.Upgrader{},
		Broadcaster: newBroadcaster(4),
	}
}

func (ch *Chat) Chatify(c echo.Context) error {
	ws, err := ch.Upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		return err
	}

	messages := make(chan Message)
	ch.Broadcaster.Register(messages)
	defer ch.Broadcaster.Unregister(messages)
	defer ws.Close()

	go func() {
		for msg := range messages {
			text, err := json.Marshal(msg)
			if err != nil {
				c.Logger().Error(err)
				continue
			}

			ws.WriteMessage(websocket.TextMessage, text)
		}
	}()

	for {
		var msg Message
		_, text, err := ws.ReadMessage()
		if err != nil {
			c.Logger().Error(err)
			break
		}

		err = json.Unmarshal(text, &msg)
		if err != nil {
			c.Logger().Error(err)
		}

		ch.Broadcaster.SendExcept(msg, messages)
	}
	return nil
}
