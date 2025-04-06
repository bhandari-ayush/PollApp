package payload

type VoteRequest struct {
	Id       int `json:"id"`
	PollId   int `json:"poll_id"`
	OptionId int `json:"option_id"`
	UserId   int `json:"user_id"`
}
