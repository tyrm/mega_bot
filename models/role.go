package models

import (
	"database/sql"
	"github.com/lib/pq"
	"time"
)

type Role struct {
	Name        string         `db:"name",json:"name"`
	Description sql.NullString `db:"description",json:"description"`
	SystemRole  bool           `db:"system_role",json:"system_role"`

	// metadata
	ID        string    `db:"id",json:"id"`
	CreatedAt time.Time `db:"created_at",json:"created_at"`
	UpdatedAt time.Time `db:"updated_at",json:"updated_at"`
}

func CountRolesByUserIDRoles(user_id string, roles []string) (int, error) {
	// Timing
	defer stats.NewTiming().Send("CountRolesByUserIDRoles")

	var count int
	err := client.Get(&count, "SELECT count(u.id) FROM users as u " +
		"INNER JOIN user_roles as ur ON (u.id = ur.user_id) " +
		"INNER JOIN roles as r ON (ur.role_id = r.id) " +
		"WHERE u.id = $1 AND r.name = ANY($2);", user_id, pq.Array(roles))
	if err != nil {
		return 0, err
	}

	return count, nil
}
