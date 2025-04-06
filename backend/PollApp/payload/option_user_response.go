package payload

type OptionUserResponse struct {
	OptionId  string  `json:"option_id"`
	VoteCount int     `json:"vote_count"`
	Users     []*User `json:"users"`
}

type User struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}
