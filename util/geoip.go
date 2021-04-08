// Copyright (c) 2021 Miczone Asia.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package util

import (
	"net"

	"github.com/adrg/postcode"
	geoip2 "github.com/oschwald/geoip2-golang"
)

func ValidPostalCode(postalCode string) (bool, error) {
	if err := postcode.Validate(postalCode); err != nil {
		return false, err
	}
	return true, nil
}

func ParseCountryFromIP(db *geoip2.Reader, ipv4Str string) (*geoip2.City, error) {
	ip := net.ParseIP(ipv4Str)
	record, err := db.City(ip)
	if err != nil {
		return record, err
	}
	return record, nil
}
