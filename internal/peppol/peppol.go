package peppol

import (
	"encoding/xml"
	"time"
)

var xmlns = "http://www.peppol.eu/schema/pd/businesscard-generic/201907/"

type Root struct {
	XMLName      xml.Name       `xml:"root" json:"root"`
	XMLNS        string         `xml:"xmlns,attr" json:"xmlns"`
	Version      string         `xml:"version,attr" json:"version"`
	CreationDt   time.Time      `xml:"creationdt,attr" json:"creationdt"`
	BusinessCard []BusinessCard `xml:"businesscard" json:"businesscard"`
}

type BusinessCard struct {
	Participant Participant `xml:"participant" json:"participant"`
	Entity      Entity      `xml:"entity" json:"entity"`
	DoctypeIDs  []DoctypeID `xml:"doctypeid" json:"doctypeid"`
}

type Participant struct {
	Scheme string `xml:"scheme,attr" json:"scheme"`
	Value  string `xml:"value,attr" json:"value"`
}

type Entity struct {
	CountryCode string  `xml:"countrycode,attr" json:"countrycode"`
	Name        Name    `xml:"name" json:"name"`
	Contact     Contact `xml:"contact" json:"contact"`
}

type Name struct {
	Name string `xml:"name,attr" json:"name"`
}

type Contact struct {
	Type        string `xml:"type,attr" json:"type"`
	Name        string `xml:"name,attr" json:"name"`
	PhoneNumber string `xml:"phonenumber,attr" json:"phonenumber"`
	Email       string `xml:"email,attr" json:"email"`
}

type DoctypeID struct {
	Scheme      string `xml:"scheme,attr" json:"scheme"`
	Value       string `xml:"value,attr" json:"value"`
	DisplayName string `xml:"displayname,attr" json:"displayname"`
	Deprecated  string `xml:"deprecated,attr" json:"deprecated"`
}
