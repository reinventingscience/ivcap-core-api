// $ goa gen github.com/reinventingscience/ivcap-core-api/design

package client

import (
	"fmt"
)

// ListServicePath returns the URL path to the service service list HTTP endpoint.
func ListServicePath() string {
	return "/1/services"
}

// CreateServicePath returns the URL path to the service service create HTTP endpoint.
func CreateServicePath() string {
	return "/1/services"
}

// ReadServicePath returns the URL path to the service service read HTTP endpoint.
func ReadServicePath(id string) string {
	return fmt.Sprintf("/1/services/%v", id)
}

// UpdateServicePath returns the URL path to the service service update HTTP endpoint.
func UpdateServicePath(id string) string {
	return fmt.Sprintf("/1/services/%v", id)
}

// DeleteServicePath returns the URL path to the service service delete HTTP endpoint.
func DeleteServicePath(id string) string {
	return fmt.Sprintf("/1/services/%v", id)
}
