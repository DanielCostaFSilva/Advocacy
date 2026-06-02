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
	os.Setenv("JWT_SECRET", "test-secret")
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
	assert.Equal(t, "test-secret", cfg.JWTSecret)
	assert.Equal(t, 60, cfg.JWTExpirationMinutes)
}

func TestConfig_ShouldUseDefaults(t *testing.T) {
	os.Clearenv()
	os.Setenv("DB_HOST", "localhost")
	os.Setenv("DB_PORT", "5432")
	os.Setenv("DB_USER", "postgres")
	os.Setenv("DB_PASSWORD", "postgres")
	os.Setenv("DB_NAME", "legalflow")
	os.Setenv("JWT_SECRET", "my-secret")
	defer os.Clearenv()

	cfg, err := Load()
	require.NoError(t, err)

	assert.Equal(t, "law-office-api", cfg.AppName)
	assert.Equal(t, "local", cfg.AppEnv)
	assert.Equal(t, "8080", cfg.AppPort)
	assert.Equal(t, 60, cfg.JWTExpirationMinutes)
}

func TestConfig_ShouldFailWhenRequiredVarsMissing(t *testing.T) {
	os.Clearenv()
	_, err := Load()
	require.Error(t, err)
}

func TestConfig_ShouldFailWhenJWTSecretMissing(t *testing.T) {
	os.Clearenv()
	os.Setenv("DB_HOST", "localhost")
	defer os.Clearenv()

	_, err := Load()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "JWT_SECRET")
}
