package user

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	repository "github.com/user/go-templates/template-mysql/internal/user/sqlc"
	"go.uber.org/zap"
)

// --- Domain ---

type User struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type Repository interface {
	Get(ctx context.Context, id string) (*User, error)
	Create(ctx context.Context, user *User) error
}

type Service interface {
	GetUser(ctx context.Context, id string) (*User, error)
	CreateUser(ctx context.Context, user *User) error
}

// --- Service Implementation ---

type userService struct {
	repo   Repository
	logger *zap.Logger
}

func NewService(repo Repository, logger *zap.Logger) Service {
	return &userService{
		repo:   repo,
		logger: logger,
	}
}

func (s *userService) GetUser(ctx context.Context, id string) (*User, error) {
	s.logger.Info("fetching user", zap.String("id", id))
	return s.repo.Get(ctx, id)
}

func (s *userService) CreateUser(ctx context.Context, user *User) error {
	s.logger.Info("creating user", zap.String("email", user.Email))
	return s.repo.Create(ctx, user)
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
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(user)
}

func (h *Handler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var user User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	if err := h.svc.CreateUser(r.Context(), &user); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

// --- MySQL Repository ---

type MysqlRepository struct {
	q  *repository.Queries
	db *sql.DB
}

func NewMysqlRepository(db *sql.DB) *MysqlRepository {
	return &MysqlRepository{
		q:  repository.New(db),
		db: db,
	}
}

func (r *MysqlRepository) Get(ctx context.Context, id string) (*User, error) {
	userModel, err := r.q.GetUser(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	return &User{
		ID:    userModel.ID,
		Name:  userModel.Name,
		Email: userModel.Email,
	}, nil
}

func (r *MysqlRepository) Create(ctx context.Context, user *User) error {
	user.ID = uuid.New().String()

	params := repository.CreateUserParams{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
	}

	_, err := r.q.CreateUser(ctx, params)
	return err
}
