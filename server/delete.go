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
	"strconv"
)

func (s Server) Delete(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Sec-Fetch-Site") != "" && r.Header.Get("Sec-Fetch-Site") != "same-origin" {
		http.Error(w, "request looks forged", http.StatusBadRequest)
		return
	}
	if r.Header.Get("Referer") != "" {
		u, err := url.Parse(r.Header.Get("Referer"))
		if err != nil {
			http.Error(w, fmt.Sprintf("failed to parse referer header: %s", err), http.StatusBadRequest)
			return
		}

		if u.Host != "momiji.chat" {
			http.Error(w, "request looks forged", http.StatusBadRequest)
			return
		}
	}

	result, err := s.back.GetVerificationResult(r.PathValue("site"))
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to get host verification: %s", err), http.StatusInternalServerError)
		return
	}
	if !result.Success {
		http.Error(w, "unverified host", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to decode id value: %s", err), http.StatusBadRequest)
		return
	}

	post, err := s.back.GetPost(r.PathValue("site"), int64(id))
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to get post: %s", err), http.StatusInternalServerError)
		return
	}
	if !post.IsCreatedBy(s.ident.Derive(r.PathValue("site"), GetRequestIP(r))) {
		http.Error(w, "post not owned by you", http.StatusUnauthorized)
		return
	}

	err = s.back.DeletePost(r.PathValue("site"), int64(id))
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to delete post: %s", err), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/%s", r.PathValue("site")), http.StatusSeeOther)
}
