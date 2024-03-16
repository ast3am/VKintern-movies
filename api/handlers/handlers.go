package handlers

import (
	"context"
	"github.com/ast3am/VKintern-movies/internal/models"
	"net/http"
)

const (
	MethodNotAllowed    = "invalid request method"
	UnprocessableEntity = "No valid data"
	ParsingJSONError    = "JSON parsing error"
	InvalidToken        = "Invalid token"
)

const (
	Admin = "admin"
	User  = "user"
)

type logger interface {
}

type services interface {
	Auth(ctx context.Context, email, password string) (token string, err error)
	CheckToken(token, permissionLevel string) error
	CreateActor(ctx context.Context, actor models.Actor) error
	GetActorList(ctx context.Context) (map[string][]string, error)
	DeleteActor(ctx context.Context, id string) error
	UpdateActor(ctx context.Context, id string, actor models.Actor) error
	CreateMovie(ctx context.Context, movie models.Movie) error
}

type Handler struct {
	services services
	log      logger
}

func NewHandler(serv services, log logger) *Handler {
	return &Handler{
		services: serv,
		log:      log,
	}
}

func (h *Handler) RegisterHandlers(mux *http.ServeMux) {
	mux.HandleFunc("/auth", h.Auth)                    //+
	mux.HandleFunc("/actor/create", h.CreateActor)     //+
	mux.HandleFunc("/actor/get-list", h.GetActorsList) //+
	mux.HandleFunc("/actor/update/", h.UpdateActor)    //+
	mux.HandleFunc("/actor/delete/", h.DeleteActor)    //+
	mux.HandleFunc("/movie/create", h.CreateMovie)     //+
	mux.HandleFunc("/movie/update/", h.UpdateMovie)
	mux.HandleFunc("/movie/delete/", h.DeleteMovie)
	mux.HandleFunc("/movie/get-list/", h.GetMoviesList)
	mux.HandleFunc("/movie/get-movie/", h.GetMovie)
	http.Handle("/", mux)
}
