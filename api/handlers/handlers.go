package handlers

import (
	"context"
	_ "github.com/ast3am/VKintern-movies/docs"
	"github.com/ast3am/VKintern-movies/internal/models"
	httpSwagger "github.com/swaggo/http-swagger/v2"

	"net/http"
)

const (
	MethodNotAllowed    = "invalid request method"
	UnprocessableEntity = "no valid data"
	ParsingJSONError    = "JSON parsing error"
	InvalidToken        = "invalid token"
)

const (
	Admin = "admin"
	User  = "user"
)

//go:generate mockery --name logger
type logger interface {
	HandlerErrorLog(r *http.Request, status int, msg string, err error)
	HandlerLog(r *http.Request, status int, msg string)
}

//go:generate mockery --name service
type service interface {
	Auth(ctx context.Context, email, password string) (token string, err error)
	CheckToken(token, permissionLevel string) error
	CreateActor(ctx context.Context, actor models.Actor) error
	GetActorList(ctx context.Context) (map[string][]string, error)
	DeleteActor(ctx context.Context, id string) error
	UpdateActor(ctx context.Context, id string, actor models.Actor) error
	CreateMovie(ctx context.Context, movie models.Movie) error
	UpdateMovie(ctx context.Context, id string, movie models.Movie) error
	DeleteMovie(ctx context.Context, id string) error
	GetMovieList(ctx context.Context, sortby, line string) ([]*models.Movie, error)
	GetMovie(ctx context.Context, actor, movie string) ([]*models.Movie, error)
}

type Handler struct {
	services service
	log      logger
}

func NewHandler(serv service, log logger) *Handler {
	return &Handler{
		services: serv,
		log:      log,
	}
}

func (h *Handler) RegisterHandlers(mux *http.ServeMux) {
	mux.HandleFunc("/auth", h.Auth)                    //+
	mux.HandleFunc("/actor/create", h.CreateActor)     //+
	mux.HandleFunc("/actor/get-list", h.GetActorsList) //+
	mux.HandleFunc("/actor/update/", h.UpdateActor)
	mux.HandleFunc("/actor/delete/", h.DeleteActor)
	mux.HandleFunc("/movie/create", h.CreateMovie) //+
	mux.HandleFunc("/movie/update/", h.UpdateMovie)
	mux.HandleFunc("/movie/delete/", h.DeleteMovie)
	mux.HandleFunc("/movie/get-list", h.GetMoviesList) //+
	mux.HandleFunc("/movie/get-movie", h.GetMovie)     //+
	mux.Handle("/swagger/", httpSwagger.Handler(
		httpSwagger.URL("http://localhost:8080/swagger/doc.json"),
	))
}
