package pg

import (
	"net"

	"github.com/cyberark/secretless-broker/pkg/secretless/plugin/connector"
	"github.com/cyberark/secretless-broker/pkg/secretless/plugin/connector/tcp"
)

// NewConnector returns a tcp.Connector which returns an authenticated
// connection to a target service for each incoming client connection. It is a
// required method on the tcp.Plugin interface.
func NewConnector(conRes connector.Resources) tcp.Connector {
	newConnectorFunc := func(
		clientConn net.Conn,
		credentialValuesByID connector.CredentialValuesByID,
	) (backendConn net.Conn, err error) {
		// singleUseConnector is responsible for generating the authenticated connection
		// to the target service for each incoming client connection
		singleUseConnector := &SingleUseConnector{
			logger:   conRes.Logger(),
		}

		return singleUseConnector.Connect(clientConn, credentialValuesByID)
	}

	return tcp.ConnectorFunc(newConnectorFunc)
}

// PluginInfo is required as part of the Secretless plugin spec. It provides
// important metadata about the plugin.
func PluginInfo() map[string]string {
	return map[string]string{
		"pluginAPIVersion": "0.1.0",
		"type":             "connector.tcp",
		"id":               "pg",
		"description":      "returns an authenticated connection to a PostgreSQL database",
	}
}

// GetTCPPlugin is required as part of the Secretless plugin spec for TCP connector
// plugins. It returns the TCP plugin.
func GetTCPPlugin() tcp.Plugin {
	return tcp.ConnectorConstructor(NewConnector)
}
