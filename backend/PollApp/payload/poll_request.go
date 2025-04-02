package payload

type PollRequest struct {
	Title       string       `json:"title"`
	Description string       `json:"description"`
	CreatorID   int          `json:"creator_id"`
	Options     []PollOption `json:"options"`
}
