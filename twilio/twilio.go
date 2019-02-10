package twilio

import (
	"fmt"
	"net/http"
)

// SMSReceived is all the data that twilio sends via SMS
type SMSReceived struct {
	SMSStatus     string
	FromCity      string
	ToCountry     string
	SmsMessageSid string
	ToCity        string
	SmsSid        string
	AccountSid    string
	FromZip       string
	Body          string
	To            string
	MessageSid    string
	From          string
	APIVersion    string
	NumMedia      string
	FromCountry   string
	ToZip         string
	NumSegments   string
	ToState       string
	FromState     string
}

// SMSParseForm takes in the http request for an SMS endpoint that your twilio number POSTs to
func SMSParseForm(r *http.Request) SMSReceived {
	r.ParseForm()
	return SMSReceived{
		SMSStatus:     r.FormValue("SmsStatus"),
		FromCity:      r.FormValue("FromCity"),
		ToCountry:     r.FormValue("ToCountry"),
		SmsMessageSid: r.FormValue("SmsMessageSid"),
		ToCity:        r.FormValue("ToCity"),
		SmsSid:        r.FormValue("SmsSid"),
		AccountSid:    r.FormValue("AccountSid"),
		FromZip:       r.FormValue("FromZip"),
		Body:          r.FormValue("Body"),
		To:            r.FormValue("To"),
		MessageSid:    r.FormValue("MessageSid"),
		From:          r.FormValue("From"),
		APIVersion:    r.FormValue("ApiVersion"),
		NumMedia:      r.FormValue("NumMedia"),
		FromCountry:   r.FormValue("FromCountry"),
		ToZip:         r.FormValue("ToZip"),
		NumSegments:   r.FormValue("NumSegments"),
		ToState:       r.FormValue("ToState"),
		FromState:     r.FormValue("FromState"),
	}
}

// SimpleTwiML wraps a string in a basic response
func SimpleTwiML(s string) string {
	return fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8"?>
	<Response>
	<Message>
	  %s
	</Message>
	</Response>`, s)
}
