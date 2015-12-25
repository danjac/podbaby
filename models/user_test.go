package models

import "testing"

func TestSetPassword(t *testing.T) {
	user := &User{}
	user.SetPassword("testpass")
	if user.Password == "testpass" {
		t.Fail()
	}
}

func TestCheckPasswordIfEmpty(t *testing.T) {
	user := &User{}
	if user.CheckPassword("testpass") {
		t.Fail()
	}
}

func TestCheckPasswordIfWrong(t *testing.T) {
	user := &User{}
	user.SetPassword("TestPas5")
	if user.CheckPassword("testpass") {
		t.Fail()
	}
}

func TestCheckPasswordIfCorrect(t *testing.T) {
	user := &User{}
	user.SetPassword("testpass")
	if !user.CheckPassword("testpass") {
		t.Fail()
	}
}
