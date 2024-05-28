package chats

import (
	"api/internal/core"
	"database/sql"

	"github.com/jmoiron/sqlx"
)

/*

SELECT chats.id as chat_id,
	users.id as friend_id,
	users.photo as friend_photo,
	users.firstname as friend_firstname,
	users.lastname as friend_lastname,
	(SELECT sender_id FROM chats.messages WHERE chat_id = chats.id ORDER BY id DESC LIMIT 1) as last_message_sender,
	(SELECT type_id FROM chats.messages WHERE chat_id = chats.id ORDER BY id DESC LIMIT 1) as last_message_type,
	(SELECT content FROM chats.messages WHERE chat_id = chats.id ORDER BY id DESC LIMIT 1) as last_message_content

FROM chats.chats
LEFT JOIN users.users ON users.id = any(array_remove(chats.members, 29))
WHERE 29 = any(chats.members) ORDER BY (SELECT id FROM chats.messages WHERE chat_id = chats.id ORDER BY id DESC LIMIT 1) DESC;

*/

type chats struct {
	db *sqlx.DB
}

func New(db *sqlx.DB) *chats {
	return &chats{db: db}
}

func (r *chats) AddMessage(message *core.Message) error {
	query := "INSERT INTO chats.messages (chat_id, sender_id, type_id, content) VALUES($1, $2, $3, $4)"
	err := r.db.QueryRow(query, message.ChatId, message.SenderId, message.TypeId, message.Content).Scan()
	if err != nil && err != sql.ErrNoRows {
		return err
	}
	return nil
}

func (r *chats) GetAll(userId int) (*[]core.Chat, error) {
	var output []core.Chat
	query := `SELECT chats.id as chat_id,
			users.id as friend_id,
			users.photo as friend_photo,
			users.firstname as friend_firstname,
			users.lastname as friend_lastname,
			(SELECT sender_id FROM chats.messages WHERE chat_id = chats.id ORDER BY id DESC LIMIT 1) as last_message_sender,
			(SELECT type_id FROM chats.messages WHERE chat_id = chats.id ORDER BY id DESC LIMIT 1) as last_message_type,
			(SELECT content FROM chats.messages WHERE chat_id = chats.id ORDER BY id DESC LIMIT 1) as last_message_content
		FROM chats.chats
		LEFT JOIN users.users ON users.id = any(array_remove(chats.members, $1))
		WHERE $1 = any(chats.members) ORDER BY (SELECT id FROM chats.messages WHERE chat_id = chats.id ORDER BY id DESC LIMIT 1) DESC;`

	rows, err := r.db.Queryx(query, userId)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var chat core.Chat
		if err := rows.StructScan(&chat); err != nil {
			return nil, err
		}
		output = append(output, chat)
	}

	return &output, nil
}

func (r *chats) GetMessages(chatId int) (*[]core.Message, error) {
	var output []core.Message
	query := "SELECT * FROM chats.messages WHERE chat_id = $1"
	rows, err := r.db.Queryx(query, chatId)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var msg core.Message
		if err := rows.StructScan(&msg); err != nil {
			return nil, err
		}
		output = append(output, msg)
	}

	return &output, nil
}

func (r *chats) IsChatMember(userId, chatId int) error {
	var output int
	query := "SELECT id FROM chats.chats WHERE id = $1 AND $2 = ANY(members)"
	if err := r.db.QueryRowx(query, chatId, userId).Scan(&output); err != nil {
		return err
	}
	return nil
}
