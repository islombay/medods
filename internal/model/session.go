package model

type Session struct {
	ID     string  `json:"id"`
	Hash   *string `json:"-"`
	UserID *string `json:"user_id"`
	IP     *string `json:"ip"`

	At
}
