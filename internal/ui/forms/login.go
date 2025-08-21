package forms

import (
	"fmt"

	"fyne.io/fyne/v2/data/binding"
)

type LoginForm struct {
	Username binding.String
	Password binding.String
}

func NewLoginForm() *LoginForm {
	return &LoginForm{
		Username: binding.NewString(),
		Password: binding.NewString(),
	}
}

func (f *LoginForm) GetData() (string, string, error) {
	if err := f.Validate(); err != nil {
		return "", "", err
	}

	username, _ := f.Username.Get()
	password, _ := f.Password.Get()
	return username, password, nil
}

func (f *LoginForm) Validate() error {
	if username, _ := f.Username.Get(); username == "" {
		return fmt.Errorf("username cannot be empty")
	}
	if password, _ := f.Password.Get(); password == "" {
		return fmt.Errorf("password cannot be empty")
	}
	return nil
}
