package payload

type PollResponse struct {
	Id          int           `json:"id"`
	Description string        `json:"description"`
	PollOptions []*OptionData `json:"options"`
}

type OptionData struct {
	OptionId   int    `json:"option_id"`
	OptionText string `json:"option_text"`
	VoteCount  int    `json:"vote_count"`
}
