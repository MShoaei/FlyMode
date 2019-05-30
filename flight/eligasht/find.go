package eligasht

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"

	"github.com/MShoaei/FlyMode/flight"
	"github.com/gocolly/colly"
	"github.com/google/uuid"
	// "github.com/MShoaei/FlyMode/uuid"
)

const (
	queryURL = "https://www.eligasht.com/Flight/Search"
)

type query struct {
	IsSearchRequest bool   `json:"is_search"`
	Category        string `json:"category"`
	Trip            string `json:"trip"`
	Source          string `json:"source"`
	SourceText      string `json:"source_text"`
	Destination     string `json:"destination"`
	DestinationText string `json:"destination_text"`
	Child           int    `json:"child"`
	Adult           int    `json:"adult"`
	Infant          int    `json:"infant"`
	FlightClass     string `json:"flight_class"`
	FromDate        string `json:"from_date"`
	ToDate          string `json:"to_date"`
	CheckOut        string `json:"check_out"`
	CheckIn         string `json:"check_in"`
	Rooms           string `json:"rooms"`
	NationalityCode string `json:"nationality_code"`
	ResidenceCode   string `json:"residence_code"`
	PreferredClass  string `json:"preferred_class"`
	PrefferedAir    string `json:"preferred_air"`
	Departure       string `json:"departure"`
	Birthday        string `json:"birthday"`
	PreSearch       bool   `json:"pre_search"`
	SearchKey       string `json:"search_key"`
	FlightGUID      string `json:"flight_guid"`
	CurrentCulture  string `json:"current_culture"`
	ChangeDay       int    `json:"change_day"`
	Type            string `json:"type"`
	ExclusiveTrain  bool   `json:"exclusive_train"`
	Key             string `json:"key"`
}

type flightDetail struct {
	Data         string
	ErrorCode    int
	ErrorMessage string
	Success      bool
}

func (q query) createPayload() []byte {
	return []byte(fmt.Sprintf("model%%5BIsSearchRequest%%5D=%t&model%%5BCategory%%5D=%s&model%%5BTrip%%5D=%s&model%%5BSource%%5D=%s&model%%5BSourceText%%5D=%s&model%%5BDestination%%5D=%s&model%%5BDestinationText%%5D=%s&model%%5BChild%%5D=%d&model%%5BAdult%%5D=%d&model%%5BInfant%%5D=%d&model%%5BflightClass%%5D=%s&model%%5BFromDate%%5D=%s&model%%5BToDate%%5D=%s&model%%5BCheckOut%%5D=%s&model%%5BCheckIn%%5D=%s&model%%5BRooms%%5D=%s&model%%5BNationalityCode%%5D=%s&model%%5BResidenceCode%%5D=%s&model%%5BPreferredClass%%5D=%s&model%%5BPreferedAir%%5D=%s&model%%5BpackRowId%%5D=&model%%5BpackXferIds%%5D=&model%%5BflightIds%%5D=&model%%5BDeparture%%5D=%s&model%%5BDeparturePackage%%5D=&model%%5BExtraCheckDate%%5D=&model%%5BTravelDuration%%5D=&model%%5BBirthDay%%5D=%s&model%%5BHotelofferId%%5D=&model%%5BFlightOfferId%%5D=&model%%5BpreSearch%%5D=%t&model%%5BhotelId%%5D=&model%%5BpreSearchUniqueId%%5D=&model%%5BdomesticCategory%%5D=&model%%5BsearchKey%%5D=%s&model%%5BflightGuid%%5D=%s&model%%5BCurrentCulture%%5D=%s&model%%5BChangeDay%%5D=%d&model%%5BPType%%5D=%s&model%%5BExclusiveTrain%%5D=%t&model%%5BKey%%5D=%s",
		q.IsSearchRequest, q.Category, q.Trip, q.Source,
		q.SourceText, q.Destination, q.DestinationText,
		q.Child, q.Adult, q.Infant, q.FlightClass,
		q.FromDate, q.ToDate, q.CheckOut, q.CheckIn,
		q.Rooms, q.NationalityCode, q.ResidenceCode,
		q.PreferredClass, q.PrefferedAir, q.Departure,
		q.Birthday, q.PreSearch, q.SearchKey, q.FlightGUID,
		q.CurrentCulture, q.ChangeDay, q.Type, q.ExclusiveTrain, q.Key))
}

