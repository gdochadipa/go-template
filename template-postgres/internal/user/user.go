package user

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	repository "github.com/user/go-templates/template-postgres/internal/user/sqlc"
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

// --- Postgres Repository ---

type PostgresRepository struct {
	q  *repository.Queries
	db *pgxpool.Pool
}

func NewPostgresRepository(db *pgxpool.Pool) *PostgresRepository {
	return &PostgresRepository{
		q:  repository.New(db),
		db: db,
	}
}

func (r *PostgresRepository) Get(ctx context.Context, id string) (*User, error) {
	var uuid pgtype.UUID
	if err := uuid.Scan(id); err != nil {
		return nil, fmt.Errorf("invalid uuid: %w", err)
	}

	userModel, err := r.q.GetUser(ctx, uuid)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, err
	}

	return &User{
		ID:    fmt.Sprintf("%x-%x-%x-%x-%x", userModel.ID.Bytes[0:4], userModel.ID.Bytes[4:6], userModel.ID.Bytes[6:8], userModel.ID.Bytes[8:10], userModel.ID.Bytes[10:16]),
		Name:  userModel.Name,
		Email: userModel.Email,
	}, nil
}

func (r *PostgresRepository) Create(ctx context.Context, user *User) error {
	params := repository.CreateUserParams{
		Name:  user.Name,
		Email: user.Email,
	}

	userModel, err := r.q.CreateUser(ctx, params)
	if err != nil {
		return err
	}

	user.ID = fmt.Sprintf("%x-%x-%x-%x-%x", userModel.ID.Bytes[0:4], userModel.ID.Bytes[4:6], userModel.ID.Bytes[6:8], userModel.ID.Bytes[8:10], userModel.ID.Bytes[10:16])
	return nil
}
