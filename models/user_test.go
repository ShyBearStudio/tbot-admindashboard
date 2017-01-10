package models

import (
	"flag"
	"os"
	"testing"

	"github.com/ShyBearStudio/tbot-admindashboard/dbutils"
)

const (
	createDbScript = "create_tables.sql"
)

var tdb *Db

func TestMain(m *testing.M) {
	flag.Parse()
	if testing.Short() != true {
		setupDatabase()
	}
	os.Exit(m.Run())
	if !testing.Short() {
		teardownDatabase()
	}
}

func setupDatabase() {
	var err error
	tdb, err = NewDb(testDbDriver, testDbDataSourceName)
	if err != nil {
		panic(err)
	}
	err = dbutils.ExecuteSqlScript(tdb.DB, createDbScript)
	if err != nil {
		panic(err)
	}
}

func teardownDatabase() {
}

func TestManipulate(t *testing.T) {
	SkipTestIfShort(t)
	user1, user2 := AddTwoUsers(t)
	AssertUsersInDb(t, user1, user2)
	session := CreateSession(t, user1)
	AssertValidSession(t, &session)
	AssertUserFromSession(t, &session, &user1)
	AssertUserByEmailExtraction(t, &user1)
}

func AssertUserByEmailExtraction(t *testing.T, userToValidate *User) {
	user, err := tdb.UserByEmail(userToValidate.Email)
	if err != nil {
		t.Errorf("Cannot extract user by email: '%v'", user.Email)
	}
	if user != *userToValidate {
		t.Errorf("Expected user '%v' but was found '%v'", userToValidate, user)
	}
}

func AssertUserFromSession(t *testing.T, session *Session, userToValidate *User) {
	user, err := tdb.User(session)
	if err != nil {
		t.Errorf("Cannot get user for session '%v'", session)
	}
	if user != *userToValidate {
		t.Errorf("For session '%v' expected user '%v' but '%v' was found", session, userToValidate, user)
	}
}

func AssertValidSession(t *testing.T, session *Session) {
	ok, err := tdb.CheckSession(session)
	if err != nil {
		t.Errorf("Cannot check session validness: %v", err)
	}
	if !ok {
		t.Errorf("Session ain't valid: '%v'", session)
	}
}

func CreateSession(t *testing.T, user User) Session {
	session, err := tdb.CreateSession(user)
	if err != nil {
		t.Errorf("Cannot create session for user '%v': %v", user, err)
	}
	if session.UserId != user.Id || session.Email != user.Email {
		t.Errorf("Create session '%v' does not match user '%v'", session, user)
	}
	return session
}

func AssertUsersInDb(t *testing.T, users ...User) {
	allUsers, err := tdb.Users()
	if err != nil {
		t.Errorf("Cannot get all users: %v", err)
	}
	if len(allUsers) != len(users) {
		t.Errorf("Expected to read 2 users from database")
	}
	for _, user := range users {
		assertContainsUser(t, allUsers, user)
	}
}

func AddTwoUsers(t *testing.T) (user1, user2 User) {
	const (
		user_1_name  = "user1"
		user_1_email = "user1email"
		user_1_pass  = "user1pass"
		user_1_role  = AdminRole
	)
	user1, err := tdb.AddUser(user_1_name, user_1_email, user_1_pass, user_1_role)
	if err != nil {
		t.Errorf("Cannot add user: %v", err)
	}

	const (
		user_2_name  = "user2"
		user_2_email = "user2email"
		user_2_pass  = "user2pass"
		user_2_role  = AdminRole
	)
	user2, err = tdb.AddUser(user_2_name, user_2_email, user_2_pass, user_2_role)
	if err != nil {
		t.Errorf("Cannot add user: %v", err)
	}
	return
}

func assertContainsUser(t *testing.T, users []User, userToFind User) {
	for _, user := range users {
		if user == userToFind {
			return
		}
	}
	t.Errorf("Does not contain user '%v'", userToFind)
}
