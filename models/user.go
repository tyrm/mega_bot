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

type UserListItem struct {
	User

	Administrator sql.NullString `db:"role_administrator",json:"role_administrator"`
	Operator      sql.NullString `db:"role_operator",json:"role_operator"`
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

func ReadUserListPage(index, count int) (*[]UserListItem, error) {
	// Timing
	defer stats.NewTiming().Send("ReadUserListPage")

	var userlist []UserListItem

	offset := index * count
	err := client.Select(&userlist, "SELECT u.id, u.email, u.nick, r_a.name as role_administrator, "+
		"r_o.name as role_operator FROM users as u "+
		"INNER JOIN user_roles as ur ON (u.id = ur.user_id) "+
		"LEFT JOIN roles as r_a ON (ur.role_id = r_a.id AND r_a.name = 'administrator') "+
		"LEFT JOIN roles as r_o ON (ur.role_id = r_o.id AND r_o.name = 'operator') "+
		"ORDER BY u.email LIMIT $1 OFFSET $2;",
		count, offset)
	if err != nil {
		return nil, err
	}

	return &userlist, nil
}
