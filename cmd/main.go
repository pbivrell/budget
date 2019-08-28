package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/kjk/betterguid"
	"github.com/pbivrell/budget/influx"
	"github.com/pbivrell/budget/log"
)

const (
	STATIC_HTML = "./static/"
	DB_IP       = "http://localhost:8086"
	database    = "budget"
	measurement = "test"
)

type app struct {
}

func main() {
	r := mux.NewRouter()
	a := &app{}
	r.HandleFunc("/expenses", a.Expenses())
	r.PathPrefix("/").Handler(http.StripPrefix("/", http.FileServer(http.Dir(STATIC_HTML))))
	fmt.Println(http.ListenAndServe(":4231", r))
}

type ExpenseTransaction struct {
	Date            time.Time `json:"date"`
	PrimaryCategory string    `json:"primaryCategory"`
	SubCategory     string    `json:"subCategory"`
	Description     string    `json:"description"`
	Amount          float64   `json:"amount"`
}

func (a *app) Expenses() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var request []ExpenseTransaction
		err := json.NewDecoder(r.Body).Decode(&request)
		if err != nil {
			log.Logger.Errorf("Failed to parse response body: %+v", err)
			http.Error(w, "Invalid json request", http.StatusBadRequest)
			return
		}
		connection, err := influx.New(DB_IP)
		if err != nil {
			log.Logger.Errorf("Failed to connection to database: %+v", err)
			http.Error(w, "Failed to connect to database", http.StatusInternalServerError)
			return
		}
		connection.Series(database).Measurement(measurement).Precision("ns")
		for _, v := range request {
			connection.Point(
				map[string]string{
					"primaryCategory": v.PrimaryCategory,
					"subCategory":     v.SubCategory,
					"id":              betterguid.New(),
				},
				map[string]interface{}{
					"amount":      v.Amount,
					"description": v.Description,
				},
				v.Date,
			)
		}
		err = connection.BatchWrite()
		if err != nil {
			log.Logger.Errorf("Failed to write data to database: %+v", err)
			http.Error(w, "Failed to write to database", http.StatusInternalServerError)
		}

		fmt.Fprintf(w, "Successfully recorded %d transactions", len(request))
	}
}

// https://blog.gopheracademy.com/advent-2016/advanced-encoding-decoding/
// http://choly.ca/post/go-json-marshalling/
func (r *ExpenseTransaction) UnmarshalJSON(data []byte) error {
	type Alias ExpenseTransaction
	aux := &struct {
		Date string `json:"date"`
		*Alias
	}{
		Alias: (*Alias)(r),
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	var err error
	r.Date, err = time.Parse("1/2/06", aux.Date)
	if err != nil {
		return err
	}

	return nil
}
