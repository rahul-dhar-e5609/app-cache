package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gomodule/redigo/redis"
)

// HMSET album:1 title "Electric Ladyland" artist "Jimi Hendrix" price 4.95 likes 8
// HMSET album:2 title "Back in Black" artist "AC/DC" price 5.95 likes 3
// HMSET album:3 title "Rumours" artist "Fleetwood Mac" price 7.95 likes 12
// HMSET album:4 title "Nevermind" artist "Nirvana" price 5.95 likes 8
// ZADD likes 8 1 3 2 12 3 8 4

func main() {
	// Initialize a connection pool and assign it to the pool global
	// variable.
	pool = &redis.Pool{
		MaxIdle:     10,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", "localhost:6379")
		},
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/album", showAlbum)
	log.Println("Listening on :4000...")
	http.ListenAndServe(":4000", mux)
}

func showAlbum(w http.ResponseWriter, r *http.Request) {
	// Unless the request is using the GET method, return a 405 'Method
	// Not Allowed' response.
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodGet)
		http.Error(w, http.StatusText(405), 405)
		return
	}

	// Retrieve the id from the request URL query string. If there is
	// no id key in the query string then Get() will return an empty
	// string. We check for this, returning a 400 Bad Request response
	// if it's missing.
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, http.StatusText(400), 400)
		return
	}
	// Validate that the id is a valid integer by trying to convert it,
	// returning a 400 Bad Request response if the conversion fails.
	if _, err := strconv.Atoi(id); err != nil {
		http.Error(w, http.StatusText(400), 400)
		return
	}

	// Call the FindAlbum() function passing in the user-provided id.
	// If there's no matching album found, return a 404 Not Found
	// response. In the event of any other errors, return a 500
	// Internal Server Error response.
	bk, err := FindAlbum(id)
	if err == ErrNoAlbum {
		http.NotFound(w, r)
		return
	} else if err != nil {
		http.Error(w, http.StatusText(500), 500)
		return
	}

	// Write the album details as plain text to the client.
	fmt.Fprintf(w, "%s by %s: £%.2f [%d likes] \n", bk.Title, bk.Artist, bk.Price, bk.Likes)
}
