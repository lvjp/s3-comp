package client

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	client, err := New(Config{})
	require.NoError(t, err)
	require.NotNil(t, client)
	require.Same(t, http.DefaultClient, client.config.HTTPClient)
}
