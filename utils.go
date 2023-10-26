package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly"
)

type BusData struct {
	Time    string `json:"time"`
	Seats   string `json:"seats"`
	Content string `json:"content"`
}

func fetchBuses(fromPlaceName, startPlaceId, toPlaceName, endPlaceId, journeyDate string) ([]BusData, error) {
	var busesAvailable []BusData

	c := colly.NewCollector()
	c.OnRequest(func(r *colly.Request) {
		r.Headers.Set("User-Agent", "Mozilla/5.0 (Linux; Android 10; SM-G975F) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/87.0.4280.141 Mobile Safari/537.36")
	})

	c.OnHTML(".rSetForward", func(e *colly.HTMLElement) {
		startTm := strings.Join(strings.Fields(e.DOM.Find(".StrtTm").Text()), " ")
		pTags := e.DOM.Find("p").Map(func(i int, s *goquery.Selection) string {
			return strings.TrimSpace(s.Text())
		})
		availCs := strings.TrimSpace(e.DOM.Find(".availCs").Text())

		busesAvailable = append(busesAvailable, BusData{
			Time:    startTm,
			Seats:   availCs,
			Content: strings.Join(pTags, " "),
		})
	})

	startPlaceId = url.QueryEscape(startPlaceId)
	endPlaceId = url.QueryEscape(endPlaceId)
	fromPlaceName = url.QueryEscape(fromPlaceName)
	toPlaceName = url.QueryEscape(toPlaceName)
	journeyDate = url.QueryEscape(journeyDate)

	webURL := fmt.Sprintf("https://m.ksrtc.in/oprs-mobile/forward/booking/avail/services.do?startPlaceId=%s&endPlaceId=%s&fromPlaceName=%s&toPlaceName=%s&txtJourneyDate=%s", startPlaceId, endPlaceId, fromPlaceName, toPlaceName, journeyDate)

	err := c.Visit(webURL)
	if err != nil {
		return nil, err
	}

	return busesAvailable, nil
}

type PlacesData struct {
	Id    string `json:"id"`
	Value string `json:"value"`
}

func fetchAllBusData() ([]PlacesData, error) {
	var data []PlacesData

	c := colly.NewCollector()
	c.OnRequest(func(r *colly.Request) {
		r.Headers.Set("User-Agent", "Mozilla/5.0 (Linux; Android 10; SM-G975F) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/87.0.4280.141 Mobile Safari/537.36")
	})

	c.OnHTML("#booking > div > script:nth-child(1)", func(e *colly.HTMLElement) {
		scriptText := e.Text
		if scriptText != "" {
			// Extract the JSON data from the script tag
			jsonText := scriptText[len("var jsondata = "):]
			err := json.Unmarshal([]byte(jsonText), &data)
			if err != nil {
				log.Printf("Error parsing JSON: %v", err)
			}
		}
	})

	webURL := "https://m.ksrtc.in/oprs-mobile/?OS=null"

	err := c.Visit(webURL)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func respondWithJSON(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	encodedJSON, err := json.Marshal(data)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "An error occurred while encoding JSON.")
		return
	}
	w.Write(encodedJSON)
}
