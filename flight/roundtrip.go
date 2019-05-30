package flight

// RoundTrip is a round trip ticket
type RoundTrip struct {
	GoFlight     *Flight `json:"go"`
	ReturnFlight *Flight `json:"return"`
}
