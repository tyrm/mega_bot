package models

import (
	"database/sql"
	"regexp"
	"time"
)

type ResponderMatcher struct {
	AlwaysRespond bool           `db:"always_respond",json:"always_respond"` // if false search for me string or DM = true.
	Description   string         `db:"description",json:"description"`
	Enabled       bool           `db:"enabled",json:"enabled"`
	MatcherString string         `db:"matcher_re",json:"matcher_re"`
	Response      string         `db:"response",json:"response"`
	MatcherRE     *regexp.Regexp `json:"-"` // regex to match message

	// metadata
	ID        string    `db:"id",json:"id"`
	CreatedAt time.Time `db:"created_at",json:"created_at"`
	UpdatedAt time.Time `db:"updated_at",json:"updated_at"`
}

func (rm *ResponderMatcher) BuildRE() error {
	re, err := regexp.Compile(rm.MatcherString)
	if err != nil {
		return err
	}

	rm.MatcherRE = re
	return nil
}

func CountResponderMatchers() (uint64, error) {
	// Timing
	defer stats.NewTiming().Send("CountResponderMatchers")

	var count uint64
	err := client.Get(&count, "SELECT count(id) FROM responder_matchers;")
	if err != nil {
		return 0, err
	}

	return count, nil
}

func CreateResponderMatcher(rm *ResponderMatcher) error {
	// Timing
	defer stats.NewTiming().Send("CreateResponderMatcher")

	err := client.
		QueryRowx(`INSERT INTO public.responder_matchers(always_respond, description, enabled, matcher_re, response) 
        	VALUES ($1, $2, $3, $4, $5) RETURNING id, created_at, updated_at;`, rm.AlwaysRespond, rm.Description, rm.Enabled, rm.MatcherString, rm.Response).
		Scan(&rm.ID, &rm.CreatedAt, &rm.UpdatedAt)
	return err
}

func ReadEnabledResponderMatchers() (*[]ResponderMatcher, error) {
	// Timing
	defer stats.NewTiming().Send("ReadEnabledResponderMatchers")

	var rms []*ResponderMatcher
	err := client.Select(&rms, "SELECT * FROM responder_matchers WHERE enabled = TRUE ORDER BY id;")
	if err != nil {
		return nil, err
	}

	var rmsResponse []ResponderMatcher
	for _, rm := range rms {
		err := rm.BuildRE()
		if err != nil {
			logger.Warningf("bad regex %s: %s", rm.MatcherString, err.Error())
		}
		rmsResponse = append(rmsResponse, *rm)
	}

	return &rmsResponse, nil
}

func ReadResponderMatcher(id string) (*ResponderMatcher, error) {
	// Timing
	defer stats.NewTiming().Send("ReadResponderMatcher")

	var rm ResponderMatcher
	err := client.Get(&rm, "SELECT * FROM responder_matchers WHERE id = $1;", id)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	return &rm, nil
}

func ReadResponderMatchersPage(index, count int) (*[]ResponderMatcher, error) {
	// Timing
	defer stats.NewTiming().Send("ReadResponderMatchersPage")

	var rms []*ResponderMatcher

	offset := index * count
	err := client.Select(&rms, "SELECT * FROM responder_matchers ORDER BY description LIMIT $1 OFFSET $2;",
		count, offset)
	if err != nil {
		return nil, err
	}

	var rmsResponse []ResponderMatcher
	for _, rm := range rms {
		err := rm.BuildRE()
		if err != nil {
			logger.Warningf("bad regex %s: %s", rm.MatcherString, err.Error())
		}
		rmsResponse = append(rmsResponse, *rm)
	}

	return &rmsResponse, nil
}

func UpdateResponderMatcher(rm *ResponderMatcher) error {
	// Timing
	defer stats.NewTiming().Send("UpdateResponderMatcher")

	err := client.
		QueryRowx(`UPDATE public.responder_matchers
			SET always_respond=$1, description=$2, enabled=$3, matcher_re=$4, response=$5, updated_at=CURRENT_TIMESTAMP
			WHERE id=$6 RETURNING created_at, updated_at;`, rm.AlwaysRespond, rm.Description, rm.Enabled,
			rm.MatcherString, rm.Response, rm.ID).Scan(&rm.CreatedAt, &rm.UpdatedAt)
	return err
}