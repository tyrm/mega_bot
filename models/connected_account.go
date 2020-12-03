package models

import (
	"database/sql"
	"time"
)

type ConnectedAccount struct {
	UserID           string         `db:"user_id",json:"user_id"`
	Provider         string         `db:"provider",json:"provider"`
	ProviderUserID   string         `db:"provider_user_id",json:"provider_user_id"`
	ProviderUsername sql.NullString `db:"provider_username",json:"provider_username"`
	ProviderAvatar   sql.NullString `db:"provider_avatar",json:"provider_avatar"`

	// metadata
	ID        string    `db:"id",json:"id"`
	CreatedAt time.Time `db:"created_at",json:"created_at"`
	UpdatedAt time.Time `db:"updated_at",json:"updated_at"`
}

func CreateConnectedAccount(ca *ConnectedAccount) error {
	// Timing
	defer stats.NewTiming().Send("CreateConnectedAccount")

	err := client.
		QueryRowx(`INSERT INTO public.connected_accounts(user_id, provider, provider_user_id, provider_username, 
 			provider_avatar) VALUES ($1, $2, $3, $4, $5) RETURNING id, created_at, updated_at;`,
			ca.UserID, ca.Provider, ca.ProviderUserID, ca.ProviderUsername, ca.ProviderAvatar).
		Scan(&ca.ID, &ca.CreatedAt, &ca.UpdatedAt)
	return err
}

func ReadConnectedAccount(provider, providerUserID string) (*ConnectedAccount, error) {
	// Timing
	defer stats.NewTiming().Send("ReadConnectedAccount")

	var ca ConnectedAccount

	err := client.Get(&ca, `SELECT * FROM connected_accounts WHERE provider=$1 AND provider_user_id=$2;`, provider, providerUserID)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	return &ca, nil
}
