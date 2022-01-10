package forms

import (
	"fmt"
	"net/url"
	"regexp"
	"strings"
	"unicode/utf8"
)

var EmailRX = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

type Form struct {
	url.Values
	Errors errors
}

func New(data url.Values) *Form {
	return &Form{
		data,
		errors(map[string][]string{}),
	}
}

func (f *Form) Required(fields ...string) {
	for _, field := range fields {
		value := f.Get(field)
		if strings.TrimSpace(value) == "" {
			f.Errors.Add(field, "This field cannot be blank")
		}
	}
}

func (f *Form) MaxLength(field string, maxLenght int) {
	value := f.Get(field)
	if strings.TrimSpace(value) == "" {
		return
	}

	if utf8.RuneCountInString(value) > maxLenght {
		f.Errors.Add(field, fmt.Sprintf("This field is too long (max length is %d)", maxLenght))
	}
}

func (f *Form) PermittedValues(field string, ops ...string) {
	value := f.Get(field)
	if strings.TrimSpace(value) == "" {
		return
	}

	for _, op := range ops {
		if value == op {
			return
		}
	}

	f.Errors.Add(field, "This field is invalid")
}

func (f *Form) MinLength(field string, minLength int) {
	value := f.Get(field)
	if strings.TrimSpace(value) == "" {
		return
	}

	if utf8.RuneCountInString(value) < minLength {
		f.Errors.Add(field, fmt.Sprintf("This field is too short (min length is %d)", minLength))
	}
}

func (f *Form) MatchesPattern(field string, pattern *regexp.Regexp) {
	value := f.Get(field)
	if strings.TrimSpace(value) == "" {
		return
	}

	if !pattern.MatchString(value) {
		f.Errors.Add(field, "This field is invalid")
	}
}

func (f *Form) Valid() bool {
	return len(f.Errors) == 0
}
