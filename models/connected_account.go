package models

import (
	"database/sql"
	"time"
)

type ConnectedAccount struct {
	Provider       string `db:"provider",json:"provider"`
	ProviderUserID string `db:"provider_user_id",json:"provider_user_id"`
	UserID         string `db:"user_id",json:"-"`

	// metadata
	ID        string    `db:"id",json:"id"`
	CreatedAt time.Time `db:"created_at",json:"created_at"`
	UpdatedAt time.Time `db:"updated_at",json:"updated_at"`
}

func CreateConnectedAccount(ca *ConnectedAccount) error {
	// Timing
	start := time.Now()
	defer logger.Tracef("CreateConnectedAccount() took %s", time.Since(start))

	err := client.
		QueryRowx(`INSERT INTO public.connected_accounts(provider, provider_user_id, user_id) VALUES (:provider, :provider_user_id, :user_id) RETURNING id, created_at, updated_at;`, &ca).
		Scan(&ca.ID, &ca.CreatedAt, &ca.UpdatedAt)
	return err
}

func ReadConnectedAccount(provider, providerUserID string) (*ConnectedAccount, error){
	// Timing
	start := time.Now()
	defer logger.Tracef("ReadConnectedAccount() took %s", time.Since(start))

	var ca ConnectedAccount

	err := client.Get(&ca, `SELECT * FROM connected_accounts WHERE provider=$1 AND provider_user_id=$2;`, provider, providerUserID)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	return &ca, nil
}