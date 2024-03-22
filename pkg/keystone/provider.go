package keystone

import (
	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack"
	"github.com/selectel/external-dns-stackit-webhook/pkg/httpdefault"
	"go.uber.org/zap"
)

func defaultOSClient(endpoint string) (*gophercloud.ProviderClient, error) {
	client, err := openstack.NewClient(endpoint)
	client.HTTPClient = httpdefault.Client()
	client.UserAgent.Prepend(httpdefault.UserAgent)

	return client, err
}

type Credentials struct {
	// IdentityEndpoint is an API endpoint to authorization. It is OS_AUTH_URL variable from rc.sh.
	IdentityEndpoint string
	// AccountID is Selectel account ID of the user. It is OS_PROJECT_DOMAIN_NAME variable from rc.sh.
	AccountID string
	// ProjectID is Selectel project ID of the user. It is OS_PROJECT_ID variable from rc.sh.
	ProjectID string
	// Username is service user's name. It is OS_USERNAME variable from rc.sh.
	Username string
	// Password is service user's password. It is OS_PASSWORD variable from rc.sh.
	Password string
}

type Provider struct {
	logger *zap.Logger
	// credentials contains data to access openstack identity API.
	credentials Credentials
}

// GetToken returns keystone token that may be used to authorize requests to Selectel API.
// It generates new token for each call.
func (p Provider) GetToken() (string, error) {
	p.logger.Info(
		"getting keystone token",
		zap.String("identity_endpoint", p.credentials.IdentityEndpoint),
		zap.String("username", p.credentials.Username),
		zap.String("account_id", p.credentials.AccountID),
		zap.String("project_id", p.credentials.ProjectID),
	)

	opts := gophercloud.AuthOptions{
		IdentityEndpoint: p.credentials.IdentityEndpoint,
		Username:         p.credentials.Username,
		Password:         p.credentials.Password,
		DomainName:       p.credentials.AccountID,
		Scope: &gophercloud.AuthScope{
			ProjectID: p.credentials.ProjectID,
		},
	}

	p.logger.Debug("connecting to identity endpoint")
	client, err := defaultOSClient(p.credentials.IdentityEndpoint)
	if err != nil {
		p.logger.Error("error during creating default openstack client", zap.Error(err))

		return "", err
	}
	err = openstack.Authenticate(client, opts)

	return client.Token(), err
}

func NewProvider(logger *zap.Logger, credentials Credentials) *Provider {
	return &Provider{
		logger:      logger,
		credentials: credentials,
	}
}
