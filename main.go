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

	"github.com/patapancakes/momiji/identity"
	"github.com/patapancakes/momiji/server"
	"github.com/patapancakes/momiji/storage"
)

func main() {
	// flags
	port := flag.Int("port", 80, "web server port")
	flag.Parse()

	back, err := storage.NewFilesystemBackend("data")
	if err != nil {
		log.Fatalf("failed to initialize storage backend: %s", err)
	}

	ident, err := identity.New("identity.key")
	if err != nil {
		log.Fatalf("failed to initialize identity package: %s", err)
	}

	s := server.New(back, ident)

	// start http server
	err = http.ListenAndServe(fmt.Sprintf(":%d", *port), s)
	if err != nil {
		log.Fatal(err)
	}
}
