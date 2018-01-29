package gotwilio

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

// ValidationResponse is returned after a text/sms message is posted to Twilio
type LookupResponse struct {
    countryCode    string `json:"country_code"`
    phoneNumber    string `json:"phone_number"`
    nationalFormat string `json:"national_format"`
    url            string `json:"url"`
}

// Validate phone number uses Twilio LookUp service.
// See https://www.twilio.com/lookup for more information.
func (twilio *Twilio) ValidatePhoneNumber(phone string) (lookupResponse *LookupResponse, exception *Exception, err error) {

	lookupUrl := "https://lookups.twilio.com/v1"
	twilioUrl := lookupUrl + "/PhoneNumbers/" + phone

    if twilio.TestMode {
        return &LookupResponse{countryCode: "UA", phoneNumber: phone, nationalFormat: phone, url: twilioUrl}, nil, nil
    }

	res, err := twilio.get(twilioUrl)
	if err != nil {
		return lookupResponse, exception, err
	}
	defer res.Body.Close()

	responseBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return lookupResponse, exception, err
	}

	if res.StatusCode != http.StatusOK {
		exception = new(Exception)
		err = json.Unmarshal(responseBody, exception)

		// We aren't checking the error because we don't actually care.
		// It's going to be passed to the client either way.
		return lookupResponse, exception, err
	}

	lookupResponse = new(LookupResponse)
	err = json.Unmarshal(responseBody, lookupResponse)
	return lookupResponse, exception, err
}