package forms

type errors map[string][]string

func (e errors) Add(field, message string) {
	e[field] = append(e[field], message)
}

func (e errors) Get(field string) string {
	messages := e[field]

	if len(messages) == 0 {
		return ""
	}

	return messages[0]
}
