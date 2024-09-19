//go:build unit

package api

import (
	"context"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

type testResponseWriter struct {
	w          http.ResponseWriter
	StatusCode int
}

func (trw *testResponseWriter) Header() http.Header {
	return trw.w.Header()
}

func (trw *testResponseWriter) Write(b []byte) (int, error) {
	return trw.w.Write(b)
}

func (trw *testResponseWriter) WriteHeader(statusCode int) {
	trw.StatusCode = statusCode

	trw.w.WriteHeader(statusCode)
}

func NewTestResponseWriter(w http.ResponseWriter) *testResponseWriter {
	return &testResponseWriter{w: w}
}

func wrapHandlerFunction(t *testing.T, wrappedHandler http.Handler, wantStatus int) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		trw := NewTestResponseWriter(writer)
		wrappedHandler.ServeHTTP(trw, request)

		if trw.StatusCode != wantStatus {
			t.Errorf("AuthenticateUser() = %d, want %d", trw.StatusCode, wantStatus)
		}
	})
}

func TestAuthenticateUser(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	userUUID, err := uuid.NewUUID()
	require.NoError(t, err)

	wrongUserUUID, err := uuid.NewUUID()
	require.NoError(t, err)

	correctHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	wrongUserHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	noHeaderHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	users := map[string]bool{
		userUUID.String(): true,
	}

	tests := map[string]struct {
		h          http.HandlerFunc
		clientID   uuid.UUID
		wantStatus int
	}{
		"successfully processed ": {
			h:          correctHandler,
			clientID:   userUUID,
			wantStatus: http.StatusOK,
		},
		"failed due to lacking clientID in header": {
			h:          noHeaderHandler,
			clientID:   uuid.Nil,
			wantStatus: http.StatusForbidden,
		},
		"failed due to the fact that user in header is not in users map": {
			h:          wrongUserHandler,
			clientID:   wrongUserUUID,
			wantStatus: http.StatusForbidden,
		},
	}
	for tName, tt := range tests {
		t.Run(tName, func(t *testing.T) {
			request, innerErr := http.NewRequestWithContext(context.Background(), http.MethodGet, "test", nil)
			require.NoError(t, innerErr)

			if tt.clientID != uuid.Nil {
				request.Header.Set("Client-Id", tt.clientID.String())
			}

			responseRecorder := httptest.NewRecorder()

			wrapHandlerFunction(
				t,
				AuthenticateUser(tt.h, logger, users),
				tt.wantStatus,
			).ServeHTTP(responseRecorder, request)
		})
	}
}
