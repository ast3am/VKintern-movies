package db

import (
	"context"
	"github.com/ast3am/VKintern-movies/internal/db/mocks"
	"github.com/ast3am/VKintern-movies/internal/models"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

const (
	newActorUUID = "b0482c7a-1a4c-4a3c-9463-35f0036a0d62"
	newMovieUUID = "f44d4a8f-7f16-4c1d-836b-02e0b8de4a00"
	userEmail    = "testuser@mail.com"
)

func GetDate(date string) time.Time {
	layout := "2006-01-02"
	t, _ := time.Parse(layout, date)
	return t
}

var testCfg = models.Config{
	ListenPort: "",
	SqlConfig: models.SqlConfig{
		UsernameDB: "postgres",
		PasswordDB: "password",
		HostDB:     "localhost",
		PortDB:     "5430",
		DBName:     "film_library_db",
		DelayTime:  10,
	},
	LogLevel: "",
}

func TestNewClient(t *testing.T) {
	ctx := context.Background()
	log := mocks.NewLogger(t)

	testTable := []struct {
		name string
		cfg  models.Config
	}{
		{
			name: "positive",
			cfg:  testCfg,
		},
		{
			name: "negative",
			cfg: models.Config{
				ListenPort: "",
				SqlConfig: models.SqlConfig{
					UsernameDB: "postgres",
					PasswordDB: "password",
					HostDB:     "localhost",
					PortDB:     "54320",
					DBName:     "organization_db",
					DelayTime:  1,
				},
				LogLevel: "",
			},
		},
	}
	for _, test := range testTable {
		if test.name == "positive" {
			log.On("DebugMsg", "fail connect to DB, try again").Return()
			log.On("DebugMsg", "connection to DB is OK").Return()
			db, _ := NewClient(ctx, &test.cfg, log)
			require.NotNil(t, db)
			defer db.dbConnect.Close(ctx)
		} else {
			log.On("DebugMsg", "fail connect to DB, try again").Return()
			db, _ := NewClient(ctx, &test.cfg, log)
			require.Nil(t, db)
		}
	}

}

func TestDB_CreateActor(t *testing.T) {
	ctx := context.Background()
	log := mocks.NewLogger(t)
	log.On("DebugMsg", "connection to DB is OK").Return()
	db, _ := NewClient(ctx, &testCfg, log)
	defer db.dbConnect.Close(ctx)

	testTable := []struct {
		name      string
		id        uuid.UUID
		testActor models.Actor
	}{
		{
			name: "positive",
			id:   uuid.MustParse(newActorUUID),
			testActor: models.Actor{
				Name:      "test_actor_1",
				Gender:    "male",
				BirthDate: GetDate("2002-01-01"),
			},
		}, {
			name: "negative",
			id:   uuid.MustParse(newActorUUID),
			testActor: models.Actor{
				Name:      "test_actor_2",
				Gender:    "male",
				BirthDate: GetDate("2002-01-01"),
			},
		},
	}
	for _, test := range testTable {
		if test.name == "positive" {
			err := db.CreateActor(ctx, test.id, test.testActor)
			assert.Nil(t, err)
		} else {
			err := db.CreateActor(ctx, test.id, test.testActor)
			assert.NotNil(t, err)
		}
	}
}

func TestDB_UpdateActor(t *testing.T) {
	ctx := context.Background()
	log := mocks.NewLogger(t)
	log.On("DebugMsg", "connection to DB is OK").Return()
	db, _ := NewClient(ctx, &testCfg, log)
	defer db.dbConnect.Close(ctx)
	id := uuid.MustParse(newActorUUID)
	testActor := models.Actor{
		Name:      "test_actor_1",
		Gender:    "female",
		BirthDate: GetDate("2002-01-01"),
	}
	err := db.UpdateActor(ctx, id, testActor)
	assert.Nil(t, err, "")

}

func TestDB_GetActorByUUID(t *testing.T) {
	ctx := context.Background()
	log := mocks.NewLogger(t)
	log.On("DebugMsg", "connection to DB is OK").Return()
	db, _ := NewClient(ctx, &testCfg, log)
	defer db.dbConnect.Close(ctx)
	id := uuid.MustParse(newActorUUID)
	expected := &models.Actor{
		Name:      "test_actor_1",
		Gender:    "female",
		BirthDate: GetDate("2002-01-01"),
	}
	result, err := db.GetActorByUUID(ctx, id)
	assert.Equal(t, expected, result)
	assert.Nil(t, err)
}

func TestDB_DeleteActor(t *testing.T) {
	ctx := context.Background()
	log := mocks.NewLogger(t)
	log.On("DebugMsg", "connection to DB is OK").Return()
	db, _ := NewClient(ctx, &testCfg, log)
	defer db.dbConnect.Close(ctx)
	id := uuid.MustParse(newActorUUID)
	err := db.DeleteActor(ctx, id)
	assert.Nil(t, err)
}

func TestDB_GetActorList(t *testing.T) {
	ctx := context.Background()
	log := mocks.NewLogger(t)
	log.On("DebugMsg", "connection to DB is OK").Return()
	db, _ := NewClient(ctx, &testCfg, log)
	defer db.dbConnect.Close(ctx)
	expected := map[string][]string{
		"Brad Pitt":         {"Inception", "The Devil Wears Prada"},
		"Cate Blanchett":    {"Inception", "Pretty Woman"},
		"Julia Roberts":     {"Forrest Gump", "The Devil Wears Prada"},
		"Leonardo DiCaprio": {"Training Day"},
		"Tom Hanks":         {"Forrest Gump", "Pretty Woman", "Training Day"},
	}
	result, err := db.GetActorList(ctx)
	assert.Equal(t, expected, result)
	assert.Nil(t, err)
}

func TestDB_GetUserByEmail(t *testing.T) {
	ctx := context.Background()
	log := mocks.NewLogger(t)
	log.On("DebugMsg", "connection to DB is OK").Return()
	db, _ := NewClient(ctx, &testCfg, log)
	defer db.dbConnect.Close(ctx)
	email := userEmail
	expected := &models.User{
		Email:    "testuser@mail.com",
		Password: "userPassword",
		Role:     "user",
	}
	result, err := db.GetUserByEmail(ctx, email)
	assert.Equal(t, expected, result)
	assert.Nil(t, err)
}

func TestDB_CreateMovie(t *testing.T) {
	ctx := context.Background()
	log := mocks.NewLogger(t)
	log.On("DebugMsg", "connection to DB is OK").Return()
	db, _ := NewClient(ctx, &testCfg, log)
	defer db.dbConnect.Close(ctx)

	testTable := []struct {
		name      string
		id        uuid.UUID
		testMovie models.Movie
	}{
		{
			name: "positive",
			id:   uuid.MustParse(newMovieUUID),
			testMovie: models.Movie{
				Name:        "test_movie_1",
				Description: "Description 1",
				ReleaseDate: GetDate("2002-01-01"),
				Rating:      6.2,
				ActorList:   []string{"test_actor_1", "test_actor_2"},
			},
		}, {
			name: "negative",
			id:   uuid.MustParse(newMovieUUID),
			testMovie: models.Movie{
				Name:        "test_movie_1",
				Description: "Description 1",
				ReleaseDate: GetDate("2002-01-01"),
				Rating:      6.2,
				ActorList:   []string{"test_actor_1", "test_actor_2"},
			},
		},
	}

	idActor := uuid.MustParse(newActorUUID)
	testActor := models.Actor{
		Name:      "test_actor_1",
		Gender:    "male",
		BirthDate: GetDate("2002-01-01"),
	}

	for _, test := range testTable {
		if test.name == "positive" {
			err := db.CreateActor(ctx, idActor, testActor)
			assert.Nil(t, err)
			err = db.CreateMovie(ctx, test.id, test.testMovie)
			assert.Nil(t, err)
			err = db.DeleteActor(ctx, idActor)
			assert.Nil(t, err)
		} else {
			err := db.CreateMovie(ctx, test.id, test.testMovie)
			assert.NotNil(t, err)
		}
	}
}

func TestDB_UpdateMovie(t *testing.T) {
	ctx := context.Background()
	log := mocks.NewLogger(t)
	log.On("DebugMsg", "connection to DB is OK").Return()
	db, _ := NewClient(ctx, &testCfg, log)
	defer db.dbConnect.Close(ctx)
	id := uuid.MustParse(newMovieUUID)
	idActor := uuid.MustParse(newActorUUID)
	testActor := models.Actor{
		Name:      "test_actor_1",
		Gender:    "male",
		BirthDate: GetDate("2002-01-01"),
	}
	testMovie := models.Movie{
		Name:        "test_movie_1",
		Description: "Description 2",
		ReleaseDate: GetDate("2001-01-01"),
		Rating:      9.2,
		ActorList:   []string{"test_actor_1", "test_actor_3"},
	}
	err := db.CreateActor(ctx, idActor, testActor)
	assert.Nil(t, err)
	err = db.UpdateMovie(ctx, id, testMovie)
	assert.Nil(t, err, "")

}

func TestDB_GetMovieByUUID(t *testing.T) {
	ctx := context.Background()
	log := mocks.NewLogger(t)
	log.On("DebugMsg", "connection to DB is OK").Return()
	db, _ := NewClient(ctx, &testCfg, log)
	defer db.dbConnect.Close(ctx)
	id := uuid.MustParse(newMovieUUID)
	expected := &models.Movie{
		Name:        "test_movie_1",
		Description: "Description 2",
		ReleaseDate: GetDate("2001-01-01"),
		Rating:      9.2,
		ActorList:   []string{"test_actor_1", "test_actor_3"},
	}
	result, err := db.GetMovieByUUID(ctx, id)
	assert.Equal(t, expected, result)
	assert.Nil(t, err)
	idActor := uuid.MustParse(newActorUUID)
	err = db.DeleteActor(ctx, idActor)
	assert.Nil(t, err)
}

func TestDB_DeleteMovie(t *testing.T) {
	ctx := context.Background()
	log := mocks.NewLogger(t)
	log.On("DebugMsg", "connection to DB is OK").Return()
	db, _ := NewClient(ctx, &testCfg, log)
	defer db.dbConnect.Close(ctx)
	id := uuid.MustParse(newMovieUUID)
	err := db.DeleteMovie(ctx, id)
	assert.Nil(t, err)
}

func TestDB_GetMovieList(t *testing.T) {
	ctx := context.Background()
	log := mocks.NewLogger(t)
	log.On("DebugMsg", "connection to DB is OK").Return()
	db, _ := NewClient(ctx, &testCfg, log)
	defer db.dbConnect.Close(ctx)
	testSort := "release_date"
	testLine := "asc"
	expected := []*models.Movie{
		{
			Name:        "Pretty Woman",
			Description: "Description 4",
			ReleaseDate: time.Date(1990, 3, 23, 0, 0, 0, 0, time.UTC),
			Rating:      7,
			ActorList:   []string{"Cate Blanchett", "Tom Hanks"},
		},
		{
			Name:        "Forrest Gump",
			Description: "Description 1",
			ReleaseDate: time.Date(1994, 7, 6, 0, 0, 0, 0, time.UTC),
			Rating:      8.8,
			ActorList:   []string{"Julia Roberts", "Tom Hanks"},
		},
		{
			Name:        "Training Day",
			Description: "Description 3",
			ReleaseDate: time.Date(2001, 10, 5, 0, 0, 0, 0, time.UTC),
			Rating:      7.7,
			ActorList:   []string{"Leonardo DiCaprio", "Tom Hanks"},
		},
		{
			Name:        "The Devil Wears Prada",
			Description: "Description 2",
			ReleaseDate: time.Date(2006, 6, 30, 0, 0, 0, 0, time.UTC),
			Rating:      6.9,
			ActorList:   []string{"Brad Pitt", "Julia Roberts"},
		},
		{
			Name:        "Inception",
			Description: "Description 5",
			ReleaseDate: time.Date(2010, 7, 16, 0, 0, 0, 0, time.UTC),
			Rating:      8.8,
			ActorList:   []string{"Brad Pitt", "Cate Blanchett"},
		},
	}

	result, err := db.GetMovieList(ctx, testSort, testLine)
	assert.Equal(t, expected, result)
	assert.Nil(t, err)
}

func TestDB_GetMovie(t *testing.T) {
	ctx := context.Background()
	log := mocks.NewLogger(t)
	log.On("DebugMsg", "connection to DB is OK").Return()
	db, _ := NewClient(ctx, &testCfg, log)
	defer db.dbConnect.Close(ctx)
	testActor := "%Jul%"
	testMovie := "%%"
	expected := []*models.Movie{
		{
			Name:        "Forrest Gump",
			Description: "Description 1",
			ReleaseDate: time.Date(1994, 7, 6, 0, 0, 0, 0, time.UTC),
			Rating:      8.8,
			ActorList:   []string{"Julia Roberts"},
		},
		{
			Name:        "The Devil Wears Prada",
			Description: "Description 2",
			ReleaseDate: time.Date(2006, 6, 30, 0, 0, 0, 0, time.UTC),
			Rating:      6.9,
			ActorList:   []string{"Julia Roberts"},
		},
	}
	result, err := db.GetMovie(ctx, testActor, testMovie)
	assert.Equal(t, expected, result)
	assert.Nil(t, err)
}
