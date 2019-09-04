package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

// ErrStatusCode an error thrown when a given http call responds with a bad http status
type ErrStatusCode struct {
	Code                   int
	Message                string
	SuggestedRetryDuration time.Duration
}

func (err *ErrStatusCode) Error() string {
	return fmt.Sprintf(
		"Http call failed with a status code: %d. Reason: %s",
		err.Code,
		err.Message,
	)
}

func checkResponse(res *http.Response) error {
	status := res.StatusCode
	if status >= 200 && status < 300 {
		return nil
	}

	statusErr := &ErrStatusCode{Code: status}
	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}
	if len(data) > 0 {
		statusErr.Message = string(data)
	}
	if statusErr.Code == 429 {
		rawRetrySecs := res.Header.Get("Retry-After")
		if rawRetrySecs != "" {
			retrySecs, _ := strconv.ParseInt(rawRetrySecs, 10, 64)
			statusErr.SuggestedRetryDuration = time.Duration(retrySecs) * time.Second
		}
	}

	return statusErr
}
