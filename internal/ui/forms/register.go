package forms

import (
	"errors"

	"fyne.io/fyne/v2/data/binding"
	log "github.com/sirupsen/logrus"
)

type RegisterForm struct {
	Username        binding.String
	Password        binding.String
	ConfirmPassword binding.String
}

func NewRegisterForm() *RegisterForm {
	return &RegisterForm{
		Username:        binding.NewString(),
		Password:        binding.NewString(),
		ConfirmPassword: binding.NewString(),
	}
}

func (f *RegisterForm) Validate() error {
	_, err := f.Username.Get()
	if err != nil {
		log.Errorf("failed to get username: %v", err)
		return errors.New("username is required")
	}

	password, err := f.Password.Get()
	if err != nil {
		log.Errorf("failed to get password: %v", err)
		return errors.New("password is required")
	}

	confirmPassword, err := f.ConfirmPassword.Get()
	if err != nil {
		log.Errorf("failed to get confirm password: %v", err)
		return errors.New("confirm password is required")
	}

	if password != confirmPassword {
		log.Error("password and confirm password do not match")
		return errors.New("password and confirm password do not match")
	}
	return nil
}

func (f *RegisterForm) GetData() (string, string, string, error) {
	err := f.Validate()
	if err != nil {
		log.Errorf("form validation failed: %v", err)
		return "", "", "", err
	}
	username, _ := f.Username.Get()
	password, _ := f.Password.Get()
	confirmPassword, _ := f.ConfirmPassword.Get()
	return username, password, confirmPassword, nil
}
