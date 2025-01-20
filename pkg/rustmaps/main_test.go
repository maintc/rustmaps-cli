// Package rustmaps provides functionality for generating and managing Rust game maps
package rustmaps

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/maintc/rustmaps-cli/pkg/api"
	"github.com/maintc/rustmaps-cli/pkg/types"
	"go.uber.org/zap"
)

func NewMockedGenerator(t *testing.T, other *Generator) *Generator {

	if other == nil {
		fmt.Fprintf(os.Stderr, "NewMockedGenerator: other is nil\n")
		panic("NewMockedGenerator: other is nil")
	}

	mocked := &Generator{
		configPath:   filepath.Join(t.TempDir(), "config.json"),
		importsDir:   t.TempDir(),
		downloadsDir: t.TempDir(),
		logPath:      filepath.Join(t.TempDir(), "generator.log"),
		backoffTime:  0,
	}

	if other.configPath != "" {
		mocked.configPath = other.configPath
	}

	if other.importsDir != "" {
		mocked.importsDir = other.importsDir
	}

	if other.downloadsDir != "" {
		mocked.downloadsDir = other.downloadsDir
	}

	if other.logPath != "" {
		mocked.logPath = other.logPath
	}

	mocked.config = other.config
	mocked.maps = other.maps
	mocked.rmcli = other.rmcli
	mocked.target = other.target
	mocked.baseDir = other.baseDir

	return mocked
}

type MockedRustMapsCLI struct {
	api.RustMapsClientBase
	ApiUrl            string
	apiKey            string
	rateLimiter       *api.RateLimiter
	MockedServerUrl   string
	ConcurrentCurrent int
	ConcurrentAllowed int
	MonthlyCurrent    int
	MonthlyAllowed    int
	LimitsError       bool
}

func (c *MockedRustMapsCLI) GetStatus(log *zap.Logger, m *types.Map) (*api.RustMapsStatusResponse, error) {
	canDownload := true
	switch m.Seed {
	case "0":
		return nil, fmt.Errorf("error")
	case "2":
		canDownload = false
	}

	if c.MockedServerUrl == "" {
		c.MockedServerUrl = "http://localhost"
	}

	return &api.RustMapsStatusResponse{
		Meta: api.RustMapsStatusResponseMeta{
			Status:     m.Status,
			StatusCode: 200,
			Errors:     []string{},
		},
		Data: api.RustMapsStatusResponseData{
			Seed:         parseInt(m.Seed),
			Size:         m.Size,
			CanDownload:  canDownload,
			DownloadURL:  fmt.Sprintf("%s/%s_%d.map", c.MockedServerUrl, m.Seed, m.Size),
			ImageURL:     fmt.Sprintf("%s/%s_%d.png", c.MockedServerUrl, m.Seed, m.Size),
			RawImageURL:  fmt.Sprintf("%s/%s_%d_raw.png", c.MockedServerUrl, m.Seed, m.Size),
			ThumbnailURL: fmt.Sprintf("%s/%s_%d_thumbnail.png", c.MockedServerUrl, m.Seed, m.Size),
			ImageIconURL: fmt.Sprintf("%s/%s_%d_icons.png", c.MockedServerUrl, m.Seed, m.Size),
			URL:          fmt.Sprintf("%s/%s_%d", c.MockedServerUrl, m.Seed, m.Size),
		},
	}, nil
}

func (c *MockedRustMapsCLI) GetLimits(log *zap.Logger) (*api.RustMapsLimitsResponse, error) {
	if c.LimitsError {
		return nil, fmt.Errorf("error")
	}
	return &api.RustMapsLimitsResponse{
		Meta: api.RustMapsLimitsResponseMeta{
			Status:     "complete",
			StatusCode: 200,
			Errors:     []string{},
		},
		Data: api.RustMapsLimitsResponseData{
			Concurrent: api.RustMapsLimitsResponseDataConcurrent{
				Current: c.ConcurrentCurrent,
				Allowed: c.ConcurrentAllowed,
			},
			Monthly: api.RustMapsLimitsResponseDataMonthly{
				Current: c.MonthlyCurrent,
				Allowed: c.MonthlyAllowed,
			},
		},
	}, nil
}

func (c *MockedRustMapsCLI) GenerateProcedural(log *zap.Logger, m *types.Map) (*api.RustMapsGenerateResponse, error) {
	return &api.RustMapsGenerateResponse{
		Meta: api.RustMapsGenerateResponseMeta{
			Status:     "complete",
			StatusCode: 200,
			Errors:     []string{},
		},
		Data: api.RustMapsGenerateResponseData{
			MapID: "123",
		},
	}, nil
}

