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
	contentType := r.Header.Get("Content-Type")
	if strings.Contains(contentType, "multipart/form-data") {
		return decodeMultipartForm(r, dst, v)
	}

	// If not multipart, proceed with normal JSON decoding
	if err := json.NewDecoder(r.Body).Decode(dst); err != nil {
		return errors.New("invalid request payload: " + err.Error())
	}

	// Validate the struct
	if err := v.Struct(dst); err != nil {
		return errors.New("validation failed: " + err.Error())
	}

	return nil
}

func decodeMultipartForm(r *http.Request, dst interface{}, v *validator.Validate) error {
	if err := r.ParseMultipartForm(10 << 20); err != nil { // 10MB limit
		return errors.New("failed to parse multipart form: " + err.Error())
	}

	// Convert request struct (dst) into a map to populate fields dynamically
	data := make(map[string]string)

	// Populate form fields (excluding files)
	for key, values := range r.Form {
		if len(values) > 0 {
			data[key] = values[0] // Use first value for simplicity
		}
	}

	// Convert `data` back to JSON and decode it into `dst`
	jsonData, _ := json.Marshal(data)
	if err := json.Unmarshal(jsonData, dst); err != nil {
		return errors.New("failed to decode form data: " + err.Error())
	}

	// Validate the struct
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
