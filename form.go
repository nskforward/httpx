package httpx

import (
	"fmt"
	"strings"
)

type validateStrategy struct {
	required  bool
	minLength int
	maxLength int
	lowercase bool
	uppercase bool
}

func (s validateStrategy) Validate(f, v string) (string, error) {
	if s.required && len(v) == 0 {
		return "", fmt.Errorf("field '%s' cannot be empty", f)
	}
	if len(v) < s.minLength {
		return "", fmt.Errorf("field '%s' is too short (min %d characters)", f, s.minLength)
	}
	if s.maxLength > 0 && len(v) > s.maxLength {
		return "", fmt.Errorf("field '%s' is too large (max %d characters)", f, s.maxLength)
	}
	if s.lowercase {
		v = strings.ToLower(v)
	}
	if s.uppercase {
		v = strings.ToUpper(v)
	}
	return v, nil
}

type ValidateOpt func(*validateStrategy)

func (ctx *Ctx) FormParam(field string, opts ...ValidateOpt) (string, error) {
	if !ctx.formParsed {
		err := ctx.Request().ParseMultipartForm(1024 * 1024)
		if err != nil {
			panic(err)
		}
		ctx.formParsed = true
	}
	var s validateStrategy
	for _, opt := range opts {
		opt(&s)
	}
	return s.Validate(field, ctx.Request().FormValue(field))
}

func Required(s *validateStrategy) {
	s.required = true
}

func MinLength(n int) ValidateOpt {
	return func(s *validateStrategy) {
		s.minLength = n
	}
}

func MaxLength(n int) ValidateOpt {
	return func(s *validateStrategy) {
		s.maxLength = n
	}
}

func Lowercase(s *validateStrategy) {
	s.lowercase = true
}

func Uppercase(s *validateStrategy) {
	s.uppercase = true
}
