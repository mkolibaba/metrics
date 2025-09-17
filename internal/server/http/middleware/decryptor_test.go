package middleware

import (
	"fmt"
	"github.com/mkolibaba/metrics/internal/server/testutils"
	"go.uber.org/zap/zaptest"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestDecryptor(t *testing.T) {
	run := func(t *testing.T, decryptor BodyDecryptor, body io.Reader) (*httptest.ResponseRecorder, string) {
		t.Helper()

		var decryptedBody string
		h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			data, err := io.ReadAll(r.Body)
			testutils.AssertNoError(t, err)
			decryptedBody = string(data)
		})
		mw := Decryptor(decryptor, zaptest.NewLogger(t).Sugar())

		request := httptest.NewRequest("POST", "/test", body)
		recorder := httptest.NewRecorder()

		mw(h).ServeHTTP(recorder, request)

		return recorder, decryptedBody
	}
	successfulDecryptor := &BodyDecryptorMock{
		DecryptFunc: func(data []byte) ([]byte, error) {
			return []byte("decrypted: " + string(data)), nil
		},
	}

	t.Run("success", func(t *testing.T) {
		requestBody := "some data"
		wantRequestBody := "decrypted: some data"
		var decryptedRequestBody string

		recorder, decryptedBody := run(t, successfulDecryptor, strings.NewReader(requestBody))

		testutils.AssertResponseStatusCode(t, http.StatusOK, recorder.Code)
		if decryptedBody != wantRequestBody {
			t.Errorf("error decrypted request body: want %s, got %s", wantRequestBody, decryptedRequestBody)
		}
	})
	t.Run("fail_read_body", func(t *testing.T) {
		recorder, _ := run(t, successfulDecryptor, testutils.AlwaysFailingReader)

		testutils.AssertResponseStatusCode(t, http.StatusInternalServerError, recorder.Code)
	})
	t.Run("fail_decrypt", func(t *testing.T) {
		dcr := &BodyDecryptorMock{
			DecryptFunc: func(data []byte) ([]byte, error) {
				return nil, fmt.Errorf("failed to decrypt")
			},
		}

		recorder, _ := run(t, dcr, strings.NewReader("some data"))

		testutils.AssertResponseStatusCode(t, http.StatusInternalServerError, recorder.Code)
	})
}
