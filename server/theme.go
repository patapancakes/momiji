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
	"net/url"
	"regexp"
)

type Theme struct {
	Style   string
	Even    string
	Odd     string
	Border  string
	Text    string
	Message string
	Link    string
}

var IsValidHexColor = regexp.MustCompile("^[a-f0-9]{6}$").MatchString

func ThemeFromURLValues(v url.Values) Theme {
	var t Theme

	if IsValidHexColor(v.Get("style")) {
		t.Style = v.Get("style")
	}
	if IsValidHexColor(v.Get("bgeven")) {
		t.Even = v.Get("bgeven")
	}
	if IsValidHexColor(v.Get("bgodd")) {
		t.Odd = v.Get("bgodd")
	}
	if IsValidHexColor(v.Get("border")) {
		t.Border = v.Get("border")
	}
	if IsValidHexColor(v.Get("txt")) {
		t.Text = v.Get("txt")
	}
	if IsValidHexColor(v.Get("msg")) {
		t.Message = v.Get("msg")
	}
	if IsValidHexColor(v.Get("link")) {
		t.Link = v.Get("link")
	}

	return t
}

func (t Theme) Encode() string {
	v := make(url.Values)

	if t.Style != "" {
		v.Set("style", t.Style)
	}
	if t.Even != "" {
		v.Set("bgeven", t.Even)
	}
	if t.Odd != "" {
		v.Set("bgodd", t.Odd)
	}
	if t.Border != "" {
		v.Set("border", t.Border)
	}
	if t.Text != "" {
		v.Set("txt", t.Text)
	}
	if t.Message != "" {
		v.Set("msg", t.Message)
	}
	if t.Link != "" {
		v.Set("link", t.Link)
	}

	return v.Encode()
}
