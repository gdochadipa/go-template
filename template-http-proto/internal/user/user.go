package user

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	userv1 "github.com/user/go-templates/template-http-proto/gen/go/user/v1"
	"go.uber.org/zap"
	"google.golang.org/protobuf/encoding/protojson"
)

// --- Domain/Service ---

type Service interface {
	GetUser(ctx context.Context, id string) (*userv1.User, error)
	CreateUser(ctx context.Context, name, email string) (*userv1.User, error)
}

type userService struct {
	logger *zap.Logger
}

func NewService(logger *zap.Logger) Service {
	return &userService{logger: logger}
}

func (s *userService) GetUser(ctx context.Context, id string) (*userv1.User, error) {
	s.logger.Info("fetching user", zap.String("id", id))
	return &userv1.User{Id: id, Name: "John Doe", Email: "john@example.com"}, nil
}

func (s *userService) CreateUser(ctx context.Context, name, email string) (*userv1.User, error) {
	s.logger.Info("creating user", zap.String("email", email))
	return &userv1.User{Id: "new-uuid", Name: name, Email: email}, nil
}

// --- Handler ---

type Handler struct {
	svc Service
}

func NewHandler(svc Service) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) RegisterRoutes(r chi.Router) {
	r.Get("/users/{id}", h.GetUser)
	r.Post("/users", h.CreateUser)
}

func (h *Handler) GetUser(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	user, err := h.svc.GetUser(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	b, _ := protojson.Marshal(user)
	w.Write(b)
}

func (h *Handler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var req userv1.CreateUserRequest

	// Unmarshal from JSON to Proto message
	dec := json.NewDecoder(r.Body)
	var raw map[string]interface{}
	if err := dec.Decode(&raw); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	jsonBytes, _ := json.Marshal(raw)
	if err := protojson.Unmarshal(jsonBytes, &req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	user, err := h.svc.CreateUser(r.Context(), req.GetName(), req.GetEmail())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	b, _ := protojson.Marshal(user)
	w.Write(b)
}
