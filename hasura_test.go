package hjwt

import (
	"testing"
	"time"
)

type TestUser struct {
	Id    string   `json:"id"`
	Name  string   `json:"name"`
	Email string   `json:"email"`
	Roles []string `json:"roles"`
}

func TestValidate(t *testing.T) {
	token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1aWQiOiJtdW5nY2hhdSIsInJvbGUiOiJhZG1pbiIsInVzZXIiOnsiaWQiOiJtdW5nY2hhdSIsIm5hbWUiOiJEcmFjbyBDaGF1IiwiZW1haWwiOiJkY0BhYmMuY29tIiwicm9sZXMiOlsiYWRtaW4iXX0sImh0dHBzOi8vaGFzdXJhLmlvL2p3dC9jbGFpbXMiOnsieC1oYXN1cmEtYWxsb3dlZC1yb2xlcyI6WyJhZG1pbiJdLCJ4LWhhc3VyYS1kZWZhdWx0LXJvbGUiOiJhZG1pbiIsIngtaGFzdXJhLXVzZXItaGFsbHMiOlsiU1IwMSIsIlNSMDIiXSwieC1oYXN1cmEtdXNlci1pZCI6Im11bmdjaGF1IiwieC1oYXN1cmEtdXNlci1yb2xlcyI6ImFkbWluIiwieC1oYXN1cmEtdXNlci1zdXBlciI6dHJ1ZX0sImV4cCI6MTU4NzY5OTM3MiwiaWF0IjoxNTg3Njk1NzEyLCJpc3MiOiJIQVNVUkEtSldULUdFTkVSQVRPUiJ9.Mkz0vleXzuDlc4UkOjz4pTHDZBe2FcF_Z4IjSQ9Z-og"
	err := Validate(token)
	if err != nil {
		t.Errorf("error while validate token: %v", err)
	}
}

func TestGenerate(t *testing.T) {
	u := &TestUser{
		Id:    "mungchau",
		Name:  "Draco Chau",
		Email: "dc@abc.com",
		Roles: []string{"admin"},
	}
	extra := make(map[string]interface{})
	extra["halls"] = []string{"SR01", "SR02"}
	extra["super"] = true
	token, err := Generate(u.Id, u.Roles, u, extra, time.Hour)
	if err != nil {
		t.Errorf("error while generate token: %v", err)
	}
	t.Log(token)
}
