package rustmaps

import (
	"testing"

	"github.com/maintc/rustmaps-cli/pkg/api"
	"github.com/maintc/rustmaps-cli/pkg/types"
)

func TestGenerator_LoadConfig(t *testing.T) {
	tests := []struct {
		name      string
		generator *Generator
		wantErr   bool
	}{
		{
			name: "Test LoadConfig",
			generator: NewMockedGenerator(t, &Generator{
				config: types.Config{},
				maps:   []*types.Map{},
				rmcli:  &api.RustMapsClient{},
			}),
			wantErr: false,
		},
		{
			name: "Test LoadConfig fail",
			generator: NewMockedGenerator(t, &Generator{
				configPath: "/tmp/asdfjasdlf/sadf432323/./.23",
				config:     types.Config{},
				maps:       []*types.Map{},
				rmcli:      &api.RustMapsClient{},
			}),
			wantErr: true,
		},
		{
			name: "Test LoadConfig fail tmp",
			generator: NewMockedGenerator(t, &Generator{
				configPath: "/tmp/",
				config:     types.Config{},
				maps:       []*types.Map{},
				rmcli:      &api.RustMapsClient{},
			}),
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.generator.LoadConfig(); (err != nil) != tt.wantErr {
				t.Errorf("Generator.LoadConfig() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGenerator_SaveConfig(t *testing.T) {
	tests := []struct {
		name      string
		generator *Generator
		wantErr   bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.generator.SaveConfig(); (err != nil) != tt.wantErr {
				t.Errorf("Generator.SaveConfig() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
