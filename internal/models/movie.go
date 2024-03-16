package models

import (
	"encoding/json"
	"errors"
	"time"
)

type Movie struct {
	Name        string    `json:"name"`
	Description string    `json:"description"`
	ReleaseDate time.Time `json:"release_date"`
	Rating      float64   `json:"rating"`
	ActorList   []string  `json:"actor_list"`
}

func (m *Movie) Validate() error {
	if len(m.Name) > 150 || len(m.Name) < 1 || len(m.Description) > 1000 || m.Rating < 0 || m.Rating > 10 {
		return errors.New("no valid data")
	}
	return nil
}

func (m *Movie) UnmarshalJSON(data []byte) error {
	type Alias Movie
	aux := &struct {
		ReleaseDate string `json:"release_date"`
		*Alias
	}{
		Alias: (*Alias)(m),
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	var err error
	if aux.ReleaseDate == "" {
		m.ReleaseDate = time.Time{}
	} else {
		m.ReleaseDate, err = time.Parse("2006-01-02", aux.ReleaseDate)
		if err != nil {
			return err
		}
	}

	return nil
}
