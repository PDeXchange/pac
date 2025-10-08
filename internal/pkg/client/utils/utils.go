package utils

import (
	"context"
	"fmt"
	"regexp"

	"github.com/IBM/go-sdk-core/v5/core"
	"github.com/IBM/platform-services-go-sdk/iamidentityv1"
	"github.com/PDeXchange/pac/internal/pkg/client/platform"
)

// GetAccountID returns IBM cloud account ID of API key used.
func GetAccountID(ctx context.Context, auth *core.IamAuthenticator) (string, error) {
	iamv1, err := platform.NewIAMIdentityClient()
	if err != nil {
		return "", err
	}

	apiKeyDetailsOpt := iamidentityv1.GetAPIKeysDetailsOptions{IamAPIKey: &auth.ApiKey}
	apiKey, _, err := iamv1.GetAPIKeysDetailsWithContext(ctx, &apiKeyDetailsOpt)
	if err != nil {
		return "", err
	}
	if apiKey == nil {
		return "", fmt.Errorf("could not retrieve account id")
	}

	return *apiKey.AccountID, nil
}

// resourceCRNRegexp - regex pattern for IBM Resource CRN
var resourceCRNRegexp = regexp.MustCompile(`^crn:v[0-9]:(?P<cloudName>[^:]*):(?P<cloudType>[^:]*):(?P<serviceName>[^:]*):(?P<location>[^:]*):(?P<scope>[^:]*):(?P<guid>[^:]*):(?P<resourceType>[^:]*):(?P<resourceID>[^:]*)$`)

// ValidateResourceCRN - validates provided IBM Resource CRN string
func ValidateResourceCRN(crn string) error {
	if !resourceCRNRegexp.MatchString(crn) {
		return fmt.Errorf("provided CRN is invalid")
	}
	return nil
}
