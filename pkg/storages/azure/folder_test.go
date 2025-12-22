package azure

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/wal-g/wal-g/pkg/storages/storage"
)

func TestAzureFolder(t *testing.T) {
	t.Skip("Credentials needed to run Azure Storage tests")

	storageFolder, err := ConfigureFolder("azure://test-container/test-folder/Sub0",
		make(map[string]string))

	assert.NoError(t, err)

	storage.RunFolderTest(storageFolder, t)
}

var ConfigureAuthType = configureAuthType

func TestConfigureAccessKeyAuthType(t *testing.T) {
	settings := map[string]string{AccessKeySetting: "foo"}
	authType, accountToken, accessKey := ConfigureAuthType(settings)
	assert.Equal(t, authType, AzureAccessKeyAuth)
	assert.Empty(t, accountToken)
	assert.Equal(t, accessKey, "foo")
}

func TestConfigureSASTokenAuth(t *testing.T) {
	settings := map[string]string{SasTokenSetting: "foo"}
	authType, accountToken, accessKey := ConfigureAuthType(settings)
	assert.Equal(t, authType, AzureSASTokenAuth)
	assert.Equal(t, accountToken, "?foo")
	assert.Empty(t, accessKey)
}

func TestConfigureDefaultAuth(t *testing.T) {
	settings := make(map[string]string)
	authType, accountToken, accessKey := ConfigureAuthType(settings)
	assert.Empty(t, authType)
	assert.Empty(t, accountToken)
	assert.Empty(t, accessKey)
}

func TestConfigureManagedIdentityAuth(t *testing.T) {
	settings := map[string]string{MiTokenSetting: "foo"}
	authType, accountToken, accessKey := ConfigureAuthType(settings)
	assert.Equal(t, authType, AzureManagedIdentityAuth)
	assert.Equal(t, accountToken, "foo")
	assert.Empty(t, accessKey)
}

func TestConfigureManagedIdentityAuthWithClientID(t *testing.T) {
	settings := map[string]string{ClientIDSetting: "client-id-123"}
	authType, accountToken, accessKey := ConfigureAuthType(settings)
	assert.Equal(t, AzureManagedIdentityAuth, authType)
	assert.Empty(t, accountToken)
	assert.Empty(t, accessKey)
}

func TestGetContainerClientWithManagedIdentity_RequiresOneSetting(t *testing.T) {
	// Neither client ID nor MI token provided: should error
	client, err := getContainerClientWithManagedIndetity(
		"acct", "core.windows.net", "container", time.Minute, "", "",
	)
	assert.Error(t, err)
	assert.Nil(t, client)
}

func TestGetContainerClientWithManagedIdentity_WithClientID(t *testing.T) {
	// With client ID: should construct a client (no network calls involved)
	client, err := getContainerClientWithManagedIndetity(
		"acct", "core.windows.net", "container", time.Minute, "", "client-id-123",
	)
	assert.NoError(t, err)
	assert.NotNil(t, client)
}
