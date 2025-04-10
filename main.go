/*
	momiji - a free and simple embedded message box
	Copyright (C) 2025  Pancakes <patapancakes@pagefault.games>

	This program is free software: you can redistribute it and/or modify
	it under the terms of the GNU Affero General Public License as published by
	the Free Software Foundation, either version 3 of the License, or
	(at your option) any later version.

	This program is distributed in the hope that it will be useful,
	but WITHOUT ANY WARRANTY; without even the implied warranty of
	MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
	GNU Affero General Public License for more details.

	You should have received a copy of the GNU Affero General Public License
	along with this program.  If not, see <http://www.gnu.org/licenses/>.
*/

package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/patapancakes/momiji/db"
	"github.com/patapancakes/momiji/identity"
	"github.com/patapancakes/momiji/server"
)

func main() {
	// flags
	port := flag.Int("port", 80, "web server port")
	flag.Parse()

	// init identity package
	err := identity.Init("identity.key")
	if err != nil {
		log.Fatalf("failed to initialize identity package: %s", err)
	}

	// init storage backend
	server.Backend, err = db.NewFilesystemBackend("data")
	if err != nil {
		log.Fatalf("failed to initialize storage backend: %s", err)
	}

	// http routes
	http.HandleFunc("GET /", server.View)
	http.HandleFunc("GET /{site}", server.View)
	http.HandleFunc("POST /", server.Post)

	http.HandleFunc("GET /{site}/delete/{id}", server.Delete)

	http.HandleFunc("GET /{site}/feed", server.Feed)
	http.HandleFunc("GET /{site}/feed/{type}", server.Feed)

	// static files
	http.Handle("GET /data/static/", http.StripPrefix("/data/static/", http.FileServer(http.Dir("static/"))))

	// start http server
	err = http.ListenAndServe(fmt.Sprintf(":%d", *port), nil)
	if err != nil {
		log.Fatal(err)
	}
}
