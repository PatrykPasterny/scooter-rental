//go:build unit

package config

import (
	"context"
	"reflect"
	"testing"
)

func TestNewConfig(t *testing.T) {
	tests := map[string]struct {
		configPath string
		want       *Config
		wantErr    bool
	}{
		"successful run": {
			configPath: "test_vars/valid_vars.env",
			want: &Config{
				HTTP:  8081,
				Name:  "scootin_aboot",
				Users: "8212d8ba-74d1-49af-8a84-6d6c392ec71c,897737a8-77f1-4f53-8a51-6f9edaee6ed9",
			},
			wantErr: false,
		},
		"failed run": {
			configPath: "test_vars/invalid_vars.env",
			want:       nil,
			wantErr:    true,
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := NewConfig(context.Background(), tt.configPath)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewConfig() got = %v, want %v", got, tt.want)
			}
		})
	}
}
