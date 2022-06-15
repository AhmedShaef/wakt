// Package feed geolocation provides a simple function to get the geolocation of an IP address.
package feed

import (
	"encoding/json"
	"net/http"
)

// Geolocation is the IP geolocation service.
type Geolocation struct {
	IP          string  `json:"ip"`
	CountryCode string  `json:"country_code"`
	CountryName string  `json:"country_name"`
	RegionCode  string  `json:"region_code"`
	RegionName  string  `json:"region_name"`
	City        string  `json:"city"`
	ZipCode     string  `json:"zip_code"`
	TimeZone    string  `json:"time_zone"`
	Latitude    float64 `json:"latitude"`
	Longitude   float64 `json:"longitude"`
	MetroCode   int     `json:"metro_code"`
}

// GetGeolocation returns geolocation data using IP address.
func GetGeolocation(ip string) (*Geolocation, error) {
	res, err := http.Get("https://api.freegeoip.app/json/" + ip + "?apikey=" + "42e5cd90-bfcc-11ec-91be-c30d312f695e")
	if err != nil {
		return nil, err
	}
	geo := &Geolocation{}
	err = json.NewDecoder(res.Body).Decode(geo)
	if err != nil {
		return nil, err
	}
	return geo, nil
}
