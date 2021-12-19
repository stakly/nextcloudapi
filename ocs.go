package nextcloudapi

import "encoding/xml"

// OCS XML server answer structs
type OCS struct {
	XMLName xml.Name `xml:"ocs"`
	Meta    Meta     `xml:"meta"`
	Data    Data     `xml:"data"`
}
type Meta struct {
	Status       string `xml:"status"`
	StatusCode   int    `xml:"statuscode"`
	Message      string `xml:"message"`
	TotalItems   string `xml:"totalitems"`
	ItemsPerPage string `xml:"itemsperpage"`
}
type Data struct {
	Users           Users       `xml:"users"`
	StorageLocation string      `xml:"storageLocation"`
	LastLogin       int64       `xml:"lastLogin"`
	Backend         string      `xml:"backend"`
	Id              string      `xml:"id"`
	Enabled         bool        `xml:"enabled"`
	Quota           Quota       `xml:"quota"`
	Email           string      `xml:"email"`
	Displayname     string      `xml:"displayname"`
	Phone           string      `xml:"phone"`
	Address         string      `xml:"address"`
	Website         string      `xml:"website"`
	Twitter         string      `xml:"twitter"`
	Datafield       []Datafield `xml:"element"`
	Groups          Groups      `xml:"groups"`
}
type Datafield struct {
	Field string `xml:",chardata"`
}
type Users struct {
	Username []Username `xml:"element"`
}
type Groups struct {
	Groupname []Groupname `xml:"element"`
}
type Groupname struct {
	Name string `xml:",chardata"`
}
type Username struct {
	Name string `xml:",chardata"`
}
type Quota struct {
	Free     int64   `xml:"free"`
	Used     int64   `xml:"used"`
	Total    int64   `xml:"total"`
	Relative float64 `xml:"relative"`
	Quota    int64   `xml:"quota"`
}
