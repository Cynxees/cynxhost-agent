package middleware

import (
	"context"
	"cynxhostagent/internal/helper"
	"cynxhostagent/internal/model/response"
	"cynxhostagent/internal/model/response/responsecode"
	"fmt"
	"net/http"
	"reflect"
	"strconv"
	"time"
)

type HandlerFuncWithHelper func(w http.ResponseWriter, r *http.Request) (context.Context, response.APIResponse)

func WrapHandler(handler HandlerFuncWithHelper, debug bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Execute the handler and capture response or error
		ctx, apiResponse := handler(w, r)

		if !debug && apiResponse.Error != "" {
			// apiResponse.Error = "Hidden Error, Debug Mode is Disabled"
		}

		if apiResponse.Code != "" {
			apiResponse.CodeName = helper.GetResponseCodeName(apiResponse.Code)
		}

		if apiResponse.Code != responsecode.CodeSuccess {
			helper.WriteJSONResponse(w, http.StatusOK, apiResponse)
			return
		}

		// Apply visibility filter to the response
		level := getUserVisibilityLevel(ctx)
		filteredData, err := applyVisibilityFilter(apiResponse.Data, level)
		if err != nil {
			apiResponse.Error = "Failed to process response data"
		} else {
			apiResponse.Data = filteredData
		}

		// Write the successful response
		helper.WriteJSONResponse(w, http.StatusOK, apiResponse)
	}
}

func getUserVisibilityLevel(ctx context.Context) int {
	// Get the visibility level from the request header
	level, ok := helper.GetVisibilityLevelFromContext(ctx)
	if !ok {
		return 1 // Default to the lowest level if not provided
	}

	return level
}

// applyVisibilityFilter automatically filters struct fields based on visibility levels.
func applyVisibilityFilter(input interface{}, level int) (interface{}, error) {
	// Use reflection to handle nested structs dynamically
	v := reflect.ValueOf(input)
	t := reflect.TypeOf(input)

	// Handle pointers
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
		t = t.Elem()
	}

	// If it's not a struct, slice, array, or map, return it as-is
	if v.Kind() != reflect.Struct && v.Kind() != reflect.Slice && v.Kind() != reflect.Array && v.Kind() != reflect.Map {
		return input, nil
	}

	// Map to hold filtered fields
	filtered := make(map[string]interface{})

	// Iterate over struct fields or map entries
	switch v.Kind() {
	case reflect.Struct:
		for i := 0; i < v.NumField(); i++ {
			field := t.Field(i)
			fieldValue := v.Field(i)

			visibilityTag := field.Tag.Get("visibility")
			if visibilityTag == "" || level >= parseVisibility(visibilityTag) {
				// Special case for time.Time
				if fieldValue.Kind() == reflect.Struct && fieldValue.Type() == reflect.TypeOf(time.Time{}) {
					// Handle time.Time separately, we just pass it as-is
					filtered[field.Name] = fieldValue.Interface()
				} else if fieldValue.Kind() == reflect.Struct {
					// Recursively handle nested structs
					nestedFiltered, err := applyVisibilityFilter(fieldValue.Interface(), level)
					if err != nil {
						return nil, err
					}
					filtered[field.Name] = nestedFiltered
				} else if fieldValue.Kind() == reflect.Slice || fieldValue.Kind() == reflect.Array {
					// Handle slices or arrays
					var filteredSlice []interface{}
					for j := 0; j < fieldValue.Len(); j++ {
						element := fieldValue.Index(j).Interface()
						filteredElement, err := applyVisibilityFilter(element, level)
						if err != nil {
							return nil, err
						}
						filteredSlice = append(filteredSlice, filteredElement)
					}
					filtered[field.Name] = filteredSlice
				} else if fieldValue.Kind() == reflect.Map {
					// Handle map types
					filteredMap := make(map[string]interface{})
					for _, key := range fieldValue.MapKeys() {
						mapKey := key.Interface()
						mapValue := fieldValue.MapIndex(key).Interface()

						// Apply the filter for both keys and values
						filteredKey, err := applyVisibilityFilter(mapKey, level)
						if err != nil {
							return nil, err
						}
						filteredValue, err := applyVisibilityFilter(mapValue, level)
						if err != nil {
							return nil, err
						}
						filteredMap[fmt.Sprintf("%v", filteredKey)] = filteredValue
					}
					filtered[field.Name] = filteredMap
				} else {
					// Add non-struct, non-array, non-map types directly
					filtered[field.Name] = fieldValue.Interface()
				}
			}
		}
	case reflect.Map:
		// Handle map filtering recursively
		filteredMap := make(map[string]interface{})
		for _, key := range v.MapKeys() {
			mapKey := key.Interface()
			mapValue := v.MapIndex(key).Interface()

			// Apply the filter for both keys and values
			filteredKey, err := applyVisibilityFilter(mapKey, level)
			if err != nil {
				return nil, err
			}
			filteredValue, err := applyVisibilityFilter(mapValue, level)
			if err != nil {
				return nil, err
			}
			filteredMap[fmt.Sprintf("%v", filteredKey)] = filteredValue
		}
		return filteredMap, nil
	case reflect.Slice, reflect.Array:
		// Handle slice and array filtering recursively
		var filteredSlice []interface{}
		for j := 0; j < v.Len(); j++ {
			element := v.Index(j).Interface()
			filteredElement, err := applyVisibilityFilter(element, level)
			if err != nil {
				return nil, err
			}
			filteredSlice = append(filteredSlice, filteredElement)
		}
		return filteredSlice, nil
	}

	return filtered, nil
}

func parseVisibility(tag string) int {
	// Parse the "visibility" tag and return the corresponding visibility level
	level, err := strconv.Atoi(tag)
	if err != nil {
		return 3 // Default to the lowest level if parsing fails
	}
	return level
}
