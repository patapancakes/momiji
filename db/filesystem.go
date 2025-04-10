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

package db

import (
	"bytes"
	"encoding/json"
	"io"
	"os"
	"path"
	"slices"

	"github.com/patapancakes/momiji/identity"
)

type Filesystem struct {
	Path string
}

func NewFilesystemBackend(path string) (Filesystem, error) {
	_, err := os.Stat(path)
	if err != nil {
		if !os.IsNotExist(err) {
			return Filesystem{}, err
		}

		err = os.Mkdir(path, 0755)
		if err != nil {
			return Filesystem{}, err
		}
	}

	return Filesystem{Path: path}, nil
}

func (fs Filesystem) GetVerificationResults() (map[string]VerificationResult, error) {
	f, err := os.Open(path.Join(fs.Path, "verified.json"))
	if err != nil {
		if os.IsNotExist(err) {
			return make(map[string]VerificationResult), nil
		}

		return nil, err
	}

	defer f.Close()

	var data map[string]VerificationResult
	err = json.NewDecoder(f).Decode(&data)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (fs Filesystem) GetVerificationResult(host string) (VerificationResult, error) {
	results, err := fs.GetVerificationResults()
	if err != nil {
		return VerificationResult{}, err
	}

	result, _ := results[host]

	return result, nil
}

func (fs Filesystem) AddVerificationResult(host string, result VerificationResult) error {
	results, err := fs.GetVerificationResults()
	if err != nil {
		return err
	}

	results[host] = result

	f, err := os.OpenFile(path.Join(fs.Path, "verified.json"), os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		return err
	}

	defer f.Close()

	_, err = f.Seek(0, io.SeekStart)
	if err != nil {
		return err
	}

	err = json.NewEncoder(f).Encode(results)
	if err != nil {
		return err
	}

	return nil
}

func (fs Filesystem) GetLatestVerificationResultByID(id identity.ID) (VerificationResult, error) {
	results, err := fs.GetVerificationResults()
	if err != nil {
		return VerificationResult{}, err
	}

	for _, result := range results {
		if bytes.Equal(result.Requester, id) {
			return result, nil
		}
	}

	return VerificationResult{}, nil
}

func (fs Filesystem) GetPosts(host string) ([]Post, error) {
	f, err := os.Open(path.Join(fs.Path, host, "posts.json"))
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}

		return nil, err
	}

	defer f.Close()

	var data []Post
	err = json.NewDecoder(f).Decode(&data)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (fs Filesystem) GetPost(host string, id int64) (Post, error) {
	posts, err := fs.GetPosts(host)
	if err != nil {
		return Post{}, err
	}

	for _, post := range posts {
		if post.ID() == id {
			return post, nil
		}
	}

	return Post{}, ErrNonExistentPost
}

func (fs Filesystem) AddPost(host string, post Post) error {
	_, err := os.Stat(path.Join(fs.Path, host))
	if err != nil {
		if !os.IsNotExist(err) {
			return err
		}

		err = os.Mkdir(path.Join(fs.Path, host), 0755)
		if err != nil {
			return err
		}
	}

	posts, err := fs.GetPosts(host)
	if err != nil {
		return err
	}

	posts = append([]Post{post}, posts...)

	err = fs.WritePostsFile(host, posts)
	if err != nil {
		return err
	}

	return nil
}

func (fs Filesystem) DeletePost(host string, id int64) error {
	posts, err := fs.GetPosts(host)
	if err != nil {
		return err
	}

	var found bool
	for i, post := range posts {
		if post.ID() == id {
			found = true
			posts = slices.Delete(posts, i, i+1)
			break
		}
	}
	if !found {
		return ErrNonExistentPost
	}

	err = fs.WritePostsFile(host, posts)
	if err != nil {
		return err
	}

	return nil
}

func (fs Filesystem) GetLatestPostByID(host string, id identity.ID) (Post, error) {
	posts, err := fs.GetPosts(host)
	if err != nil {
		return Post{}, err
	}

	var latest Post
	for _, post := range posts {
		if !bytes.Equal(post.Author, id) {
			continue
		}

		if latest.Created.Before(post.Created) {
			latest = post
		}
	}

	return latest, nil
}

func (fs Filesystem) WritePostsFile(host string, posts []Post) error {
	f, err := os.OpenFile(path.Join(fs.Path, host, "posts.json"), os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		return err
	}

	defer f.Close()

	err = f.Truncate(0)
	if err != nil {
		return err
	}

	_, err = f.Seek(0, io.SeekStart)
	if err != nil {
		return err
	}

	err = json.NewEncoder(f).Encode(posts)
	if err != nil {
		return err
	}

	return nil
}
