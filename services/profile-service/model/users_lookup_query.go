package model

import "net/url"

type UsersLookupQuery struct {
	DelayMs string `form:"usersDelayMs"`
	Fail    string `form:"usersFail"`
}

func (q *UsersLookupQuery) SetDefaults() {
	if q.Fail == "" {
		q.Fail = "false"
	}
}

func (q UsersLookupQuery) ToUsersServiceQuery() url.Values {
	values := url.Values{}

	if q.DelayMs != "" {
		values.Set("delayMs", q.DelayMs)
	}
	if q.Fail == "true" {
		values.Set("fail", "true")
	}

	return values
}
