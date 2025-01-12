package helper

import (
	"bytes"
	"cynxhostagent/internal/model/response/responsecode"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"net/http"
	"strings"
	"text/template"

	"github.com/go-playground/validator/v10"
)

func DecodeAndValidateRequest(r *http.Request, dst interface{}, v *validator.Validate) error {
	if err := json.NewDecoder(r.Body).Decode(dst); err != nil {
		return errors.New("invalid request payload: " + err.Error())
	}

	if err := v.Struct(dst); err != nil {
		return errors.New("validation failed: " + err.Error())
	}

	return nil
}

func WriteJSONResponse(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if data == nil {
		return
	}

	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, "Failed to decode request: "+err.Error(), http.StatusBadRequest)
	}
}

func GetResponseCodeName(code responsecode.ResponseCode) string {
	if name, exists := responsecode.ResponseCodeNames[code]; exists {
		return name
	}
	return "Unknown Code"
}

// replacePlaceholders replaces {{}} placeholders in a script with real values.
func ReplacePlaceholders(script string, variables map[string]string) (string, error) {
	tmpl, err := template.New("script").Parse(script)
	if err != nil {
		return "", fmt.Errorf("failed to parse template: %w", err)
	}

	var output bytes.Buffer
	if err := tmpl.Execute(&output, variables); err != nil {
		return "", fmt.Errorf("failed to execute template: %w", err)
	}

	return output.String(), nil
}

func GetClientIP(r *http.Request) string {
	// If the request is behind a reverse proxy, the IP address might be forwarded in the X-Forwarded-For header.
	// First, check for the X-Forwarded-For header.
	ips := r.Header.Get("X-Forwarded-For")
	if ips != "" {
		// The X-Forwarded-For header contains a comma-separated list of IPs
		// The first IP in the list is the original client IP.
		return strings.Split(ips, ",")[0]
	}

	// Otherwise, fallback to the remote address.
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return ip
}
