package dr

import (
	"strings"

	"github.com/cyverse/go-irodsclient/irods/types"
)

type Config struct {
	IrodsHost           string
	IrodsPort           int
	IrodsZone           string
	IrodsUser           string
	IrodsPass           string
	IrodsAuthScheme     string
	IrodsSslCacert      string
	IrodsSslKeysize     int
	IrodsSslAlgorithm   string
	IrodsSslSaltSize    int
	IrodsHashRounds     int
	OrganisationalUnits []ServiceAccount
}

// ServiceAccount defines the configuration data structure for
// OU-specific service account data-access credential.
type ServiceAccount struct {
	Name      string
	IrodsUser string
	IrodsPass string
}

func (c Config) AuthSchemeType() types.AuthScheme {
	switch strings.ToLower(c.IrodsAuthScheme) {
	case "pam":
		return types.AuthSchemePAM
	case "gsi":
		return types.AuthSchemeGSI
	default:
		return types.AuthSchemeNative
	}
}

// NewAccount returned a configured `types.IRODSAccount` iRODS client that is
// ready for making iRODS connections.
func NewAccount(config Config) (*types.IRODSAccount, error) {

	account, err := types.CreateIRODSAccount(
		config.IrodsHost,
		config.IrodsPort,
		config.IrodsUser,
		config.IrodsZone,
		config.AuthSchemeType(),
		config.IrodsPass,
		"",
	)

	if err != nil {
		return nil, err
	}

	sslConfig, err := types.CreateIRODSSSLConfig(
		config.IrodsSslCacert,
		config.IrodsSslKeysize,
		config.IrodsSslAlgorithm,
		config.IrodsSslSaltSize,
		config.IrodsHashRounds,
	)
	if err != nil {
		return nil, err
	}

	account.SSLConfiguration = sslConfig
	account.CSNegotiationPolicy = types.CSNegotiationRequireSSL
	account.ClientServerNegotiation = true

	return account, nil
}

// NewServiceAccounts returned a map of configured `types.IRODSAccount` iRODS clients that is
// ready for making iRODS connections.
//
// The key of the map is the organizational unit name provided by the `config.OrganizationalUnits`;
// and the value is the iRODS client initiated with the organizational unit's service account credential.
func NewServiceAccounts(config Config) (map[string]*types.IRODSAccount, error) {

	accounts := make(map[string]*types.IRODSAccount)

	sslConfig, err := types.CreateIRODSSSLConfig(
		config.IrodsSslCacert,
		config.IrodsSslKeysize,
		config.IrodsSslAlgorithm,
		config.IrodsSslSaltSize,
		config.IrodsHashRounds,
	)

	if err != nil {
		return nil, err
	}

	for _, sa := range config.OrganisationalUnits {
		account, err := types.CreateIRODSAccount(
			config.IrodsHost,
			config.IrodsPort,
			sa.IrodsUser,
			config.IrodsZone,
			config.AuthSchemeType(),
			sa.IrodsPass,
			"",
		)

		if err != nil {
			return nil, err
		}

		account.SSLConfiguration = sslConfig
		account.CSNegotiationPolicy = types.CSNegotiationRequireSSL
		account.ClientServerNegotiation = true

		accounts[sa.Name] = account
	}

	return accounts, nil
}
