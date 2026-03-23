package delete_test

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/require"

	"url-shortener/internal/http-server/handlers/url/delete"
	"url-shortener/internal/http-server/handlers/url/delete/mocks"
	"url-shortener/internal/lib/api/response"
	"url-shortener/internal/lib/logger/handlers/slogdiscard"
	"url-shortener/internal/storage"
)

func TestDeleteHandler(t *testing.T) {
	cases := []struct {
		name      string
		alias     string
		url       string
		expectedStatus int
		respError string
		mockError error
	}{
		{
			name:  "Success",
			alias: "test_alias",
			expectedStatus: http.StatusOK,
		},
		{
			name:      "Alias not found",
			alias:     "non_existent_alias",
			respError: "alias not found",
			mockError: storage.ErrAliasNotFound,
			expectedStatus: http.StatusNotFound,
		},
		{
			name:      "DeleteURL Error",
			alias:     "test_alias",
			respError: "failed to delete url",
			mockError: errors.New("unexpected error"),
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tc := range cases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			urlDeleterMock := mocks.NewURLDeleter(t)

			if tc.respError == "" || tc.mockError != nil {
				urlDeleterMock.On("DeleteURL", tc.alias).
					Return(tc.mockError).
					Once()
			}

			handler := delete.New(slogdiscard.NewDiscardLogger(), urlDeleterMock)

			req, err := http.NewRequest(http.MethodDelete, "/url/"+tc.alias, nil)
			require.NoError(t, err)

			router := chi.NewRouter()
			router.Delete("/url/{alias}", handler)

			rr := httptest.NewRecorder()
			router.ServeHTTP(rr, req)

			require.Equal(t, rr.Code, tc.expectedStatus)

			body := rr.Body.String()

			var resp response.Response

			require.NoError(t, json.Unmarshal([]byte(body), &resp))

			require.Equal(t, tc.respError, resp.Error)

			// TODO: add more checks
		})
	}
}