package main

import (
	"encoding/json"
	"net/http"
	"time"
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
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(err)
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newService)
}
