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

package identity

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"net"
	"os"
)

type ID []byte

func (id ID) String() string {
	return base64.StdEncoding.EncodeToString(id)
}

type Identity struct {
	key []byte
}

func New(keyfile string) (Identity, error) {
	var err error

	var i Identity
	i.key, err = os.ReadFile(keyfile)
	if err == nil {
		return Identity{}, nil
	}
	if !os.IsNotExist(err) {
		return Identity{}, err
	}

	buf := make([]byte, 256)
	_, err = rand.Read(buf)
	if err != nil {
		return Identity{}, err
	}

	err = os.WriteFile(keyfile, buf, 0644)
	if err != nil {
		return Identity{}, err
	}

	return i, nil
}

func (i Identity) Derive(realm string, ip net.IP) ID {
	hash := sha256.New()

	hash.Write(i.key)
	hash.Write([]byte(realm))
	hash.Write(ip)

	return hash.Sum(nil)
}
