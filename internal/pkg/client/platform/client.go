package platform

import (
	"context"
	"fmt"

	"github.com/IBM/platform-services-go-sdk/iamidentityv1"
	"github.com/IBM/platform-services-go-sdk/resourcecontrollerv2"

	"github.com/PDeXchange/pac/internal/pkg/client/iam"
)

type Client struct {
	resourceControllerClient *resourcecontrollerv2.ResourceControllerV2
	iamIdentityClient        *iamidentityv1.IamIdentityV1
}

func (c *Client) GetResourceInstance(ctx context.Context, id string) (*resourcecontrollerv2.ResourceInstance, error) {
	instance, _, err := c.resourceControllerClient.GetResourceInstanceWithContext(ctx, &resourcecontrollerv2.GetResourceInstanceOptions{ID: &id})
	return instance, err
}

func NewClient() (*Client, error) {
	auth, err := iam.GetIAMAuth()
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve authenticator: %w", err)
	}

	rcClient, err := resourcecontrollerv2.NewResourceControllerV2(&resourcecontrollerv2.ResourceControllerV2Options{Authenticator: auth})
	if err != nil {
		return nil, fmt.Errorf("failed to create resource controller client: %w", err)
	}

	iamClient, err := iamidentityv1.NewIamIdentityV1(&iamidentityv1.IamIdentityV1Options{
		Authenticator: auth,
		URL:           iamidentityv1.DefaultServiceURL,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create iam identity client: %w", err)
	}

	return &Client{
		resourceControllerClient: rcClient,
		iamIdentityClient:        iamClient,
	}, nil
}

// NewIAMIdentityClient creates iam identity client.
func NewIAMIdentityClient() (*iamidentityv1.IamIdentityV1, error) {
	auth, err := iam.GetIAMAuth()
	if err != nil {
		return nil, err
	}
	return iamidentityv1.NewIamIdentityV1(&iamidentityv1.IamIdentityV1Options{
		Authenticator: auth,
		URL:           iamidentityv1.DefaultServiceURL,
	})
}
