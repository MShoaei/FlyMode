package flight

import (
	"time"
)

// Flight is a complete flight struct with at most 2 stops.
// Fields without omitempty tag are the minimum requirement for a valid flight
type Flight struct {
	Source         string        `json:"source"`
	Destination    string        `json:"destination"`
	TakeoffDate    string        `json:"takeoff_date"`
	TakeoffTime    string        `json:"takeoff_time"`
	LandingDate    string        `json:"landing_date"`
	LandingTime    string        `json:"landing_time"`
	Duration       time.Duration `json:"duration"`
	RemainingSeats int           `json:"remaining_seats,omitempty"`
	Website        string        `json:"website"`
	Price          float64       `json:"price"`
	Rating         float32       `json:"rating,omitempty"`
	// one stop & two stop
	FirstAirline      string `json:"first_airline"`
	FirstFlightNumber string `json:"first_flight_number"`
	FirstPlaneModel   string `json:"first_plane_model,omitempty"`
	FirstFlightClass  string `json:"first_flight_class"`
	FirstCatering     string `json:"first_catering,omitempty"`
	FirstCargoWeight  int    `json:"first_cargo_weight,omitempty"`
	FirstHandCargo    int    `json:"first_hand_cargo,omitempty"`

	FirstMiddleCountry  string    `json:"first_middle_country,omitempty"`
	FirstChangeTerminal bool      `json:"first_change_terminal,omitempty"`
	FirstChangeAirport  bool      `json:"first_change_airport,omitempty"`
	FirstStopDuration   time.Time `json:"first_stop_duration,omitempty"`
	// two stop
	SecondAirline      string `json:"second_airline"`
	SecondFlightNumber string `json:"second_flight_number"`
	SecondPlaneModel   string `json:"second_plane_model,omitempty"`
	SecondFlightClass  string `json:"second_flight_class"`
	SecondCatering     string `json:"second_catering,omitempty"`
	SecondCargoWeight  int    `json:"second_cargo_weight,omitempty"`
	SecondHandCargo    int    `json:"second_hand_cargo,omitempty"`

	SecondMiddleCountry  string    `json:"second_middle_country,omitempty"`
	SecondChangeTerminal bool      `json:"second_change_terminal,omitempty"`
	SecondChangeAirport  bool      `json:"second_change_airport,omitempty"`
	SecondStopDuration   time.Time `json:"second_stop_duration,omitempty"`

	ThirdAirline      string `json:"third_airline"`
	ThirdFlightNumber string `json:"third_flight_number"`
	ThirdPlaneModel   string `json:"third_plane_model,omitempty"`
	ThirdFlightClass  string `json:"third_flight_class"`
	ThirdCatering     string `json:"third_catering,omitempty"`
	ThirdCargoWeight  int    `json:"third_cargo_weight,omitempty"`
	ThirdHandCargo    int    `json:"third_hand_cargo,omitempty"`
}
