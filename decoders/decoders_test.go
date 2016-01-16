package decoders

import "testing"

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
