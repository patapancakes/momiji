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
	"encoding/base64"
	"fmt"
	"html/template"
	"net/http"
	"net/url"
	"os"

	"github.com/patapancakes/momiji/identity"
	"github.com/patapancakes/momiji/storage"

	"github.com/xeonx/timeago"
)

type ViewData struct {
	Requester identity.ID
	Referer   *url.URL
	Posts     []storage.Post
}

var viewT = template.Must(template.New("main.html").Funcs(template.FuncMap{"timeago": timeago.English.Format, "b64": base64.StdEncoding.EncodeToString}).ParseGlob("templates/*.html"))

func View(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Referer") == "" {
		http.Redirect(w, r, "https://github.com/patapancakes/momiji", http.StatusSeeOther)
		return
	}

	var vd ViewData
	var err error

	vd.Referer, err = url.Parse(r.Header.Get("Referer"))
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to parse referer header: %s", err), http.StatusBadRequest)
		return
	}
	if vd.Referer.Host == "" {
		http.Error(w, "failed to derive host", http.StatusBadRequest)
		return
	}

	if r.PathValue("site") != "" && vd.Referer.Host == "momiji.chat" {
		vd.Referer, err = url.Parse(fmt.Sprintf("https://%s", r.PathValue("site")))
		if err != nil {
			http.Error(w, fmt.Sprintf("failed to parse site: %s", err), http.StatusBadRequest)
			return
		}
	}

	vd.Requester = identity.Derive(vd.Referer.Host, GetRequestIP(r))

	vd.Posts, err = Backend.GetPosts(vd.Referer.Host)
	if err != nil && !os.IsNotExist(err) {
		http.Error(w, fmt.Sprintf("failed to get posts: %s", err), http.StatusInternalServerError)
		return
	}

	err = viewT.Execute(w, vd)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to execute page template: %s", err), http.StatusInternalServerError)
		return
	}
}
