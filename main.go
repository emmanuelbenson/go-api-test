package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

type SlactDetails struct {
	SlackName     string `json:"slack_name"`
	CurrentDay    string `json:"current_day"`
	UTCTime       string `json:"utc_time"`
	Track         string `json:"track"`
	GithubFileUrl string `json:"github_file_url"`
	GithubRepoUrl string `json:"github_repo_url"`
	StatusCode    int `json:"status_code"`
}

type detailsParam struct {
	SlackName string  `json:"slack_name"`
	Track     string `json:"track"`
}

type Response struct {
	Status  string `json:"status"`
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data"`
}

func main() {

	//create router
	router := mux.NewRouter()
	router.HandleFunc("/", getHome()).Methods("GET")
	router.HandleFunc("/api", getSlackDetails()).Methods("GET")

	//start server
	log.Fatal(http.ListenAndServe(":8000", jsonContentTypeMiddleware(router)))
}

func jsonContentTypeMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}

func SendJSONResponse(w http.ResponseWriter, status string, code int, message string, data any) {
	response := Response{Status: status, Code: code, Message: message, Data: data}
	json.NewEncoder(w).Encode(response)
}

func getHome() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		responseData := [0]int{}
		SendJSONResponse(w, "success", 200, "This API will get your slack details", responseData)
	}
}

func getSlackDetails() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		slackName := r.URL.Query().Get("slack_name")
		track := r.URL.Query().Get("track")

		if slackName == "" || track == "" {
			responseData := detailsParam{SlackName: slackName, Track: track}
			SendJSONResponse(w, "error", 422, "slack_name or track required", responseData)
			return
		}

		t := time.Now()

		var sd SlactDetails
		sd.SlackName = slackName
		sd.CurrentDay = t.Weekday().String()
		sd.UTCTime = t.UTC().String()
		sd.Track = track
		sd.GithubFileUrl = ""
		sd.GithubRepoUrl = ""
		sd.StatusCode = 200
		json.NewEncoder(w).Encode(sd)
	}
}
