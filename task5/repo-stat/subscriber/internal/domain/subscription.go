package domain

import "errors"

type Subscription struct {
	ID    int32
	Owner string
	Repo  string
}

func (sub Subscription) IsValid() bool {
	if sub.Owner == "" || sub.Repo == "" {
		return false
	}
	return true
}

var ErrDuplicateSubscription = errors.New("subscription already exists")
