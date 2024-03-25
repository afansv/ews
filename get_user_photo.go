package ews

import (
	"encoding/xml"
	"errors"
)

type GetUserPhotoRequest struct {
	XMLName       struct{} `xml:"m:GetUserPhoto"`
	Email         string   `xml:"m:Email"`
	SizeRequested string   `xml:"m:SizeRequested"`
}

type GetUserPhotoResponse struct {
	Response
	Etag        string `xml:"-"`
	HasChanged  bool   `xml:"HasChanged"`
	PictureData string `xml:"PictureData"`
}

type getUserPhotoResponseEnvelop struct {
	XMLName struct{}                 `xml:"Envelope"`
	Body    getUserPhotoResponseBody `xml:"Body"`
}
type getUserPhotoResponseBody struct {
	GetUserPhotoResponse GetUserPhotoResponse `xml:"GetUserPhotoResponse"`
}

// GetUserPhoto
// https://docs.microsoft.com/en-us/exchange/client-developer/web-service-reference/getuserphoto-operation
func GetUserPhoto(c Client, r *GetUserPhotoRequest) (*GetUserPhotoResponse, error) {

	xmlBytes, err := xml.MarshalIndent(r, "", "  ")
	if err != nil {
		return nil, err
	}

	bb, respHeaders, err := c.SendAndReceiveWithHeaders(xmlBytes)
	if err != nil {
		return nil, err
	}

	var soapResp getUserPhotoResponseEnvelop
	err = xml.Unmarshal(bb, &soapResp)
	if err != nil {
		return nil, err
	}

	if soapResp.Body.GetUserPhotoResponse.ResponseClass == ResponseClassError {
		return nil, errors.New(soapResp.Body.GetUserPhotoResponse.MessageText)
	}

	soapResp.Body.GetUserPhotoResponse.Etag = respHeaders.Get("etag")

	return &soapResp.Body.GetUserPhotoResponse, nil
}