func fill() {
	// json.Unmarshal()
}

// FindFlights finds all flights related to the query
func FindFlights(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		fmt.Fprintf(w, "%s is not supported", r.Method)
		return
	}

	var items []interface{}

	q := &query{}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprint(w, err)
		return
	}
	json.Unmarshal(body, q)
	// u, _ := uuid.NewV4()
	u, _ := uuid.NewRandom()
	q.Key = u.String()

	log.Println("Starting colly")
	c := colly.NewCollector(colly.UserAgent("Mozilla/5.0 (X11; Ubuntu; Linux x86_64; rv:66.0) Gecko/20100101 Firefox/66.0"))

	c.OnRequest(func(r *colly.Request) {
		r.Headers.Set("ost", "www.eligasht.com")
		r.Headers.Set("Referer", fmt.Sprintf("https://www.eligasht.com/Search?Category=%s&Trip=%s&Adult=%d&Child=%d&Infant=%d&FlightClass=%s&Key=%s",
			q.Category, q.Trip, q.Adult, q.Child, q.Infant, q.FlightClass, q.Key))
	})

	c.OnHTML("#SearchKey", func(e *colly.HTMLElement) {
		q.SearchKey = e.Attr("value")
	})

	c.OnHTML("#SearchResultFlight", func(e *colly.HTMLElement) {
		// fmt.Println("Hooray")
		d := c.Clone()
		e.ForEach("li", func(_ int, el *colly.HTMLElement) {
			stops := el.Attr("data-stop")
			if stops == "0" {
				if len(el.DOM.Filter("div.F-Info-departure").Children().Nodes) == 0 {
					f := &flight.OneWay{GoFlight: &flight.Flight{}}
					code := el.ChildText(".r-code:last-child")
					f.GoFlight.Source = code[1 : len(code)-1]
					f.GoFlight.Destination = el.Attr("data-airport")
					f.GoFlight.TakeoffTime = el.Attr("data-time")
					f.GoFlight.LandingTime = el.ChildText(`span[data-original-title="زمان ورود"]`)
					f.GoFlight.Website = "https://eligasht.com"
					p, _ := strconv.ParseFloat(el.Attr("data-cost"), 64)
					f.GoFlight.Price = p
					f.GoFlight.FirstAirline = el.Attr("data-airline")
					f.GoFlight.FirstFlightClass = el.Attr("data-class")

					d.OnResponse(func(cr *colly.Response) {
						details := &flightDetail{}
						json.Unmarshal(cr.Body, details)
						dom, err := goquery.NewDocumentFromReader(strings.NewReader(details.Data))
						if err != nil {
							log.Println(err)
							return
						}
						flightNum := dom.Find("div.segment-details-grid span:nth-child(3)").Text()
						f.GoFlight.FirstFlightNumber = strings.TrimSpace(strings.Split(flightNum, ":")[1])
					})

					d.Post("https://www.eligasht.com/Flight/GetFlightDetails",
						map[string]string{
							"FlightGUID": e.ChildAttr("input", "value"),
							"SearchKey":  q.SearchKey,
						})
					items = append(items, f)
				} else {
					// f := &flight.RoundTrip{}

				}
			} else if stops == "1" {
			} else if stops == "2" {
			}

		})
	})
	c.SetRequestTimeout(time.Second * 60)

	if err := c.PostRaw(queryURL, q.createPayload()); err != nil {
		log.Fatal(err)
		return
	}
	json.NewEncoder(w).Encode(items)
}
