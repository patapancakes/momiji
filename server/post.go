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
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/patapancakes/momiji/identity"
	"github.com/patapancakes/momiji/storage"
)

func Post(w http.ResponseWriter, r *http.Request) {
	ip := GetRequestIP(r)

	result, err := Backend.GetVerificationResult(r.PathValue("site"))
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to get host verification: %s", err), http.StatusInternalServerError)
		return
	}
	if !result.Success {
		if result.Created.Add(time.Minute * 10).After(time.Now().UTC()) {
			http.Error(w, "rate limited", http.StatusTooManyRequests)
			return
		}

		result.Requester = identity.Derive("verification", ip)

		latest, err := Backend.GetLatestVerificationResultByID(result.Requester)
		if err != nil {
			http.Error(w, fmt.Sprintf("failed to get latest verification result: %s", err), http.StatusInternalServerError)
			return
		}
		if latest.Created.Add(time.Minute * 10).After(time.Now().UTC()) {
			http.Error(w, "rate limited", http.StatusTooManyRequests)
			return
		}

		u, err := url.Parse(r.PostFormValue("referer"))
		if err != nil {
			http.Error(w, fmt.Sprintf("failed to parse referer value: %s", err), http.StatusBadRequest)
			return
		}
		if u.Host != r.PathValue("site") {
			http.Error(w, "referer does not match", http.StatusBadRequest)
			return
		}

		result.Success, err = VerifyHost(u)
		if err != nil {
			http.Error(w, fmt.Sprintf("failed to check host verification: %s", err), http.StatusInternalServerError)
			return
		}

		result.Created = time.Now().UTC()

		err = Backend.AddVerificationResult(u.Host, result)
		if err != nil {
			http.Error(w, fmt.Sprintf("failed to add verified host: %s", err), http.StatusInternalServerError)
			return
		}

		if !result.Success {
			http.Error(w, "unverified host", http.StatusBadRequest)
			return
		}
	}

	author := identity.Derive(r.PathValue("site"), ip)

	latest, err := Backend.GetLatestPostByID(r.PathValue("site"), author)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to get latest user post: %s", err), http.StatusInternalServerError)
		return
	}
	if latest.Created.Add(time.Second * 10).After(time.Now().UTC()) {
		http.Error(w, "rate limited", http.StatusTooManyRequests)
		return
	}

	name := strings.TrimSpace(r.PostFormValue("name"))
	if len(name) > 12 {
		http.Error(w, "invalid name", http.StatusBadRequest)
		return
	}

	body := strings.TrimSpace(r.PostFormValue("body"))
	if body == "" || len(body) > 140 {
		http.Error(w, "invalid body", http.StatusBadRequest)
		return
	}

	err = Backend.AddPost(r.PathValue("site"), storage.Post{Author: author, Persona: name, Body: body, Created: time.Now().UTC()})
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to insert post: %s", err), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, r.PathValue("site"), http.StatusSeeOther)
}
