package pusher

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/gorilla/websocket"
)

const (
	WebSocketStateNew = iota
	WebSocketStateConnected
	WebSocketStateDisconnected
	WebSocketStateClosed
)

type EventHandler func(interface{})

type Pusher struct {
	conn        *websocket.Conn
	Url         string
	HandleEvent EventHandler
}

func NewPusher(url string, handleEvent EventHandler) *Pusher {
	return &Pusher{
		Url:         url,
		HandleEvent: handleEvent,
	}
}

// https://github.com/gorilla/websocket/blob/master/examples/echo/client.go
// wss://stargate.aiursoft.com/Listen/Channel?Id=&Key=
// This is a synchronize call, it returns when connection closed.
// No auto reconnect provided.
func (p *Pusher) Connect(interrupt <-chan struct{}) error {
	var err error
	p.conn, _, err = websocket.DefaultDialer.Dial(p.Url, nil)
	if err != nil {
		return err
	}
	// CLose connection when return
	defer p.conn.Close()

	// Close done to exit main loop
	done := make(chan struct{})
	errChan := make(chan error)
	// CLose errChan when return
	defer close(errChan)
	// Message loop
	go func() {
		// Close done when message loop exit, notify main loop exit.
		defer close(done)
		for {
			_, message, err := p.conn.ReadMessage()
			if err != nil {
				// Send error and exit
				errChan <- err
				return
			}
			event, err := DecodePusherEvent(message)
			if err != nil {
				// Send error and exit
				errChan <- err
				return
			}
			p.HandleEvent(event)
		}
	}()

	ticker := time.NewTicker(45 * time.Second)
	defer ticker.Stop()

	// wait connection close or interrupt
	for {
		select {
		case <-done:
			// Message loop exit
			select {
			case err := <-errChan:
				// Message loop exit with error
				return err
			default:
				panic("cannot reach")
				return nil
			}
		case <-ticker.C:
			// Heartbeat
			err := p.conn.WriteMessage(websocket.PingMessage, []byte{})
			if err != nil {
				return err
			}
		case <-interrupt:
			// Cleanly close the connection by sending a close message and then
			// waiting (with timeout) for the server to close the connection.
			err := p.conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				return err
			}
			select {
			case <-done:
			case <-time.After(time.Second):
			}
			return nil
		}
	}
}

func DecodePusherEvent(message []byte) (interface{}, error) {
	var err error
	event1 := &Pusher_Event{}
	err = json.Unmarshal(message, event1)
	if err != nil {
		return event1, err
	}
	var event interface{}
	switch int32(event1.Type) {
	case Pusher_EventType_NewMessage:
		event = &Pusher_NewMessageEvent{}
	case Pusher_EventType_NewFriendRequest:
		event = &Pusher_NewFriendRequestEvent{}
	case Pusher_EventType_WereDeleted:
		event = &Pusher_WereDeletedEvent{}
	case Pusher_EventType_FriendAccepted:
		event = &Pusher_FriendAcceptedEvent{}
	case Pusher_EventType_TimerUpdated:
		event = &Pusher_TimerUpdatedEvent{}
	case Pusher_EventType_NewMember:
		event = &Pusher_NewMemberEvent{}
	case Pusher_EventType_SomeoneLeft:
		event = &Pusher_SomeoneLeftEvent{}
	case Pusher_EventType_Dissolve:
		event = &Pusher_DissolveEvent{}
	default:
		return event1, &InvalidEventTypeError{event1.Type}
	}
	err = json.Unmarshal(message, event)
	if err != nil {
		return event, err
	}
	return event, nil
}

type InvalidEventTypeError struct {
	EventType uint32
}

func (i *InvalidEventTypeError) Error() string {
	return fmt.Sprintf("invalid event type: %d", i.EventType)
}
