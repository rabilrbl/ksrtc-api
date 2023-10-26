package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
)

var (
	dataCache     sync.Map
	cacheDuration = 24 * time.Hour // Adjust the cache duration as needed
	cacheMutex    sync.Mutex
	lastCacheTime time.Time
)

func indexHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte("API is up and running!"))
}

func allHandler(w http.ResponseWriter, r *http.Request) {
	from := r.URL.Query().Get("from")
	to := r.URL.Query().Get("to")

	data, err := fetchAllBusDataWithCache()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "An error occurred while fetching all bus data.")
		return
	}

	if from != "" && to != "" {
		var filteredData []PlacesData
		for _, bus := range data {
			if strings.Contains(strings.ToLower(bus.Value), strings.ToLower(from)) || strings.Contains(strings.ToLower(bus.Value), strings.ToLower(to)) {
				filteredData = append(filteredData, bus)
			}
		}
		respondWithJSON(w, filteredData)
	} else {
		respondWithJSON(w, data)
	}
}

func fetchAllBusDataWithCache() ([]PlacesData, error) {
	// Check if the cache has expired
	if time.Since(lastCacheTime) >= cacheDuration {
		// Lock to prevent multiple fetches
		cacheMutex.Lock()
		defer cacheMutex.Unlock()

		// Double-check in case another goroutine has updated the cache
		if time.Since(lastCacheTime) >= cacheDuration {
			newData, err := fetchAllBusData()
			if err != nil {
				return nil, err
			}

			// Store the fetched data in the cache and update the lastCacheTime
			dataCache.Store("allBusData", newData)
			lastCacheTime = time.Now()
		}
	}

	// Try to get data from the cache
	cachedData, found := dataCache.Load("allBusData")
	if found {
		return cachedData.([]PlacesData), nil
	}

	// If not found in the cache, this should not happen due to double-checking
	return nil, fmt.Errorf("data not found in cache")
}

func busHandler(w http.ResponseWriter, r *http.Request) {
	fromPlaceName := r.URL.Query().Get("fromPlaceName")
	startPlaceId := r.URL.Query().Get("startPlaceId")
	toPlaceName := r.URL.Query().Get("toPlaceName")
	endPlaceId := r.URL.Query().Get("endPlaceId")
	journeyDate := r.URL.Query().Get("journeyDate")

	if fromPlaceName == "" || startPlaceId == "" || toPlaceName == "" || endPlaceId == "" || journeyDate == "" {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "Missing query params")
		return
	}

	data, err := fetchBuses(fromPlaceName, startPlaceId, toPlaceName, endPlaceId, journeyDate)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "An error occurred while fetching buses: %v", err)
		return
	}

	respondWithJSON(w, data)
}

func main() {
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/all", allHandler)
	http.HandleFunc("/bus", busHandler)

	port := "8080" // Default port
	if envPort := os.Getenv("PORT"); envPort != "" {
		port = envPort
	}

	addr := ":" + port
	fmt.Printf("Server listening on port %s...\n", port)
	err := http.ListenAndServe(addr, nil)
	if err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
