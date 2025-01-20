package rustmaps

import (
	"path/filepath"
	"testing"

	"github.com/maintc/rustmaps-cli/pkg/types"
	"go.uber.org/zap"
)

func TestGenerator_ValidateAuthentication(t *testing.T) {
	type args struct {
		log *zap.Logger
	}
	tests := []struct {
		name      string
		generator *Generator
		args      args
		wantErr   bool
	}{
		{
			name:      "API key not set should fail to authenticate",
			generator: NewMockedGenerator(t, &Generator{}),
			args: args{
				log: zap.NewNop(),
			},
			wantErr: true,
		},
		{
			name: "Tier not set should fail to authenticate",
			generator: NewMockedGenerator(t, &Generator{
				config: types.Config{
					APIKey: "test",
				},
			}),
			args: args{
				log: zap.NewNop(),
			},
			wantErr: true,
		},
		{
			name: "API key and tier set should authenticate",
			generator: NewMockedGenerator(t, &Generator{
				config: types.Config{
					APIKey: "test",
					Tier:   "test",
				},
			}),
			args: args{
				log: zap.NewNop(),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.generator.ValidateAuthentication(tt.args.log); (err != nil) != tt.wantErr {
				t.Errorf("Generator.ValidateAuthentication() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGenerator_DetermineTier(t *testing.T) {
	type args struct {
		log *zap.Logger
	}
	tests := []struct {
		name      string
		generator *Generator
		args      args
		want      string
		want1     bool
	}{
		{
			name: "Test DetermineTier",
			generator: NewMockedGenerator(t, &Generator{
				config: types.Config{
					APIKey: "test",
					Tier:   "Free",
				},
				maps:   nil,
				target: "",
				rmcli: &MockedRustMapsCLI{
					ConcurrentAllowed: 3,
					MonthlyAllowed:    250,
				},
				configPath:   filepath.Join(t.TempDir(), "config.json"),
				importsDir:   t.TempDir(),
				downloadsDir: t.TempDir(),
				logPath:      filepath.Join(t.TempDir(), "generator.log"),
			}),
			args: args{
				log: zap.NewNop(),
			},
			want:  "Free",
			want1: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := tt.generator.DetermineTier(tt.args.log)
			if got != tt.want {
				t.Errorf("Generator.DetermineTier() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("Generator.DetermineTier() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
