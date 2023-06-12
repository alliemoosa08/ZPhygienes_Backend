package main

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/lib/pq"
)

type Service struct {
	Image                  string `json:"image"`
	Title                  string `json:"title"`
	Price                  string `json:"price"`
	IsSelected             bool   `json:"isSelected"`
	ServiceType            int    `json:"serviceType"`
	ServiceDescription     string `json:"serviceDescription"`
	DateCreated            string `json:"dateCreated"`
	DateModified           string `json:"dateModified"`
	RemovedFromBookingList bool   `json:"removedFromBookingList"`
}

func InsertService(w http.ResponseWriter, r *http.Request) {
	var newService Service

	db, err := ConnectDB()
	if err != nil {
		CheckError(err)
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode(err)
	}

	// Close the database connection when you're done
	defer db.Close()

	err = json.NewDecoder(r.Body).Decode(&newService)
	if err != nil {
		CheckError(err)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(err)
	}

	stmt, err := db.Prepare(`INSERT INTO services(image, title, price, is_selected, service_type, 
		service_description, date_created, date_modified, removed_from_booking_list) VALUES($1,$2,$3,$4,$5,$6,$7,$8,$9)`)

	if err != nil {
		CheckError(err)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(err)
	}
	defer stmt.Close()
	newService.IsSelected = false

	currentTime := time.Now()
	formattedTime := currentTime.Format("2006-01-02 15:04:05")

	newService.DateCreated = formattedTime
	newService.DateModified = formattedTime
	newService.RemovedFromBookingList = false

	_, err = stmt.Exec(&newService.Image, &newService.Title, &newService.Price, &newService.IsSelected, &newService.ServiceType,
		&newService.ServiceDescription, &newService.DateCreated, &newService.DateModified, &newService.RemovedFromBookingList)

	if err != nil {
		CheckError(err)
		if err.(*pq.Error).Code == "23505" {
			w.WriteHeader(http.StatusAlreadyReported)
			json.NewEncoder(w).Encode(err)
		} else {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(err)
		}

	} else {
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(newService)
	}

}

func GetServiceByServiceType(w http.ResponseWriter, r *http.Request) {
	var getService Service
	var getServices []Service
	serviceType := mux.Vars(r)["serviceType"]

	db, err := ConnectDB()
	if err != nil {
		CheckError(err)
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode(err)
	}

	// Close the database connection when you're done
	//defer db.Close()

	query := `SELECT image, title, price, is_selected, service_type, 
	service_description, date_created, date_modified, removed_from_booking_list FROM
	services  WHERE service_type = $1`

	// Execute the query
	rows, err := db.Query(query, serviceType)

	if err != nil {
		CheckError(err)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(err)
	}

	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&getService.Image, &getService.Title, &getService.Price, &getService.IsSelected, &getService.ServiceType,
			&getService.ServiceDescription, &getService.DateCreated, &getService.DateModified, &getService.RemovedFromBookingList)

		if err != nil {
			CheckError(err)
			if err == sql.ErrNoRows {
				w.WriteHeader(http.StatusNoContent)
				json.NewEncoder(w).Encode(err)
			}
		}
		getServices = append(getServices, getService)
	}

	if len(getServices) == 0 {
		w.WriteHeader(http.StatusNoContent)
		json.NewEncoder(w).Encode(err)
	} else {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(getServices)
	}

}
