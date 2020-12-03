package models

import (
	"database/sql"
	"time"
)

type User struct {
	Email    string         `db:"email",json:"email"`
	Password sql.NullString `db:"password",json:"password"`
	Nick     sql.NullString `db:"nick",json:"nick"`

	// metadata
	ID        string    `db:"id",json:"id"`
	CreatedAt time.Time `db:"created_at",json:"created_at"`
	UpdatedAt time.Time `db:"updated_at",json:"updated_at"`
}

func (u *User) HasOneOfRoles(roles []string) (bool, error) {
	roleCount, err := CountRolesByUserIDRoles(u.ID, roles)
	if err != nil {
		return false, err
	}

	if roleCount == 0 {
		return false, nil
	}

	return true, nil
}

func CreateUser(u *User) error {
	// Timing
	defer stats.NewTiming().Send("CreateUser")

	err := client.
		QueryRowx(`INSERT INTO public.users(email, password, nick) 
        	VALUES ($1, $2, $3) RETURNING id, created_at, updated_at;`, u.Email, u.Password, u.Nick).
		Scan(&u.ID, &u.CreatedAt, &u.UpdatedAt)
	return err
}

func ReadUser(id string) (*User, error) {
	// Timing
	defer stats.NewTiming().Send("ReadUser")

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
	// Timing
	defer stats.NewTiming().Send("ReadUserByEmail")

	var user User

	err := client.Get(&user, "SELECT * FROM users WHERE email = $1;", e)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	return &user, nil
}
