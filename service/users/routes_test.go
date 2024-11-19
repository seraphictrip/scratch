package users_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"scratch/service/users"
	"scratch/types"
	"testing"

	"github.com/gorilla/mux"
)

var HandleRegisterTests = []struct {
	name string
	// input
	payload types.RegisterUserPayload

	// output
	statusCode int
	body       any
	// mock
	userStore types.UserStore
}{
	{"BadRequest-email",
		types.RegisterUserPayload{
			FirstName: "Stephen",
			LastName:  "King",
			Email:     "",
			Password:  "asd",
		},
		http.StatusBadRequest,
		nil,
		&mockUserStore{}},
	{"BadRequest-missing-firstname",
		types.RegisterUserPayload{
			FirstName: "",
			LastName:  "King",
			Email:     "stepking@example.com",
			Password:  "asd",
		},
		http.StatusBadRequest,
		nil,
		&mockUserStore{}},
	{"BadRequest-lastName",
		types.RegisterUserPayload{
			FirstName: "Stephen",
			LastName:  "",
			Email:     "tepking@example.com",
			Password:  "asd",
		},
		http.StatusBadRequest,
		nil,
		&mockUserStore{}},
	{"BadRequest-password",
		types.RegisterUserPayload{
			FirstName: "Stephen",
			LastName:  "King",
			Email:     "tepking@example.com",
			Password:  "",
		},
		http.StatusBadRequest,
		nil,
		&mockUserStore{}},
	{"Success",
		types.RegisterUserPayload{
			FirstName: "Stephen",
			LastName:  "King",
			Email:     "tepking@example.com",
			Password:  "password",
		},
		http.StatusCreated,
		nil,
		&mockUserStore{}},
}

func TestHandleRegister(t *testing.T) {
	for _, e := range HandleRegisterTests {
		e := e
		t.Run(e.name, func(t *testing.T) {
			// code
			marshalled, _ := json.Marshal(e.payload)
			handler := users.NewHandler(e.userStore)
			request, err := http.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(marshalled))
			if err != nil {
				t.Fatalf("failed to instantiate test: %v", err)
			}

			rr := httptest.NewRecorder()
			router := mux.NewRouter()
			router.HandleFunc("/register", handler.HandleRegister)
			router.ServeHTTP(rr, request)
			if rr.Result().StatusCode != e.statusCode {
				t.Fatalf("StatusCode: %v, want %v", rr.Result().StatusCode, e.statusCode)
			}

		})
	}
}

type mockUserStore struct{}

func (s mockUserStore) GetUserByEmail(email string) (*types.User, error) {
	return nil, users.ErrUserExists
}

func (s mockUserStore) GetUserByID(id int) (*types.User, error) {
	panic("not implemented") // TODO: Implement
}

func (s mockUserStore) CreateUser(user types.User) error {
	return nil
}
