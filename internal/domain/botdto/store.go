package botdto

type UserInfoDTO struct {
	Registered bool     `json:"registered,omitempty"`
	Groups     []string `json:"groups,omitempty"`
	ChatId     string   `json:"chat_id,omitempty"`
}
