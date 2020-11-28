package models

import (
	"database/sql"
	"regexp"
	"time"
)

type ResponderMatcher struct {
	AlwaysRespond bool           `db:"always_respond",json:"always_respond"` // if false search for me string or DM = true.
	Enabled       bool           `db:"enabled",json:"enabled"`
	Description   string         `db:"description",json:"description"`
	MatcherString string         `db:"matcher_re",json:"matcher_re"`
	Response      string         `db:"repsonse",json:"repsonse"`
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
	start := time.Now()
	defer logger.Tracef("CountResponderMatchers() took %s", time.Since(start))

	var count uint64
	err := client.Get(&count, "SELECT count(id) FROM responder_matchers;")
	if err != nil {
		return 0, err
	}

	return count, nil
}

func ReadEnabledResponderMatchers() (*[]ResponderMatcher, error) {
	// Timing
	start := time.Now()
	defer logger.Tracef("ReadEnabledResponderMatchers() took %s", time.Since(start))

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
	start := time.Now()
	defer logger.Tracef("ReadEnabledResponderMatcher() took %s", time.Since(start))

	var rm ResponderMatcher
	err := client.Get(&rm, "SELECT * FROM responder_matchers WHERE id = $1;", id)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	return &rm, nil
}

func ReadResponderMatchersPage(index, count uint) (*[]ResponderMatcher, error) {
	// Timing
	start := time.Now()
	defer logger.Tracef("ReadResponderMatchersPage(%d, %d) took %s",index, count, time.Since(start))

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
