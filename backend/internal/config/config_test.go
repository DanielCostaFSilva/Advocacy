package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConfig_ShouldLoadFromEnv(t *testing.T) {
	os.Clearenv()
	os.Setenv("APP_NAME", "test-api")
	os.Setenv("APP_ENV", "testing")
	os.Setenv("APP_PORT", "9090")
	os.Setenv("DB_HOST", "dbhost")
	os.Setenv("DB_PORT", "5432")
	os.Setenv("DB_USER", "dbuser")
	os.Setenv("DB_PASSWORD", "dbpass")
	os.Setenv("DB_NAME", "dbname")
	defer os.Clearenv()

	cfg, err := Load()
	require.NoError(t, err)

	assert.Equal(t, "test-api", cfg.AppName)
	assert.Equal(t, "testing", cfg.AppEnv)
	assert.Equal(t, "9090", cfg.AppPort)
	assert.Equal(t, "dbhost", cfg.DBHost)
	assert.Equal(t, "5432", cfg.DBPort)
	assert.Equal(t, "dbuser", cfg.DBUser)
	assert.Equal(t, "dbpass", cfg.DBPassword)
	assert.Equal(t, "dbname", cfg.DBName)
}

func TestConfig_ShouldUseDefaults(t *testing.T) {
	os.Clearenv()
	os.Setenv("DB_HOST", "localhost")
	os.Setenv("DB_PORT", "5432")
	os.Setenv("DB_USER", "postgres")
	os.Setenv("DB_PASSWORD", "postgres")
	os.Setenv("DB_NAME", "legalflow")
	defer os.Clearenv()

	cfg, err := Load()
	require.NoError(t, err)

	assert.Equal(t, "law-office-api", cfg.AppName)
	assert.Equal(t, "local", cfg.AppEnv)
	assert.Equal(t, "8080", cfg.AppPort)
}

func TestConfig_ShouldFailWhenDBHostMissing(t *testing.T) {
	os.Clearenv()
	_, err := Load()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "DB_HOST")
}
