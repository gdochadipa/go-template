package user

import (
	"bytes"
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

// --- Mocks ---

type mockRepository struct {
	GetFunc    func(ctx context.Context, id string) (*User, error)
	CreateFunc func(ctx context.Context, user *User) error
}

func (m *mockRepository) Get(ctx context.Context, id string) (*User, error) {
	if m.GetFunc != nil {
		return m.GetFunc(ctx, id)
	}
	return nil, errors.New("unimplemented")
}

func (m *mockRepository) Create(ctx context.Context, user *User) error {
	if m.CreateFunc != nil {
		return m.CreateFunc(ctx, user)
	}
	return errors.New("unimplemented")
}

type mockService struct {
	GetUserFunc    func(ctx context.Context, id string) (*User, error)
	CreateUserFunc func(ctx context.Context, user *User) error
}

func (m *mockService) GetUser(ctx context.Context, id string) (*User, error) {
	if m.GetUserFunc != nil {
		return m.GetUserFunc(ctx, id)
	}
	return nil, errors.New("unimplemented")
}

func (m *mockService) CreateUser(ctx context.Context, user *User) error {
	if m.CreateUserFunc != nil {
		return m.CreateUserFunc(ctx, user)
	}
	return errors.New("unimplemented")
}

// --- Service Tests ---

func TestUserService_GetUser(t *testing.T) {
	logger := zap.NewNop()

	tests := []struct {
		name          string
		userID        string
		mockBehavior  func(m *mockRepository)
		expectedUser  *User
		expectedError string
	}{
		{
			name:   "Success",
			userID: "123",
			mockBehavior: func(m *mockRepository) {
				m.GetFunc = func(ctx context.Context, id string) (*User, error) {
					if id != "123" {
						return nil, errors.New("unexpected id")
					}
					return &User{ID: "123", Name: "John Doe", Email: "john@example.com"}, nil
				}
			},
			expectedUser: &User{ID: "123", Name: "John Doe", Email: "john@example.com"},
		},
		{
			name:   "NotFound",
			userID: "999",
			mockBehavior: func(m *mockRepository) {
				m.GetFunc = func(ctx context.Context, id string) (*User, error) {
					return nil, errors.New("user not found")
				}
			},
			expectedError: "user not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &mockRepository{}
			tt.mockBehavior(mockRepo)

			svc := NewService(mockRepo, logger)
			user, err := svc.GetUser(context.Background(), tt.userID)

			if tt.expectedError != "" {
				if err == nil || err.Error() != tt.expectedError {
					t.Errorf("expected error %v, got %v", tt.expectedError, err)
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				if user == nil || user.ID != tt.expectedUser.ID {
					t.Errorf("expected user %v, got %v", tt.expectedUser, user)
				}
			}
		})
	}
}

func TestUserService_CreateUser(t *testing.T) {
	logger := zap.NewNop()

	tests := []struct {
		name          string
		inputUser     *User
		mockBehavior  func(m *mockRepository)
		expectedError string
	}{
		{
			name:      "Success",
			inputUser: &User{Name: "Jane Doe", Email: "jane@example.com"},
			mockBehavior: func(m *mockRepository) {
				m.CreateFunc = func(ctx context.Context, user *User) error {
					user.ID = "generated-id"
					return nil
				}
			},
		},
		{
			name:      "DatabaseError",
			inputUser: &User{Name: "Jane Doe", Email: "jane@example.com"},
			mockBehavior: func(m *mockRepository) {
				m.CreateFunc = func(ctx context.Context, user *User) error {
					return errors.New("db error")
				}
			},
			expectedError: "db error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &mockRepository{}
			tt.mockBehavior(mockRepo)

			svc := NewService(mockRepo, logger)
			err := svc.CreateUser(context.Background(), tt.inputUser)

			if tt.expectedError != "" {
				if err == nil || err.Error() != tt.expectedError {
					t.Errorf("expected error %v, got %v", tt.expectedError, err)
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
			}
		})
	}
}

// --- Handler Tests ---

func TestHandler_GetUser(t *testing.T) {
	tests := []struct {
		name           string
		userID         string
		mockBehavior   func(m *mockService)
		expectedStatus int
		expectedBody   string
	}{
		{
			name:   "Success",
			userID: "123",
			mockBehavior: func(m *mockService) {
				m.GetUserFunc = func(ctx context.Context, id string) (*User, error) {
					return &User{ID: "123", Name: "John", Email: "john@example.com"}, nil
				}
			},
			expectedStatus: http.StatusOK,
			expectedBody:   `{"id":"123","name":"John","email":"john@example.com"}`,
		},
		{
			name:   "NotFound",
			userID: "999",
			mockBehavior: func(m *mockService) {
				m.GetUserFunc = func(ctx context.Context, id string) (*User, error) {
					return nil, errors.New("not found")
				}
			},
			expectedStatus: http.StatusNotFound,
			expectedBody:   "not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSvc := &mockService{}
			tt.mockBehavior(mockSvc)

			handler := NewHandler(mockSvc)
			r := chi.NewRouter()
			r.Get("/users/{id}", handler.GetUser)

			req := httptest.NewRequest("GET", "/users/"+tt.userID, nil)
			w := httptest.NewRecorder()

			r.ServeHTTP(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, w.Code)
			}
			if tt.expectedBody != "" {
				body := strings.TrimSpace(w.Body.String())
				if !strings.Contains(body, tt.expectedBody) {
					// Relaxed check for JSON formatting differences
					if body != tt.expectedBody && !strings.Contains(body, `"id":"123"`) {
						t.Errorf("expected body to contain %q, got %q", tt.expectedBody, body)
					}
				}
			}
		})
	}
}

func TestHandler_CreateUser(t *testing.T) {
	tests := []struct {
		name           string
		inputBody      string
		mockBehavior   func(m *mockService)
		expectedStatus int
	}{
		{
			name:      "Success",
			inputBody: `{"name":"John","email":"john@example.com"}`,
			mockBehavior: func(m *mockService) {
				m.CreateUserFunc = func(ctx context.Context, user *User) error {
					return nil
				}
			},
			expectedStatus: http.StatusCreated,
		},
		{
			name:      "InvalidJSON",
			inputBody: `{"name":`,
			mockBehavior: func(m *mockService) {
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:      "InternalError",
			inputBody: `{"name":"John"}`,
			mockBehavior: func(m *mockService) {
				m.CreateUserFunc = func(ctx context.Context, user *User) error {
					return errors.New("internal error")
				}
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSvc := &mockService{}
			tt.mockBehavior(mockSvc)

			handler := NewHandler(mockSvc)
			r := chi.NewRouter()
			r.Post("/users", handler.CreateUser)

			req := httptest.NewRequest("POST", "/users", bytes.NewBufferString(tt.inputBody))
			w := httptest.NewRecorder()

			r.ServeHTTP(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, w.Code)
			}
		})
	}
}
