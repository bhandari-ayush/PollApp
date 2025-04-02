package payload

type VoteRequest struct {
	PollID   int `json:"poll_id"`
	OptionID int `json:"option_id"`
	UserID   int `json:"user_id"`
}
