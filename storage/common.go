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

package storage

import (
	"bytes"
	"errors"
	"time"

	"github.com/patapancakes/momiji/identity"
)

var ErrNonExistentPost = errors.New("non-existent post")

type Post struct {
	Author  identity.ID `json:"author"`
	Persona string      `json:"persona"`
	Body    string      `json:"body"`
	Created time.Time   `json:"created"`
}

func (p Post) ID() int64 {
	return p.Created.UnixMilli()
}

func (p Post) IsCreatedBy(id identity.ID) bool {
	return bytes.Equal(p.Author, id)
}

type VerificationResult struct {
	Requester identity.ID `json:"requester"`
	Success   bool        `json:"success"`
	Created   time.Time   `json:"created"`
}

type Backend interface {
	// GetVerificationResult returns the Host's VerificationResult
	GetVerificationResult(host string) (VerificationResult, error)
	// AddVerificationResult adds a Host's VerificationResult
	AddVerificationResult(host string, result VerificationResult) error
	// GetLatestVerificationResultByID returns the latest VerificationResult by an ID
	GetLatestVerificationResultByID(id identity.ID) (VerificationResult, error)

	// GetPosts returns all Posts for a Host, sorted newest to oldest
	GetPosts(host string) ([]Post, error)
	// GetPost returns a Post for a Host
	GetPost(host string, id int64) (Post, error)
	// AddPost adds a new Post for a Host
	AddPost(host string, post Post) error
	// DeletePost deletes a post for a Host
	DeletePost(host string, id int64) error
	// GetLatestPostByID returns the latest Post for an ID
	GetLatestPostByID(host string, id identity.ID) (Post, error)
}
