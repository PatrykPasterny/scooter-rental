package api

import (
	"log/slog"
	"net/http"
)

func AuthenticateUser(h http.HandlerFunc, logger *slog.Logger, users map[string]bool) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		clientUUID, err := clientUUIDFromHeader(request)
		if err != nil {
			logger.Error("Failed to parse the clientID", err)

			Error(writer, http.StatusForbidden, "Failed getting clientID from header.")

			return
		}

		if _, ok := users[clientUUID.String()]; !ok {
			logger.Error("Failed to authenticate user", err)

			Error(writer, http.StatusForbidden, "Failed getting clientID from header.")

			return
		}

		h(writer, request)
	}
}
