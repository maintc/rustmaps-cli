package rustmaps

import (
	"testing"

	"github.com/maintc/rustmaps-cli/pkg/api"
)

func TestGenerator_SetApiKey(t *testing.T) {
	type args struct {
		apiKey string
	}
	tests := []struct {
		name      string
		generator *Generator
		args      args
	}{
		{
			name: "Set API key",
			generator: NewMockedGenerator(t, &Generator{
				rmcli: &api.RustMapsClient{},
			}),
			args: args{
				apiKey: "test",
			},
		},
		{
			name: "Set empty API key",
			generator: NewMockedGenerator(t, &Generator{
				rmcli: &api.RustMapsClient{},
			}),
			args: args{
				apiKey: "",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.generator.SetApiKey(tt.args.apiKey)
		})
	}
}

func TestGenerator_SetTier(t *testing.T) {
	type args struct {
		tier string
	}
	tests := []struct {
		name      string
		generator *Generator
		args      args
	}{
		{
			name:      "Set tier",
			generator: NewMockedGenerator(t, &Generator{}),
			args: args{
				tier: "test",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.generator.SetTier(tt.args.tier)
		})
	}
}
