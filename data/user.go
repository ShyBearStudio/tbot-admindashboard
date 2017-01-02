package data

import (
	"database/sql/driver"
	"errors"
	"fmt"
	"time"

	"github.com/satori/go.uuid"
)

type UserRoleType string

const (
	UnknownRole UserRoleType = "UnknownRole"
	AllRoles    UserRoleType = ""
	AdminRole   UserRoleType = "admin"
	UserRole    UserRoleType = "user"
	VendorRole  UserRoleType = "vendor"
)

func (role *UserRoleType) Scan(value interface{}) error {
	_ = "breakpoint"
	strValue, ok := value.(string)
	if !ok {
		return errors.New("Cannot scan user role, source is not string.")
	}
	switch strValue {
	case AdminRole.String():
		*role = AdminRole
	case UserRole.String():
		*role = UserRole
	case VendorRole.String():
		*role = VendorRole
	default:
		return errors.New(fmt.Sprintf("Unrecognized role value: '%s'", strValue))
	}
	return nil
}

func (role UserRoleType) Value() (driver.Value, error) {
	_ = "breakpoint"
	strValue := role.String()
	if strValue == UnknownRole.String() {
		return nil, errors.New(fmt.Sprintf("Unrecognized role type: '%d'", role))
	}

	return strValue, nil
}

func (role UserRoleType) String() string {
	if !(role == AdminRole || role == UserRole || role == VendorRole) {
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

// Create a new session for an existing user
func (user *User) CreateSession() (session Session, err error) {
	statement := "insert into sessions (uuid, email, user_id, created_at) values ($1, $2, $3, $4) returning id, uuid, email, user_id, created_at"
	stmt, err := Db.Prepare(statement)
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
func (session *Session) Check() (valid bool, err error) {
	valid = false
	err = Db.QueryRow("SELECT id, uuid, email, user_id, created_at FROM sessions WHERE uuid = $1", session.Uuid).
		Scan(&session.Id, &session.Uuid, &session.Email, &session.UserId, &session.CreatedAt)
	if err == nil && session.Id != 0 {
		valid = true
	}
	return
}

// Gets the user from the session
func (session *Session) User() (user User, err error) {
	user = User{}
	err = Db.QueryRow("SELECT id, uuid, name, email, role, created_at FROM users WHERE id = $1", session.UserId).
		Scan(&user.Id, &user.Uuid, &user.Name, &user.Email, &user.Role, &user.CreatedAt)
	return
}

// Gets a single user given the email
func UserByEmail(email string) (user User, err error) {
	user = User{}
	err = Db.QueryRow("SELECT id, uuid, name, email, password, created_at FROM users WHERE email = $1", email).
		Scan(&user.Id, &user.Uuid, &user.Name, &user.Email, &user.Password, &user.CreatedAt)
	return
}

// Creates a new user, save user info into the database
func CreateUser(name, email, password string, role UserRoleType) (userId int, err error) {
	_ = "breakpoint"
	// Postgres does not automatically return the last insert id, because it would be wrong to assume
	// you're always using a sequence.You need to use the RETURNING keyword in your insert to get this
	// information from postgres.
	statement := "insert into users (uuid, name, email, password, role, created_at) values ($1, $2, $3, $4, $5, $6) returning id"
	stmt, err := Db.Prepare(statement)
	if err != nil {
		return
	}
	defer stmt.Close()
	// use QueryRow to return a row and scan the returned id into the User struct
	err = stmt.QueryRow(uuid.NewV4(), name, email, Encrypt(password), role, time.Now()).
		Scan(&userId)
	return
}

// Get all users in the database and returns it
func Users() (users []User, err error) {
	rows, err := Db.Query("SELECT id, uuid, name, email, password, role, created_at FROM users")
	if err != nil {
		return
	}
	for rows.Next() {
		user := User{}
		if err = rows.Scan(&user.Id, &user.Uuid, &user.Name, &user.Email, &user.Password, &user.Role, &user.CreatedAt); err != nil {
			return
		}
		users = append(users, user)
	}
	rows.Close()
	return
}
