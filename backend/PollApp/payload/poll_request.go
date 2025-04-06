package payload

type PollRequest struct {
	Id          string        `json:"id"`
	Description string        `json:"description"`
	CreatorId   int           `json:"creator_id"`
	Options     []*PollOption `json:"options"`
}
