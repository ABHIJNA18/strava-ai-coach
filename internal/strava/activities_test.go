// This file verifies that Strava activity requests contain the correct filters.

package strava

import (
	"io"
	"net/http"
	"strconv"
	"strings"
	"testing"
	"time"
)

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(
	req *http.Request,
) (*http.Response, error) {
	return f(req)
}

func TestGetActivitiesPageIncludesAfterAndPagination(
	t *testing.T,
) {
	originalClient := stravaHTTPClient
	defer func() {
		stravaHTTPClient = originalClient
	}()

	var capturedRequest *http.Request

	stravaHTTPClient = &http.Client{
		Transport: roundTripFunc(
			func(req *http.Request) (*http.Response, error) {
				capturedRequest = req

				return &http.Response{
					StatusCode: http.StatusOK,
					Body: io.NopCloser(
						strings.NewReader("[]"),
					),
					Header: make(http.Header),
				}, nil
			},
		),
	}

	after := time.Unix(1700000000, 0)

	_, err := GetActivitiesPage(
		"test-access-token",
		2,
		100,
		after,
	)
	if err != nil {
		t.Fatal(err)
	}

	if capturedRequest == nil {
		t.Fatal("expected an outgoing request")
	}

	query := capturedRequest.URL.Query()

	expectedAfter := strconv.FormatInt(
		after.Unix(),
		10,
	)

	if query.Get("after") != expectedAfter {
		t.Fatalf(
			"expected after=%s, got %s",
			expectedAfter,
			query.Get("after"),
		)
	}

	if query.Get("page") != "2" {
		t.Fatalf(
			"expected page=2, got %s",
			query.Get("page"),
		)
	}

	if query.Get("per_page") != "100" {
		t.Fatalf(
			"expected per_page=100, got %s",
			query.Get("per_page"),
		)
	}

	if capturedRequest.Header.Get("Authorization") !=
		"Bearer test-access-token" {
		t.Fatal("expected bearer authorization header")
	}
}
