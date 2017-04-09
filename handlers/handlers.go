// Package handlers provides endpoints for the web client to request song, post
// playlists to Spotify, and authenticate the server to make requests
// on its behalf.
package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"mime"
	"net/http"
	"strings"
	"time"

	"github.com/snyderks/selfbots.net/configRead"
)

// Page is a basic page, with Body being an HTML doc.
type Page struct {
	Title string
	Body  []byte
}

type token struct {
	Token      string   `json:"token"`
	Selections []string `json:"selections"`
}

// configLocation is the location of the config
// (should be in the same directory as the application)
const configLocation = "config.json"

// config is the translated structure of the application's config file.
var config configRead.Config

// assetsHandler is a catch-all for any static assets that the page needs,
// such as JS dependencies, images, CSS files, etc.
// Meant to be passed to AddHandler in an http server.
func assetsHandler(w http.ResponseWriter, r *http.Request) {
	loc := r.URL.Path[len("/assets/"):]
	f, err := ioutil.ReadFile("assets/" + loc)
	var contentType string
	if strings.HasSuffix(loc, ".css") {
		contentType = "text/css"
	} else if strings.HasSuffix(loc, ".png") {
		contentType = "image/png"
	} else if strings.HasSuffix(loc, ".js") {
		contentType = "application/javascript"
	} else if strings.HasSuffix(loc, ".svg") {
		contentType = "image/svg+xml"
	} else {
		contentType = "text/plain"
	}
	if err != nil {
		w.WriteHeader(404)
		return
	}
	w.Header().Add("Content-Type", contentType)
	fmt.Fprintf(w, "%s", f)
}

// indexHandler serves up the landing page for the site. Meant to be passed to AddHandler in
// an http server.
func indexHandler(w http.ResponseWriter, r *http.Request) {
	p, err := loadPage("index")
	if err != nil {
		fmt.Fprintf(w, "Error")
		return
	}
	fmt.Fprintf(w, "%s", p.Body)
}

// notFoundHandler is for serving a 404 page. Meant to be passed to AddHandler in an http
// server for a default.
func notFoundHandler(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Path[len("/404/"):]
	renderTemplate(w, "notfound", &Page{Title: title})
}

// discordLoginURLHandler is to initially redirect the user to the Spotify
// authentication page
func discordLoginURLHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Received request for Discord Login URL")
	type loginURL struct {
		URL string `json:"URL"`
	}
	var err error
	if err != nil {
		http.Error(w, "Failed to generate state. Something went wrong "+
			"or something is vulnerable.", http.StatusInternalServerError)
	}
	callback := "http://" + config.Hostname + config.HTTPPort + config.AuthRedirectHandler
	url := config.AuthURL + config.DiscordKey + "&redirect_url=" + callback +
		"&response_type=token" + "&scope=" + strings.Join(config.Scopes, "+")

	urlJSON, err := json.Marshal(loginURL{url})
	if err == nil {
		w.Write(urlJSON)
	} else {
		w.WriteHeader(500)
		w.Write([]byte(err.Error()))
	}
}

// authHandler passes a received token to the correct page to process it in
// the browser and store for future use.
func authHandler(w http.ResponseWriter, r *http.Request) {
	requestBody, err := ioutil.ReadAll(r.Body)
	if err == nil {
		tok := token{}
		err = json.Unmarshal(requestBody, &tok)
		if err == nil {
			var c = &http.Client{
				Timeout: time.Second * 10,
			}
			c.Post(config.BotServerName+config.BotServerPath,
				mime.TypeByExtension(".json"), bytes.NewBuffer(requestBody))
			w.WriteHeader(200)
			return
		}
		fmt.Println("Failed to get token: " + err.Error())
	}
	w.WriteHeader(500)
}

func authReceiver(w http.ResponseWriter, r *http.Request) {
	p, err := loadPage("auth")
	if err != nil {
		w.WriteHeader(404)
		return
	}
	fmt.Fprintf(w, "%s", p.Body)
}

// Basic load, render functions
func loadPage(title string) (*Page, error) {
	filename := title + ".html"
	body, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return &Page{Title: title, Body: body}, nil
}

func renderTemplate(w http.ResponseWriter, templ string, p *Page) {
	t, _ := template.ParseFiles(templ + ".html")
	t.Execute(w, p)
}

// SetUpAPICalls Create handler functions for api calls
func SetUpAPICalls() {
	http.HandleFunc("/api/discordLoginUrl/", discordLoginURLHandler)
	http.HandleFunc("/sendToken", authHandler)
	// http.HandleFunc("/api/getSpotifyUser", spotifyUserHandler)
	// http.HandleFunc("/api/getPlaylist", createLastFmPlaylist)
	// http.HandleFunc("/api/createPlaylist", postPlaylistToSpotify)
}

// SetUpBasicHandlers Create handler functions for path handlers
func SetUpBasicHandlers() {
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/assets/", assetsHandler)
	http.HandleFunc("/404/", notFoundHandler)
	http.HandleFunc("/auth", authReceiver)
}

// Initial setup.
func init() {
	fmt.Println("Handlers initializing")
	var err error
	config, err = configRead.ReadConfig(configLocation)
	config.Scopes = []string{"identify", "rpc.api", "rpc.notifications.read"}
	if err != nil {
		panic("Couldn't read the config. It's either not there or isn't in the correct format.")
	}
}
