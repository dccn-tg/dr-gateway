package dr

import (
	"github.com/cyverse/go-irodsclient/irods/types"
)

type Config struct {
	IrodsHost           string
	IrodsPort           int
	IrodsZone           string
	IrodsUser           string
	IrodsPass           string
	OrganisationalUnits []string
}

func NewAccount(config Config) (*types.IRODSAccount, error) {

	account, err := types.CreateIRODSAccount(
		config.IrodsHost,
		config.IrodsPort,
		config.IrodsUser,
		config.IrodsZone,
		types.AuthSchemeNative,
		config.IrodsPass,
		"",
	)

	if err != nil {
		return nil, err
	}

	// sslConfig, err := types.CreateIRODSSSLConfig(
	// 	"/opt/irods/ssl/icat-prod.pem",
	// 	32,
	// 	"AES-256-CBC",
	// 	8,
	// 	16,
	// )
	// if err != nil {
	// 	return nil, err
	// }

	// account.SSLConfiguration = sslConfig
	// account.CSNegotiationPolicy = types.CSNegotiationRequireSSL
	// account.ClientServerNegotiation = true

	return account, nil
}

func NewProxyAccount(config Config, user string) (*types.IRODSAccount, error) {

	account, err := types.CreateIRODSProxyAccount(
		config.IrodsHost,
		config.IrodsPort,
		user,
		config.IrodsZone,
		config.IrodsUser,
		config.IrodsZone,
		types.AuthSchemeNative,
		config.IrodsPass,
		"",
	)

	if err != nil {
		return nil, err
	}

	// sslConfig, err := types.CreateIRODSSSLConfig(
	// 	"/opt/irods/ssl/icat-prod.pem",
	// 	32,
	// 	"AES-256-CBC",
	// 	8,
	// 	16,
	// )
	// if err != nil {
	// 	return nil, err
	// }

	// account.SSLConfiguration = sslConfig
	// account.CSNegotiationPolicy = types.CSNegotiationRequireSSL
	// account.ClientServerNegotiation = true

	return account, nil
}
