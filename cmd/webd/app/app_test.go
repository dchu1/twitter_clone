package app

import (
	"fmt"
	"testing"
)

func TestAddUser(t *testing.T) {
	app := MakeApp()
	expected := MakeUser("TestFirst", "TestLast", "test@test.com", "testpass", 0)
	app.AddUser("TestFirst", "TestLast", "test@test.com", "testpass")
	actual := app.getUsers()[0]
	if expected.Email != actual.Email ||
		expected.FirstName != actual.FirstName ||
		expected.LastName != actual.LastName ||
		expected.id != actual.id {
		t.Error(fmt.Sprintf("Test Failed: %s, %s, %s, %d", actual.Email, actual.FirstName, actual.LastName, actual.id))
	}
}
