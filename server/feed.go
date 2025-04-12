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

	"github.com/gorilla/feeds"
)

func (s Server) Feed(w http.ResponseWriter, r *http.Request) {
	result, err := s.back.GetVerificationResult(r.PathValue("site"))
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to get host verification: %s", err), http.StatusInternalServerError)
		return
	}
	if !result.Success {
		http.Error(w, "unverified host", http.StatusBadRequest)
		return
	}

	feed := feeds.Feed{
		Title:   fmt.Sprintf("Momiji - %s", r.PathValue("site")),
		Link:    &feeds.Link{Href: fmt.Sprintf("https://%s", r.PathValue("site"))},
		Created: result.Created,
	}

	posts, err := s.back.GetPosts(r.PathValue("site"))
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to get posts: %s", err), http.StatusInternalServerError)
		return
	}
	for _, post := range posts {
		feed.Items = append(feed.Items, &feeds.Item{
			Title:   post.Body,
			Link:    &feeds.Link{Href: fmt.Sprintf("https://momiji.chat/%s#post%d", r.PathValue("site"), post.ID())},
			Author:  &feeds.Author{Name: fmt.Sprintf("%s (%s)", post.Persona, post.Author)},
			Created: post.Created,
		})

		if feed.Updated.Before(post.Created) {
			feed.Updated = post.Created
		}
	}

	switch r.PathValue("type") {
	case "atom":
		w.Header().Set("Content-Type", "text/xml")
		feed.WriteAtom(w)
	case "json":
		w.Header().Set("Content-Type", "application/json")
		feed.WriteJSON(w)
	default: // rss
		w.Header().Set("Content-Type", "text/xml")
		feed.WriteRss(w)
	}
}
