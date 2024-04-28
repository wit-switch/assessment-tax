package domain

type FieldMessage map[string]string

type FieldMessageList struct {
	f []FieldMessage
}

func NewFieldMessageList() *FieldMessageList {
	return &FieldMessageList{
		f: []FieldMessage{},
	}
}

func (l *FieldMessageList) Add(field, message string) *FieldMessageList {
	l.f = append(l.f, FieldMessage{
		"field":   field,
		"message": message,
	})

	return l
}

func (l *FieldMessageList) Value() []FieldMessage {
	return l.f
}
