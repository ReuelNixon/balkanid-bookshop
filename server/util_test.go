package main

import (
	"bookshop/models"
	"bookshop/util"
	"testing"
)

func TestIsEmpty(t *testing.T) {
	tests := []struct {
		input string
		want  bool
	}{
		{"", false},
		{" ", true},
		{"  d", false},
		{"d", false},
	}
	for _, test := range tests {
		if got, _ := util.IsEmpty(test.input); got != test.want {
			t.Errorf("IsEmpty(%q) = %v", test.input, got)
		}
	}
}

func TestValidateRegister(t *testing.T) {
	tests := []struct {
		input      *models.User
		want       *models.UserErrors
		shouldPass bool
	}{
		{&models.User{Username: "", Email: "d", Password: "d3123214"}, &models.UserErrors{Err: true, Username: "Must not be empty", Email: "Must be a valid email", Password: "Length of password should be atleast 8 and it must be a combination of uppercase letters, lowercase letters and numbers"}, true},
		{&models.User{Username: "username", Email: "tester@gmail.com", Password: "Testing@123"}, &models.UserErrors{Err: false, Username: "Must not be empty", Email: "Must be a valid email", Password: "Length of password should be atleast 8 and it must be a combination of uppercase letters, lowercase letters and numbers"}, true},
	}

	for _, test := range tests {
		got := util.ValidateRegister(test.input)
		if test.shouldPass {
			if got.Err != test.want.Err {
				t.Errorf("ValidateRegister(%q) = %v", test.input, got)
			}
		} else {
			if got.Err == test.want.Err {
				t.Errorf("ValidateRegister(%q) = %v", test.input, got)
			}
		}
	}
}

func TestValidateAdminRegister(t *testing.T) {
	tests := []struct {
		input      *models.Admin
		want       *models.AdminErrors
		shouldPass bool
	}{
		{&models.Admin{Username: "", Email: "d", Password: "d3123214"}, &models.AdminErrors{Err: true, Username: "Must not be empty", Email: "Must be a valid email", Password: "Length of password should be atleast 8 and it must be a combination of uppercase letters, lowercase letters and numbers"}, true},
		{&models.Admin{Username: "username", Email: "tester@gmail.com", Password: "Testing@123"}, &models.AdminErrors{Err: false, Username: "Must not be empty", Email: "Must be a valid email", Password: "Length of password should be atleast 8 and it must be a combination of uppercase letters, lowercase letters and numbers"}, true},
	}

	for _, test := range tests {
		got := util.ValidateAdminRegister(test.input)
		if test.shouldPass {
			if got.Err != test.want.Err {
				t.Errorf("ValidateRegister(%q) = %v", test.input, got)
			}
		} else {
			if got.Err == test.want.Err {
				t.Errorf("ValidateRegister(%q) = %v", test.input, got)
			}
		}
	}
}
