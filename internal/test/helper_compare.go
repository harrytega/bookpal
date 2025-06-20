package test

import (
	"crypto/sha256"
	"fmt"
	"io"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"test-project/internal/api/httperrors"
)

func CompareFileHashes(t *testing.T, expectedFile string, actualFile string) {
	t.Helper()

	ef, err := os.Open(expectedFile)
	require.NoError(t, err)
	defer ef.Close()

	eh := sha256.New()
	_, err = io.Copy(eh, ef)
	require.NoError(t, err)

	af, err := os.Open(actualFile)
	require.NoError(t, err)
	defer af.Close()

	ah := sha256.New()
	_, err = io.Copy(ah, af)
	require.NoError(t, err)

	assert.Equal(t, eh.Sum(nil), ah.Sum(nil))
}

func CompareAllPayload(t *testing.T, base map[string]interface{}, toCheck map[string]string, skipKeys map[string]bool, keyConvertFunc ...func(string) string) {
	var keyFunc func(string) string
	if len(keyConvertFunc) > 0 {
		keyFunc = keyConvertFunc[0]
	} else {
		keyFunc = func(s string) string {
			return s
		}
	}
	for k, v := range base {
		if skipKeys[k] {
			continue
		}

		strV := fmt.Sprintf("%v", v)
		//revive:disable-next-line:var-naming
		//nolint:revive
		kConv := keyFunc(k)
		compareStrV := fmt.Sprintf("%v", toCheck[kConv])

		// checks with contains because of dateTime and null.String struct
		contains := strings.Contains(compareStrV, strV)
		assert.Truef(t, contains, "Expected for %s: '%s'. Got: '%s'", k, strV, compareStrV)
	}
}

func CompareAll(t *testing.T, base map[string]string, toCheck map[string]string, skipKeys map[string]bool) {
	for k, v := range base {
		if skipKeys[k] {
			continue
		}

		strV := fmt.Sprintf("%v", v)
		compareStrV := fmt.Sprintf("%v", toCheck[k])

		// checks with contains because of dateTime and null.String struct
		contains := strings.Contains(compareStrV, strV)
		assert.Truef(t, contains, "Expected for %s: '%s'. Got: '%s'", k, strV, compareStrV)
	}
}

func RequireHTTPError(t *testing.T, res *httptest.ResponseRecorder, httpError *httperrors.HTTPError) httperrors.HTTPError {
	t.Helper()

	if httpError.Code != nil {
		require.Equal(t, int(*httpError.Code), res.Result().StatusCode)
	}

	var response httperrors.HTTPError
	ParseResponseAndValidate(t, res, &response)

	require.Equal(t, httpError, &response)

	return response
}
