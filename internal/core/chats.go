package core

type Chat struct {
	ID                 int     `json:"id" db:"chat_id"`
	FriendId           *int    `json:"friend_id" db:"friend_id"`
	FriendPhoto        *string `json:"friend_photo" db:"friend_photo"`
	FriendFullname     *string `json:"friend_fullname" db:"friend_fullname"`
	LastMessageSender  *int    `json:"last_message_sender" db:"last_message_sender"`
	LastMessageType    *int    `json:"last_message_type" db:"last_message_type"`
	LastMessageContent *string `json:"last_message_content" db:"last_message_content"`
}

type Message struct {
	ID        int    `json:"id" db:"id"`
	ChatId    int    `json:"chat_id" db:"chat_id"`
	SenderId  int    `json:"sender_id" db:"sender_id"`
	TypeId    int    `json:"type_id" db:"type_id"`
	Content   string `json:"content" db:"content"`
	CreatedAt string `json:"created_at" db:"created_at"`
}
