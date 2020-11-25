package models

import (
	"database/sql"
	"time"
)

type User struct {
	Email      string `db:"email",json:"email"`
	Password   string `db:"password",json:"password"`
	Nick       string `db:"nick",json:"nick"`
	Authorized bool   `db:"authorized",json:"authorized"`
	Admin      bool   `db:"admin",json:"admin"`

	// metadata
	ID        string    `db:"id",json:"id"`
	CreatedAt time.Time `db:"created_at",json:"created_at"`
	UpdatedAt time.Time `db:"updated_at",json:"updated_at"`
}

func CreateUser(u *User) error {
	// Timing
	start := time.Now()
	defer logger.Tracef("CreateUser() took %s", time.Since(start))

	err := client.
		QueryRowx(`INSERT INTO public.users(email, password, nick, authorized, admin) 
			VALUES (:email, :password, :nick, :authorized, :admin) RETURNING id, created_at, updated_at;`, &u).
		Scan(&u.ID, &u.CreatedAt, &u.UpdatedAt)
	return err
}

func ReadUser(id string) (*User, error) {
	var user User

	err := client.Get(&user, "SELECT * FROM users WHERE id = $1;", id)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	return &user, nil
}

func ReadUserByEmail(e string) (*User, error) {
	var user User

	err := client.Get(&user, "SELECT * FROM users WHERE email = $1;", e)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	return &user, nil
}