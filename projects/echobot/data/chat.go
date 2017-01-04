package data

import (
	"database/sql"
	"time"

	"github.com/mrd0ll4r/tbotapi"
)

// Chat contains information about the chat a message originated from.
type Chat struct {
	Id        int       // Unique identifier for this chat.
	Type      string    // Type of chat, can be either "private", "group" or "channel". Check Is(PrivateChat|GroupChat|Channel)() methods.
	Title     *string   // Title for channels and group chats.
	UserName  *string   // Username for private chats and channels if available.
	FirstName *string   // First name of the other party in a private chat.
	LastName  *string   // Last name of the other party in a private chat.
	CreatedAt time.Time // Time chat was registered
	Active    bool      // True if chat is active
}

// chat registered?
func IsRegisteredChat(id int) (is bool, err error) {
	_ = "breakpoint"
	var resultId int
	err = Db.QueryRow("SELECT id FROM echobot_chats WHERE id = $1", id).
		Scan(&resultId)
	_ = "breakpoint"
	if err == sql.ErrNoRows {
		return false, nil
	} else if err != nil {
		return
	}

	return true, err
}

// add new chat
func CreateChat(chat tbotapi.Chat) error {
	_ = "breakpoint"
	statement := "insert into echobot_chats (id, type, title, username, firstname, lastname, created_at, active) values ($1, $2, $3, $4, $5, $6, $7, $8)"
	stmt, err := Db.Prepare(statement)
	if err != nil {
		return err
	}
	defer stmt.Close()
	// use QueryRow to return a row and scan the returned id into the User struct
	_, err = stmt.Query(chat.ID, chat.Type, chat.Title, chat.Username, chat.FirstName, chat.LastName, time.Now(), true)
	return err
}

// update chat status (ie from active=true to active=false)
// get chats
func Chats() (chats []Chat, err error) {
	rows, err := Db.Query("SELECT id, type, title, username, firstname, lastname, created_at, active FROM echobot_chats")
	if err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		chat := Chat{}
		if err = rows.Scan(&chat.Id, &chat.Type, &chat.Title, &chat.UserName, &chat.FirstName, &chat.LastName, &chat.CreatedAt, &chat.Active); err != nil {
			return
		}
		chats = append(chats, chat)
	}
	return
}
