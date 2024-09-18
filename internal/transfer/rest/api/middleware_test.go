package api

import (
	"log/slog"
	"net/http"
	"reflect"
	"testing"
)

func TestAuthenticateUser(t *testing.T) {
	type args struct {
		h      http.HandlerFunc
		logger *slog.Logger
		users  map[string]bool
	}
	tests := []struct {
		name string
		args args
		want http.HandlerFunc
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := AuthenticateUser(tt.args.h, tt.args.logger, tt.args.users); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("AuthenticateUser() = %v, want %v", got, tt.want)
			}
		})
	}
}
