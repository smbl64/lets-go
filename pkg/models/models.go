package models

import (
	"errors"
	"time"
)

var ErrNoRecord = errors.New("models: No matching record has been found")

type Snippet struct {
	ID      int
	Title   string
	Content string
	Created time.Time
	Expires time.Time
}
