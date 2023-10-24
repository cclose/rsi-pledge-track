package model

import (
	"errors"
	"fmt"
	"github.com/labstack/echo/v4"
	"time"
)

// RequestFormat is a custom string type for specifying request formats.
type RequestFormat string

// ValidRequestFormats is a list of valid request formats.
var ValidRequestFormats = []RequestFormat{
	"json",
	"application/json",
	"html",
	"text/html",
	"csv",
	"text/csv",
}

// RequestFormatNegotiation is a list of formats for Accepts Negotiation
// This is passed to the Negotiate function and returns the value based on
// what the client wants to receive
var RequestFormatNegotiation = map[string]string{
	"application/json": "json",
	"text/csv":         "csv",
	"text/html":        "html",
	//	"text/plain":       "plain",
}

// PledgeDataRequest is a struct representing a request with a format and timestamp.
type PledgeDataRequest struct {
	Format         RequestFormat `query:"format"`   // Format specifies the request format.
	TimeStamp      time.Time     `query:"dateTime"` // TimeStamp specifies the request timestamp.
	AfterTimestamp time.Time     `query:"startingDateTime"`
	Offset         int           `query:"offset"`
	Limit          int           `query:"limit"`
}

// ParseRequest parses the request and sets the format and timestamp properties of PledgeDataRequest.
func (pdr *PledgeDataRequest) ParseRequest(c echo.Context) error {
	format := ""
	qpBinder := echo.QueryParamsBinder(c)
	qpBinder.Time("timestamp", &pdr.TimeStamp, time.RFC3339)
	qpBinder.Time("startingDateTime", &pdr.AfterTimestamp, time.RFC3339)
	qpBinder.String("format", &format)
	qpBinder.Int("offset", &pdr.Offset)
	qpBinder.Int("limit", &pdr.Limit)
	bindErrs := qpBinder.BindErrors()

	errMsg := ""
	if len(bindErrs) > 0 {
		for _, bErr := range bindErrs {
			bindErr, ok := bErr.(*echo.BindingError)
			if ok {
				errMsg += fmt.Sprintf("\n%s=\"%s\": %s", bindErr.Field, bindErr.Values, bindErr.Internal.Error())
			} else {
				errMsg += fmt.Sprintf("\nUnknownField:%s", bErr.Error())
			}
		}
	}

	// if format was passed, call the Setter to validate it
	if format != "" {
		if err := pdr.SetFormat(format); err != nil {
			errMsg = fmt.Sprintf("%s\nformat=\"[%s]\": %s", errMsg, format, err.Error())
		}
	}

	if errMsg != "" {
		return errors.New("Error parsing request: " + errMsg)
	} //implicit else

	return nil
}

// SetFormat validates the format against the list of valid formats and sets the format property if valid.
// It returns an error if the format is not valid.
func (pdr *PledgeDataRequest) SetFormat(format string) error {
	valid := false
	for _, validFormat := range ValidRequestFormats {
		if RequestFormat(format) == validFormat {
			if formatValue, ok := RequestFormatNegotiation[format]; ok {
				format = formatValue
			}
			valid = true
			break
		}
	}

	if !valid {
		return fmt.Errorf("Invalid format: %s", format)
	}

	pdr.Format = RequestFormat(format)
	return nil
}
