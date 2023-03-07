// Copyright 2023 NJWS Inc.

package device

import (
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

type RedfishDevice struct {
	Description Description `json:"description,omitempty"`
	Api         string      `json:"api,omitempty"`
	Login       string      `json:"login,omitempty"`
	Password    string      `json:"password,omitempty"`
}

// func (s RedfishDevice) UUID() string {
// 	return strings.TrimPrefix(s.Description.Device.UDN, "uuid:")
// }

func GetDescription(url string) (d Description, err error) {
	resp, err := http.Get(url)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return
	}
	return d, xml.Unmarshal(data, &d)
}

type Description struct {
	XMLName     xml.Name `xml:"root" json:"root,omitempty"`
	Text        string   `xml:",chardata" json:"-"`
	Xmlns       string   `xml:"xmlns,attr" json:"xmlns,omitempty"`
	SpecVersion struct {
		Text  string `xml:",chardata" json:"-"`
		Major string `xml:"major"`
		Minor string `xml:"minor"`
	} `xml:"specVersion" json:"specversion,omitempty"`
	Device struct {
		Text             string `xml:",chardata" json:"-"`
		DeviceType       string `xml:"deviceType"`
		FriendlyName     string `xml:"friendlyName"`
		Manufacturer     string `xml:"manufacturer"`
		ManufacturerURL  string `xml:"manufacturerURL"`
		ModelName        string `xml:"modelName"`
		ModelNumber      string `xml:"modelNumber"`
		ModelDescription string `xml:"modelDescription"`
		ModelURL         string `xml:"modelURL"`
		UDN              string `xml:"UDN"`
		ServiceList      struct {
			Text    string `xml:",chardata" json:"-"`
			Service struct {
				Text        string `xml:",chardata" json:"-"`
				ServiceType string `xml:"serviceType"`
				ServiceId   string `xml:"serviceId"`
				ControlURL  string `xml:"controlURL"`
				EventSubURL string `xml:"eventSubURL"`
				SCPDURL     string `xml:"SCPDURL"`
			} `xml:"service" json:"service,omitempty"`
		} `xml:"serviceList" json:"servicelist,omitempty"`
		PresentationURL string `xml:"presentationURL"`
	} `xml:"device" json:"device,omitempty"`
}

func (d Description) ToDevice(u *url.URL) RedfishDevice {
	return RedfishDevice{
		Description: d,
		Api:         fmt.Sprintf("%s://%s", u.Scheme, u.Hostname()),
	}
}
