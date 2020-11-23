package models

type ResponderRequest struct {
	DirectMessage  bool
	MeString       string // string used to match an @
	Message        string // Message contents
	Service        string // Origin service
	ResponseTarget string // Where to direct response back to.
}
