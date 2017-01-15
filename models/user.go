package models

import (
	"database/sql/driver"
	"fmt"
	"time"

	"github.com/ShyBearStudio/tbot-admindashboard/dbutils"
	"github.com/satori/go.uuid"
)

type UserRoleType string

const (
	UnknownRole UserRoleType = "UnknownRole"
	AllRoles    UserRoleType = ""
	AdminRole   UserRoleType = "admin"
	UserRole    UserRoleType = "user"
)

func (role *UserRoleType) Scan(value interface{}) error {
	_ = "breakpoint"
	strValue, ok := value.(string)
	if !ok {
		return fmt.Errorf("Cannot scan user role, source is not string.")
	}
	switch strValue {
	case AdminRole.String():
		*role = AdminRole
	case UserRole.String():
		*role = UserRole
	default:
		return fmt.Errorf("Unrecognized role value: '%s'", strValue)
	}
	return nil
}

func (role UserRoleType) Value() (driver.Value, error) {
	_ = "breakpoint"
	strValue := role.String()
	if strValue == UnknownRole.String() {
		return nil, fmt.Errorf("Unrecognized role type: '%d'", role)
	}

	return strValue, nil
}

func (role UserRoleType) String() string {
	if !(role == AdminRole || role == UserRole) {
		return string(UnknownRole)
	}
	return string(role)
}

type User struct {
	Id        int
	Uuid      string
	Name      string
	Email     string
	Password  string
	Role      UserRoleType
	CreatedAt time.Time
}

type Session struct {
	Id        int
	Uuid      string
	Email     string
	UserId    int
	CreatedAt time.Time
}

func (db *Db) CreateSession(user *User) (session Session, err error) {
	statement := "insert into sessions (uuid, email, user_id, created_at) values ($1, $2, $3, $4) " +
		"returning id, uuid, email, user_id, created_at"
	stmt, err := db.Prepare(statement)
	if err != nil {
		return
	}
	defer stmt.Close()
	// use QueryRow to return a row and scan the returned id into the Session struct
	err = stmt.QueryRow(
		uuid.NewV4(), user.Email, user.Id, time.Now()).Scan(
		&session.Id, &session.Uuid, &session.Email, &session.UserId, &session.CreatedAt)
	return
}

// Check if session is valid in the database
func (db *Db) CheckSession(session *Session) (valid bool, err error) {
	valid = false
	err = db.QueryRow("SELECT id, uuid, email, user_id, created_at FROM sessions WHERE uuid = $1",
		session.Uuid).
		Scan(&session.Id, &session.Uuid, &session.Email, &session.UserId, &session.CreatedAt)
	if err == nil && session.Id != 0 {
		valid = true
	}
	return
}

func (db *Db) User(session *Session) (*User, error) {
	user := User{}
	err := db.QueryRow("SELECT id, uuid, name, email, password, role, created_at FROM users WHERE id = $1",
		session.UserId).
		Scan(&user.Id, &user.Uuid, &user.Name, &user.Email, &user.Password, &user.Role, &user.CreatedAt)
	return &user, err
}

func (db *Db) UserByEmail(email string) (user User, err error) {
	user = User{}
	err = db.QueryRow("SELECT id, uuid, name, email, password, role, created_at FROM users WHERE email = $1", email).
		Scan(&user.Id, &user.Uuid, &user.Name, &user.Email, &user.Password, &user.Role, &user.CreatedAt)
	return
}

// Adds a new user, save user info into the database
func (db *Db) AddUser(name, email, password string, role UserRoleType) (user User, err error) {
	_ = "breakpoint"
	stmt, err := db.Prepare("insert into users (uuid, name, email, password, role, created_at) " +
		"values ($1, $2, $3, $4, $5, $6) returning id, uuid, password, created_at")
	if err != nil {
		return
	}
	defer stmt.Close()
	user = User{Name: name, Email: email, Role: role}
	// use QueryRow to return a row and scan the returned id into the User struct
	err = stmt.QueryRow(uuid.NewV4(), name, email, dbutils.Encrypt(password), role, time.Now()).
		Scan(&user.Id, &user.Uuid, &user.Password, &user.CreatedAt)
	return
}

func (db *Db) Users() (users []User, err error) {
	rows, err := db.Query(
		"SELECT id, uuid, name, email, password, role, created_at FROM users")
	if err != nil {
		return
	}
	for rows.Next() {
		user := User{}
		if err = rows.Scan(
			&user.Id, &user.Uuid, &user.Name, &user.Email, &user.Password, &user.Role, &user.CreatedAt); err != nil {
			return
		}
		users = append(users, user)
	}
	rows.Close()
	return
}
