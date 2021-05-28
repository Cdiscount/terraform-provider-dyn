package api

import (
	"encoding/json"
	"fmt"
	"strconv"
)

/*
This struct represents the request body that would be sent to the DynECT API
for logging in and getting a session token for future requests.
*/
type LoginBlock struct {
	Username     string `json:"user_name"`
	Password     string `json:"password"`
	CustomerName string `json:"customer_name"`
}

// Type ResponseBlock holds the "header" information returned by any call to
// the DynECT API.
//
// All response-type structs should include this as an anonymous/embedded field.
type ResponseBlock struct {
	Status   string         `json:"status"`
	JobId    int            `json:"job_id,omitempty"`
	Messages []MessageBlock `json:"msgs,omitempty"`
}

type YNBool bool
type SBool bool
type SInt int

type PublishBlock struct {
	Publish YNBool `json:"publish"`
}

// Type MessageBlock holds the message information from the server, and is
// nested within the ResponseBlock type.
type MessageBlock struct {
	Info      string `json:"INFO"`
	Source    string `json:"SOURCE"`
	ErrorCode string `json:"ERR_CD"`
	Level     string `json:"LVL"`
}

// Type LoginResponse holds the data returned by an HTTP POST call to
// https://api.dynect.net/REST/Session/.
type LoginResponse struct {
	ResponseBlock
	Data LoginDataBlock `json:"data"`
}

// Type LoginDataBlock holds the token and API version information from an HTTP
// POST call to https://api.dynect.net/REST/Session/.
//
// It is nested within the LoginResponse struct.
type LoginDataBlock struct {
	Token   string `json:"token"`
	Version string `json:"version"`
}

// RecordRequest holds the request body for a record create/update
type RecordRequest struct {
	RData DataBlock `json:"rdata"`
	TTL   string    `json:"ttl,omitempty"`
}

// PublishZoneBlock holds the request body for a publish zone request
// https://help.dyn.com/update-zone-api/
type PublishZoneBlock struct {
	Publish bool `json:"publish"`
}

func (val *SInt) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}
	if s == "" {
		*val = 0
		return nil
	}
	parsed_int, err := strconv.ParseInt(s, 10, 0)
	if err != nil {
		return err
	}
	*val = SInt(parsed_int)
	return nil
}

func (i SInt) MarshalJSON() ([]byte, error) {
	var val string
	if i == 0 {
		val = ""
	}
	val = fmt.Sprintf("%d", i)
	return json.Marshal(val)
}

func (val *YNBool) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}
	switch s {
	case "Y":
		*val = true
	case "N":
		*val = false
	default:
		return fmt.Errorf("Unknown boolean value for '%s'", s)
	}
	return nil
}

func (p YNBool) MarshalJSON() ([]byte, error) {
	val := "N"
	if p {
		val = "Y"
	}
	return json.Marshal(val)
}

func (val *SBool) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}
	switch s {
	case "true":
		*val = true
	case "false":
		*val = false
	default:
		return fmt.Errorf("Unknown boolean value for '%s'", s)
	}
	return nil
}

func (p SBool) MarshalJSON() ([]byte, error) {
	if p {
		return json.Marshal("true")
	} else {
		return json.Marshal("false")
	}
}
