package mssql

import (
	"context"
	"net"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/cyberark/secretless-broker/internal/plugin/connectors/tcp/mssql/mock"
	logmock "github.com/cyberark/secretless-broker/pkg/secretless/log/mock"
	pluginconnector "github.com/cyberark/secretless-broker/pkg/secretless/plugin/connector"
	"github.com/denisenkom/go-mssqldb"
)

func TestHappyPath(t *testing.T) {
	logger := logmock.NewLogger()
	clientConn := mock.NewNetConn(nil)
	expectedBackendConn := mock.NewNetConn(nil)
	creds := pluginconnector.CredentialValuesByID{
		"credName": []byte("secret"),
	}

	ctor := mock.NewSuccessfulMSSQLConnectorCtor(
		func(ctx context.Context) (net.Conn, error) {
			interceptor := mssql.ConnectInterceptorFromContext(ctx)

			interceptor.ServerPreLoginResponse <- map[uint8][]byte{ 0: {0, 0} }

			<- interceptor.ClientLoginRequest

			interceptor.ServerLoginResponse <- &mssql.LoginResponse{}

			return expectedBackendConn, nil
		},
	)

	connector := NewSingleUseConnectorWithOptions(
		logger,
		ctor,
		mock.SuccessfulReadPreloginRequest,
		mock.SuccessfulWritePreloginResponse,
		mock.SuccessfulReadLoginRequest,
		mock.SuccessfulWriteLoginResponse,
		mock.SuccessfulWriteError,
		mock.FakeTdsBufferCtor,
	)
	actualBackendConn, err := connector.Connect(clientConn, creds)

	assert.Nil(t, err)
	if err != nil {
		return
	}

	assert.Equal(t, expectedBackendConn, actualBackendConn)
}
/*
Test cases not to forget

	- This is not exhaustive.  Use the method of tracing all the code paths (for
	each error condition, assuming the previous succeeded) and add a test for
	each.  If that becomes too many, use judgment to eliminate less important
	ones.

	- While we shouldn't test the logger messages extensively, here is one that
	we should test: sending an error message to the user fails.  We want to make
	sure those are logged.
*/