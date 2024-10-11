package handlers

import (
	dao "GOLANG/project/dal"
	middleware "GOLANG/project/middlewares"
	"GOLANG/project/models"
	"GOLANG/project/utils"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

var geoDatInstance *dao.GeoDataDao

func GeoRoutes(router *mux.Router) {
	geoDatInstance, _ = dao.InitializeGeoDB()
	router.Handle("/upload", middleware.AuthMiddleware(http.HandlerFunc(uploadHandler))).Methods("POST")
	router.Handle("/update/{id}", middleware.AuthMiddleware(http.HandlerFunc(updateHandler))).Methods("PATCH") // Update/{GeoId}
	router.Handle("/list", middleware.AuthMiddleware(http.HandlerFunc(listHandler))).Methods("GET")            // List/{Email}
}

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	email := r.Context().Value("email").(string)
	var geoData models.GeoData
	ctx := context.Background()
	err := json.NewDecoder(r.Body).Decode(&geoData)
	if err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	geoData.Email = email

	err = geoDatInstance.Upload(ctx, &geoData)
	if err != nil {
		http.Error(w, "Error uploading", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "Uploaded GeoData", geoData.Name)
}

type geoDataOp struct {
	models.GeoData
	Id string
}

func listHandler(w http.ResponseWriter, r *http.Request) {
	email := r.Context().Value("email").(string)
	ctx := context.Background()

	rows, err := geoDatInstance.List(ctx, email)
	if err != nil {
		http.Error(w, "Error listing", http.StatusInternalServerError)
		return
	}

	var geoDataResults []*geoDataOp

	for rows.Next() {
		var longLatStr string
		var geoData geoDataOp
		var email string
		if err := rows.Scan(&geoData.Id, &geoData.Name, &longLatStr, &email); err != nil {
			fmt.Println(err)
			http.Error(w, "Error scanning row", http.StatusInternalServerError)
			return
		}
		fmt.Println(longLatStr)
		geoData.LongLatData, _ = utils.ConvertPostgresArrayTo2D(longLatStr)
		geoDataResults = append(geoDataResults, &geoData)
	}

	if err := rows.Err(); err != nil {
		http.Error(w, "Error after scanning rows", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(geoDataResults); err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusOK)
}

func updateHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	email := r.Context().Value("email").(string)
	geoId := vars["id"]
	ctx := context.Background()
	var geoData models.GeoData
	err := json.NewDecoder(r.Body).Decode(&geoData)
	if err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	ownerEmail, err := geoDatInstance.GetEmailFromId(ctx, geoId)
	if err != nil {
		http.Error(w, "Error Patching", http.StatusInternalServerError)
		return
	}

	if *ownerEmail != email {
		http.Error(w, "Not authorized to update", http.StatusUnprocessableEntity)
		return
	}

	err = geoDatInstance.Patch(ctx, geoId, &geoData)
	if err != nil {
		http.Error(w, "Error Patching", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(geoData)
}
