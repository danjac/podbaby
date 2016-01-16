package decoders

import "testing"

func TestNewChannelIfURLEmpty(t *testing.T) {
	d := &NewChannel{}
	if err := d.Decode(); err == nil {
		t.Fatal("Channel should have URL")
	}
}

func TestNewChannelIfURLInvalid(t *testing.T) {
	d := &NewChannel{
		URL: "testing!",
	}
	if err := d.Decode(); err == nil {
		t.Fatal("Channel should have a valid URL")
	}
}

func TestNewChannelIfURLValid(t *testing.T) {
	d := &NewChannel{
		URL: "http://google.com",
	}
	if err := d.Decode(); err != nil {
		t.Fatal("Channel should have a valid URL")
	}
}

func TestNewEmailIfEmpty(t *testing.T) {
	d := &NewEmail{}
	if err := d.Decode(); err == nil {
		t.Fatal("Should be valid")
	}
}

func TestNewEmailTrimmed(t *testing.T) {
	d := &NewEmail{
		Email: "    TEST@gmail.com",
	}
	if err := d.Decode(); err != nil {
		t.Fatal("Should be valid")
	}
	if d.Email != "test@gmail.com" {
		t.Fatal("Email should be lowercase and trimmed")
	}
}

func TestNewPasswordIfEmpty(t *testing.T) {
	d := &NewPassword{}
	if err := d.Decode(); err == nil {
		t.Fatal("Should have error if no password fields")
	}
}

func TestNewPasswordIfOldPasswordMissing(t *testing.T) {
	d := &NewPassword{
		NewPassword: "testpass",
	}
	if err := d.Decode(); err == nil {
		t.Fatal("Should have error if no password fields")
	}
}

func TestNewPasswordIfNewPasswordTooShort(t *testing.T) {
	d := &NewPassword{
		OldPassword: "test",
		NewPassword: "test",
	}
	if err := d.Decode(); err == nil {
		t.Fatal("Should have error if new password too short")
	}
}

func TestNewPasswordIfValid(t *testing.T) {
	d := &NewPassword{
		OldPassword: "test",
		NewPassword: "testpass",
	}
	if err := d.Decode(); err != nil {
		t.Fatal("Should have error if new password too short")
	}
}

func TestRecoverPasswordIfNoIdentifier(t *testing.T) {
	d := &RecoverPassword{}
	if err := d.Decode(); err == nil {
		t.Fatal("Should return an error if no identifier")
	}
}

func TestRecoverPasswordIfIdentifier(t *testing.T) {
	d := &RecoverPassword{"tester"}
	if err := d.Decode(); err != nil {
		t.Fatal("Should pass if identifier")
	}
}

func TestSignupIfEmpty(t *testing.T) {
	d := &Signup{}
	if err := d.Decode(); err == nil {
		t.Fatal("Should return an error if no fields")
	}
}

func TestSignupIfInvalidEmail(t *testing.T) {
	d := &Signup{
		Name:     "tester",
		Email:    "test",
		Password: "testpass",
	}
	if err := d.Decode(); err == nil {
		t.Fatal("Should return an error if email invalid")
	}
}

func TestSignupIfInvalidPassword(t *testing.T) {
	d := &Signup{
		Name:     "tester",
		Email:    "test@gmail.com",
		Password: "test",
	}
	if err := d.Decode(); err == nil {
		t.Fatal("Should return an error if password too short")
	}
}

func TestSignupTrimsNameAndEmail(t *testing.T) {
	d := &Signup{
		Name:     "tester   ",
		Email:    "    TEST@gmail.com",
		Password: "testpass",
	}
	if err := d.Decode(); err != nil {
		t.Fatal("Should validate signup")
	}
	if d.Name != "tester" {
		t.Fatal("Name should be trimmed")
	}
	if d.Email != "test@gmail.com" {
		t.Fatal("Email should be lowercase and trimmed")
	}
}
