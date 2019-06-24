package pusher

const (
	// TODO Add 'Event' suffix
	Pusher_EventTypeDescription_NewMessage       string = "NewMessage"
	Pusher_EventTypeDescription_NewFriendRequest string = "NewFriendRequestEvent"
	Pusher_EventTypeDescription_WereDeleted      string = "WereDeletedEvent"
	Pusher_EventTypeDescription_FriendAccepted   string = "FriendAcceptedEvent"
	Pusher_EventTypeDescription_TimerUpdated     string = "TimerUpdatedEvent"
	Pusher_EventTypeDescription_NewMember        string = "NewMemberEvent"
	// TODO fix typo
	Pusher_EventTypeDescription_SomeoneLeft string = "SomeoneLeftLevent"
	Pusher_EventTypeDescription_Dissolve    string = "DissolveEvent"
)

var Pusher_EventTypeDescription_value = map[string]int32{
	"NewMessage":            Pusher_EventType_NewMessage,
	"NewFriendRequestEvent": Pusher_EventType_NewFriendRequest,
	"WereDeletedEvent":      Pusher_EventType_WereDeleted,
	"FriendAcceptedEvent":   Pusher_EventType_FriendAccepted,
	"TimerUpdatedEvent":     Pusher_EventType_TimerUpdated,
	"NewMemberEvent":        Pusher_EventType_NewMember,
	"SomeoneLeftLevent":     Pusher_EventType_SomeoneLeft,
	"DissolveEvent":         Pusher_EventType_Dissolve,
}