func TestNewGenerator(t *testing.T) {
	tmpDir := t.TempDir()
	tests := []struct {
		name    string
		want    *Generator
		baseDir *string
		wantErr bool
	}{
		{
			name:    "New generator",
			want:    &Generator{},
			baseDir: &tmpDir,
			wantErr: false,
		},
		{
			name:    "Use home dir if baseDir is nil",
			want:    &Generator{},
			baseDir: nil,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewGenerator(tt.baseDir)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewGenerator() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got.config.APIKey != tt.want.config.APIKey {
				t.Errorf("NewGenerator() API Key = %v, want %v", got.config.APIKey, tt.want.config.APIKey)
			}
		})
	}
}

func TestGenerator_InitDirs(t *testing.T) {
	tmpDir := t.TempDir()
	tests := []struct {
		name      string
		generator *Generator
		wantErr   bool
	}{
		{
			name: "Init dirs",
			generator: NewMockedGenerator(t, &Generator{
				baseDir: tmpDir,
			}),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.generator.InitDirs(); (err != nil) != tt.wantErr {
				t.Errorf("Generator.initDirs() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGenerator_CanGenerate(t *testing.T) {
	type args struct {
		log *zap.Logger
	}
	tests := []struct {
		name      string
		generator *Generator
		args      args
		want      bool
	}{
		{
			name: "Can generate",
			generator: NewMockedGenerator(t, &Generator{
				rmcli: &MockedRustMapsCLI{
					ConcurrentCurrent: 0,
					ConcurrentAllowed: 3,
					MonthlyCurrent:    0,
					MonthlyAllowed:    250,
				},
			}),
			args: args{
				log: zap.NewNop(),
			},
			want: true,
		},
		{
			name: "Can generate GetLimits error",
			generator: NewMockedGenerator(t, &Generator{
				rmcli: &MockedRustMapsCLI{
					ConcurrentCurrent: 0,
					ConcurrentAllowed: 3,
					MonthlyCurrent:    0,
					MonthlyAllowed:    250,
					LimitsError:       true,
				},
			}),
			args: args{
				log: zap.NewNop(),
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.generator.CanGenerate(tt.args.log); got != tt.want {
				t.Errorf("Generator.CanGenerate() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGenerator_GetStatus(t *testing.T) {
	type args struct {
		log *zap.Logger
		m   *types.Map
	}
	tests := []struct {
		name      string
		generator *Generator
		args      args
		want      *api.RustMapsStatusResponse
		wantErr   bool
	}{
		{
			name: "Get status",
			generator: NewMockedGenerator(t, &Generator{
				rmcli: &MockedRustMapsCLI{},
			}),
			args: args{
				m: &types.Map{
					Seed: "1",
					Size: 4000,
				},
			},
		},
		{
			name: "Get status error",
			generator: NewMockedGenerator(t, &Generator{
				rmcli: &MockedRustMapsCLI{},
			}),
			args: args{
				m: &types.Map{
					Seed: "0",
					Size: 4000,
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.generator.GetStatus(tt.args.log, tt.args.m)
			if (err != nil) != tt.wantErr {
				t.Errorf("Generator.GetStatus() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != nil && got.Data.Seed != parseInt(tt.args.m.Seed) {
				t.Errorf("Generator.GetStatus() = %v, want %v", got.Data.Seed, tt.args.m.Seed)
			}
		})
	}
}

func TestGenerator_SyncStatus(t *testing.T) {
	type args struct {
		log *zap.Logger
		m   *types.Map
	}
	tests := []struct {
		name      string
		generator *Generator
		args      args
		wantErr   bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.generator.SyncStatus(tt.args.log, tt.args.m); (err != nil) != tt.wantErr {
				t.Errorf("Generator.SyncStatus() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGenerator_AddMap(t *testing.T) {
	type args struct {
		m *types.Map
	}
	tests := []struct {
		name      string
		generator *Generator
		args      args
	}{
		{
			name:      "Add map",
			generator: NewMockedGenerator(t, &Generator{}),
			args: args{
				m: &types.Map{
					Seed:        "123",
					Size:        4000,
					SavedConfig: "default",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.generator.AddMap(tt.args.m)
		})
	}
}
