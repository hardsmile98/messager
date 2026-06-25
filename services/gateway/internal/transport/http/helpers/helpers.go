package helpers

import (
	"net/http"
)

const maxDeviceLength = 255

func DeviceFromRequest(r *http.Request) string {
	device := r.Header.Get("User-Agent")

	if device == "" {
		return "unknown"
	}

	if len(device) > maxDeviceLength {
		return device[:maxDeviceLength]
	}

	return device
}
