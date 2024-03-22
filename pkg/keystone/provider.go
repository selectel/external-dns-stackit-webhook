package keystone

import (
	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack"
	"github.com/selectel/external-dns-webhook/pkg/httpclient"
	"go.uber.org/zap"
)

func defaultOSClient(endpoint string) (*gophercloud.ProviderClient, error) {
	client, err := openstack.NewClient(endpoint)
	client.HTTPClient = httpclient.Default()
	client.UserAgent.Prepend(httpclient.DefaultUserAgent)

	return client, err
}

type Credentials struct {
	// IdentityEndpoint is OS_AUTH_URL variable from rc.sh.
	IdentityEndpoint string
	// AccountID is OS_PROJECT_DOMAIN_NAME variable from rc.sh.
	AccountID string
	// IdentityEndpoint is OS_PROJECT_ID variable from rc.sh.
	ProjectID string
	// IdentityEndpoint is OS_USERNAME variable from rc.sh.
	Username string
	// IdentityEndpoint is OS_PASSWORD variable from rc.sh.
	Password string
}

type Provider struct {
	logger *zap.Logger
	// credentials contains data to access openstack identity API.
	credentials Credentials
}

// GetToken returns valid keystone token. It will be stored for the next requests and then checked whether it is expired.
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
