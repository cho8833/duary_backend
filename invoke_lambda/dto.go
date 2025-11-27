package invoke_lambda

type SendEventReq struct {
	Title          string   `json:"title"`
	Body           string   `json:"body"`
	TargetMemberId []string `json:"target_member_id"`
	Stage          string   `json:"stage"`
}

type EventAction string

const (
	EventDeleted EventAction = "EVENT_DELETED"
	EventCreated EventAction = "EVENT_CREATED"
	EventUpdated EventAction = "EVENT_UPDATED"
)
