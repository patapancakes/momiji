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

package server

import (
	"fmt"
	"net"
	"net/http"
	"net/url"
	"strings"

	"github.com/patapancakes/momiji/identity"
	"github.com/patapancakes/momiji/storage"
	"golang.org/x/net/html"
)

type Server struct {
	mux   *http.ServeMux
	back  storage.Backend
	ident identity.Identity
}

func New(backend storage.Backend, identity identity.Identity) *Server {
	mux := http.NewServeMux()

	s := Server{mux: mux, back: backend, ident: identity}

	// http routes
	mux.HandleFunc("GET /", s.View)
	mux.HandleFunc("GET /{site}", s.View)
	mux.HandleFunc("POST /{site}", s.Post)

	mux.HandleFunc("GET /{site}/delete/{id}", s.Delete)

	mux.HandleFunc("GET /{site}/feed", s.Feed)
	mux.HandleFunc("GET /{site}/feed/{type}", s.Feed)

	// static files
	mux.Handle("GET /data/static/", http.StripPrefix("/data/static/", http.FileServer(http.Dir("static/"))))

	return &s
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.mux.ServeHTTP(w, r)
}

func GetRequestIP(r *http.Request) net.IP {
	host, _, _ := net.SplitHostPort(r.RemoteAddr) // assume this can't error
	if host == "127.0.0.1" && r.Header.Get("X-Forwarded-For") != "" {
		host = r.Header.Get("X-Forwarded-For")
	}

	return net.ParseIP(host)
}

func VerifyHost(u *url.URL) (bool, error) {
	resp, err := http.Get(fmt.Sprintf("https://%s/.well-known/momiji", u.Host))
	if err != nil {
		return false, err
	}
	if resp.StatusCode == 200 {
		return true, nil
	}

	resp, err = http.Get(u.String())
	if err != nil {
		return false, err
	}

	defer resp.Body.Close()

	t := html.NewTokenizer(resp.Body)
	for {
		tokenType := t.Next()
		if tokenType == html.ErrorToken {
			break
		}
		if tokenType != html.StartTagToken {
			continue
		}

		name, hasAttr := t.TagName()
		if string(name) != "iframe" {
			continue
		}
		if !hasAttr {
			continue
		}

		for {
			k, v, more := t.TagAttr()
			if string(k) == "src" && strings.Contains(string(v), "momiji.chat") {
				return true, nil
			}
			if !more {
				break
			}
		}
	}

	return false, nil
}
