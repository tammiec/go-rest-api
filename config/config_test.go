package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

var varName = "MYTESTVARIABLE"

func TestBadConfig(t *testing.T) {
	os.Setenv("HTTP_PORT", "NOT A VALID PORT")
	require.Panics(t, func() { Loader() })
}

func TestConfigBool(t *testing.T) {
	os.Setenv(varName, "TRUE")
	require.True(t, MustGetEnvBool(varName))
	os.Setenv(varName, "False")
	require.False(t, MustGetEnvBool(varName))
	os.Setenv(varName, "gobledgook")
	require.Panics(t, func() { MustGetEnvBool(varName) })
}

func TestConfigInt(t *testing.T) {
	os.Setenv(varName, "4")
	require.Equal(t, 4, MustGetEnvInt(varName))
	os.Setenv(varName, "-2")
	require.Equal(t, -2, MustGetEnvInt(varName))
	os.Setenv(varName, "gobledgook")
	require.Panics(t, func() { MustGetEnvInt(varName) })
}
func TestConfigString(t *testing.T) {
	os.Setenv(varName, "Böb")
	require.Equal(t, "Böb", MustGetEnvString(varName))
	os.Setenv(varName, "")
	require.Panics(t, func() { MustGetEnvString(varName) })
	os.Setenv(varName, "  ")
	require.Panics(t, func() { MustGetEnvString(varName) })
}
