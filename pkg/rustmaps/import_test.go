package rustmaps

import (
	"testing"

	"github.com/maintc/rustmaps-cli/pkg/types"
	"go.uber.org/zap"
)

func TestGenerator_LoadCSV(t *testing.T) {
	type args struct {
		log      *zap.Logger
		mapsPath string
	}
	tests := []struct {
		name      string
		generator *Generator
		args      args
		wantErr   bool
	}{
		{
			name: "Load CSV",
			generator: NewMockedGenerator(t, &Generator{
				config: types.Config{
					APIKey: "test",
					Tier:   "test",
				},
			}),
			args: args{
				log:      zap.NewNop(),
				mapsPath: "../../tests/files/test_valid.csv",
			},
			wantErr: false,
		},
		{
			name: "Fail if custom map on free tier",
			generator: NewMockedGenerator(t, &Generator{
				config: types.Config{
					APIKey: "test",
					Tier:   "Free",
				},
			}),
			args: args{
				log:      zap.NewNop(),
				mapsPath: "../../tests/files/test_valid.csv",
			},
			wantErr: true,
		},
		{
			name:      "Fail if not authed",
			generator: NewMockedGenerator(t, &Generator{}),
			args: args{
				log:      zap.NewNop(),
				mapsPath: "../../tests/files/test_valid.csv",
			},
			wantErr: true,
		},
		{
			name: "Fail if invalid csv",
			generator: NewMockedGenerator(t, &Generator{
				config: types.Config{
					APIKey: "test",
					Tier:   "test",
				},
			}),
			args: args{
				log:      zap.NewNop(),
				mapsPath: "../../tests/files/test_invalid_headers.csv",
			},
			wantErr: true,
		},
		{
			name: "Fail if csv empty",
			generator: NewMockedGenerator(t, &Generator{
				config: types.Config{
					APIKey: "test",
					Tier:   "test",
				},
			}),
			args: args{
				log:      zap.NewNop(),
				mapsPath: "../../tests/files/test_empty.csv",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.generator.LoadCSV(tt.args.log, tt.args.mapsPath); (err != nil) != tt.wantErr {
				t.Errorf("Generator.LoadCSV() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGenerator_ValidateCSV(t *testing.T) {
	type args struct {
		log      *zap.Logger
		mapsPath string
	}
	tests := []struct {
		name      string
		generator *Generator
		args      args
		wantErr   bool
	}{
		{
			name:      "Validate CSV",
			generator: NewMockedGenerator(t, &Generator{}),
			args: args{
				log:      zap.NewNop(),
				mapsPath: "../../tests/files/test_valid.csv",
			},
			wantErr: false,
		},
		{
			name:      "Invalidate not enough headers",
			generator: NewMockedGenerator(t, &Generator{}),
			args: args{
				log:      zap.NewNop(),
				mapsPath: "../../tests/files/test_not_enough_headers.csv",
			},
			wantErr: true,
		},
		{
			name:      "Invalidate invalid headers",
			generator: NewMockedGenerator(t, &Generator{}),
			args: args{
				log:      zap.NewNop(),
				mapsPath: "../../tests/files/test_invalid_headers.csv",
			},
			wantErr: true,
		},
		{
			name:      "Invalidate does not exist",
			generator: NewMockedGenerator(t, &Generator{}),
			args: args{
				log:      zap.NewNop(),
				mapsPath: "../../tests/files/test_dne.csv",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.generator.ValidateCSV(tt.args.log, tt.args.mapsPath); (err != nil) != tt.wantErr {
				t.Errorf("Generator.ValidateCSV() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGenerator_Import(t *testing.T) {
	type args struct {
		log   *zap.Logger
		force bool
	}
	// create tmp dir
	tests := []struct {
		name      string
		generator *Generator
		args      args
		wantErr   bool
	}{
		{
			name: "Import a map",
			generator: NewMockedGenerator(t, &Generator{
				maps: []*types.Map{
					{
						Seed:        "123",
						Size:        4000,
						SavedConfig: "default",
						Staging:     true,
					},
				},
			}),
			args: args{
				log:   zap.NewNop(),
				force: false,
			},
			wantErr: false,
		},
		{
			name: "Import the same map",
			generator: NewMockedGenerator(t, &Generator{
				maps: []*types.Map{
					{
						Seed:        "123",
						Size:        4000,
						SavedConfig: "default",
						Staging:     true,
					},
					{
						Seed:        "123",
						Size:        4000,
						SavedConfig: "default",
						Staging:     true,
					},
				},
			}),
			args: args{
				log:   zap.NewNop(),
				force: false,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.generator.Import(tt.args.log, tt.args.force); (err != nil) != tt.wantErr {
				t.Errorf("Generator.Import() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
