package erix

import "strings"

type MultiErr struct {
	errs []error
}

func (m *MultiErr) Error() string {
	if len(m.errs) == 0 {
		return ""
	}

	sb := strings.Builder{}

	_, _ = sb.WriteString(m.errs[0].Error())

	for _, err := range m.errs[1:] {
		_, _ = sb.WriteString("; ")
		_, _ = sb.WriteString(err.Error())
	}

	return sb.String()
}

func (m *MultiErr) Unwrap() []error {
	return m.errs
}

func NewMulti(errs []error) error {
	if len(errs) == 0 {
		return nil
	}

	return &MultiErr{errs: errs}
}
