package controller

import (
	"net/http/httptest"
	"testing"
)

func TestGetClientIP(t *testing.T) {
	tests := []struct {
		name           string
		forwardedFor   string
		remoteAddr     string
		expectedClient string
	}{
		{
			name:           "uses first forwarded address",
			forwardedFor:   "198.51.100.24, 10.0.0.3",
			remoteAddr:     "169.254.129.1:8080",
			expectedClient: "198.51.100.24",
		},
		{
			name:           "uses forwarded IPv6 address",
			forwardedFor:   "2001:db8::1, 10.0.0.3",
			remoteAddr:     "169.254.129.1:8080",
			expectedClient: "2001:db8::1",
		},
		{
			name:           "falls back to remote address",
			remoteAddr:     "169.254.129.1:8080",
			expectedClient: "169.254.129.1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := httptest.NewRequest("GET", "/", nil)
			r.RemoteAddr = tt.remoteAddr
			r.Header.Set("X-Forwarded-For", tt.forwardedFor)

			if got := getClientIP(r); got != tt.expectedClient {
				t.Errorf("getClientIP() = %q, want %q", got, tt.expectedClient)
			}
		})
	}
}
