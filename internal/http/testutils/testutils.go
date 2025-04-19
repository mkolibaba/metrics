package testutils

import (
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"testing"
)

func ReadAndCloseResponseBody(t *testing.T, response *http.Response) string {
	if response.Body == nil {
		return ""
	}

	defer response.Body.Close()
	body, err := io.ReadAll(response.Body)
	require.NoError(t, err)
	return string(body)
}
