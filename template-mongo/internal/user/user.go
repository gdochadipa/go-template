package user

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
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

// --- Mongo Repository ---

type MongoRepository struct {
	collection *mongo.Collection
}

func NewMongoRepository(db *mongo.Database) *MongoRepository {
	return &MongoRepository{
		collection: db.Collection("users"),
	}
}

type userDoc struct {
	ID    string `bson:"_id"`
	Name  string `bson:"name"`
	Email string `bson:"email"`
}

func (r *MongoRepository) Get(ctx context.Context, id string) (*User, error) {
	var doc userDoc
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&doc)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	return &User{
		ID:    doc.ID,
		Name:  doc.Name,
		Email: doc.Email,
	}, nil
}

func (r *MongoRepository) Create(ctx context.Context, user *User) error {
	user.ID = uuid.New().String()
	doc := userDoc{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
	}

	_, err := r.collection.InsertOne(ctx, doc)
	return err
}
