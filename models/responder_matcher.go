package models

import (
	"regexp"
	"time"
)

type ResponderMatcher struct {
	AlwaysRespond bool           `db:"always_respond",json:"always_respond"` // if false search for me string or DM = true.
	Enabled       bool           `db:"enabled",json:"enabled"`
	MatcherString string         `db:"matcher_re",json:"matcher_re"`
	Response      string         `db:"repsonse",json:"repsonse"`
	MatcherRE     *regexp.Regexp `json:"-"` // regex to match message

	// metadata
	ID        string      `db:"id",json:"id"`
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

func ReadEnabledResponderMatchers() ([]*ResponderMatcher, error) {
	start := time.Now()
	var rms []*ResponderMatcher

	err := client.Select(&rms, "SELECT * FROM responder_matchers WHERE enabled = TRUE ORDER BY id;")
	if err != nil {
		return nil, err
	}

	for _, rm := range rms {
		err := rm.BuildRE()
		if err != nil {
			logger.Warningf("bad regex %s: %s", rm.MatcherString, err.Error())
		}
	}

	duration := time.Since(start)
	logger.Tracef("ReadEnabledResponderMatchers() took %s", duration)
	return rms, nil
}
