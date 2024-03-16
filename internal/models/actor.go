package models

import (
	"encoding/json"
	"time"
)

type Actor struct {
	Name      string    `json:"name"`
	Gender    string    `json:"gender"`
	BirthDate time.Time `json:"birth_date"`
}

func (a *Actor) UnmarshalJSON(data []byte) error {
	type Alias Actor
	aux := &struct {
		BirthDate string `json:"birth_date"`
		*Alias
	}{
		Alias: (*Alias)(a),
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	var err error
	if aux.BirthDate == "" {
		a.BirthDate = time.Time{}
	} else {
		a.BirthDate, err = time.Parse("2006-01-02", aux.BirthDate)
		if err != nil {
			return err
		}
	}

	return nil
}
