package structs

type Messages []MessageElement

type MessageElement struct {
	Ok     bool     `json:"ok"`
	Result []Result `json:"result"`
}

type Result struct {
	UpdateID int64   `json:"update_id"`
	Message  Message `json:"message"`
}

type Message struct {
	MessageID int64  `json:"message_id"`
	From      From   `json:"from"`
	Chat      Chat   `json:"chat"`
	Date      int64  `json:"date"`
	Text      string `json:"text"`
}

type Chat struct {
	ID    int64  `json:"id"`
	Title string `json:"title"`
	Type  string `json:"type"`
}

type From struct {
	ID           int64   `json:"id"`
	IsBot        bool    `json:"is_bot"`
	FirstName    string  `json:"first_name"`
	LastName     *string `json:"last_name,omitempty"`
	Username     *string `json:"username,omitempty"`
	LanguageCode *string `json:"language_code,omitempty"`
}

type Stat struct {
	Count             int64  `json:"count" db:"count"`
	FromID            int64  `json:"from_id" db:"from_id"`
	Username          string `json:"username" db:"-"`
	FirstName         string `json:"first_name" db:"-"`
	LastName          string `json:"last_name" db:"-"`
	LengthOfWriteText int64  `json:"length_of_write_text" db:"length_of_write_text"`
}

type User struct {
	Username  string `json:"username" db:"from_username"`
	FirstName string `json:"first_name" db:"from_first_name"`
	LastName  string `json:"last_name" db:"from_last_name"`
}
