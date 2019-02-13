package helpers

import (
	"os"

	"github.com/oracle/oci-go-sdk/common"
)

var (
	privateIpAddress,publicIpAddress string
)

// ParseEnvironmentVariables parses all required variables
func ParseEnvironmentVariables() {
	privateIpAddress = os.Getenv("OCI_PRIVATE_IP_ADDRESS")
	publicIpAddress = os.Getenv("OCI_PUBLIC_IP_ADDRESS")
}

// PrivateIpAddress returns the Private IP address to be moved
func PrivateIpAddress() *string {
	return common.String(privateIpAddress)
}

// PublicIpAddress returns the Public IP address to be moved
func PublicIpAddress() *string {
	return common.String(publicIpAddress)
}

