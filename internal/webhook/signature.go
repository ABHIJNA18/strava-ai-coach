// This file verifies Strava webhook signatures using HMAC-SHA256.

package webhook

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"strconv"
	"strings"
	"time"
)

const webhookSignatureTolerance = 5 * time.Minute

func verifySignature(
	header string,
	body []byte,
	secret string,
	now time.Time,
) bool {
	if header == "" || secret == "" {
		return false
	}

	var timestamp string
	var signatureHex string

	for _, part := range strings.Split(header, ",") {
		key, value, ok := strings.Cut(strings.TrimSpace(part), "=")
		if !ok {
			continue
		}

		switch key {
		case "t":
			timestamp = value
		case "v1":
			signatureHex = value
		}
	}

	if timestamp == "" || signatureHex == "" {
		return false
	}

	timestampSeconds, err := strconv.ParseInt(timestamp, 10, 64)
	if err != nil {
		return false
	}

	eventTime := time.Unix(timestampSeconds, 0)
	difference := now.Sub(eventTime)

	if difference < 0 {
		difference = -difference
	}

	if difference > webhookSignatureTolerance {
		return false
	}

	providedSignature, err := hex.DecodeString(signatureHex)
	if err != nil {
		return false
	}

	mac := hmac.New(sha256.New, []byte(secret))
	_, _ = mac.Write([]byte(timestamp + "."))
	_, _ = mac.Write(body)

	expectedSignature := mac.Sum(nil)

	return hmac.Equal(providedSignature, expectedSignature)
}