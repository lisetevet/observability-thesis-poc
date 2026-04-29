package model

import "net/url"

type ProfileLookupQuery struct {
	UsersDelayMs   string
	UsersFail      string
	ProfileDelayMs string
	ProfileFail    string
}

func (q *ProfileLookupQuery) SetDefaults() {
	if q.UsersFail == "" {
		q.UsersFail = "false"
	}
	if q.ProfileFail == "" {
		q.ProfileFail = "false"
	}
}

func (q ProfileLookupQuery) ToProfileServiceQuery() url.Values {
	values := url.Values{}

	if q.UsersDelayMs != "" {
		values.Set("usersDelayMs", q.UsersDelayMs)
	}
	if q.UsersFail == "true" {
		values.Set("usersFail", "true")
	}
	if q.ProfileDelayMs != "" {
		values.Set("delayMs", q.ProfileDelayMs)
	}
	if q.ProfileFail == "true" {
		values.Set("fail", "true")
	}

	return values
}