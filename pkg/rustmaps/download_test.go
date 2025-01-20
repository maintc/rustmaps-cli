package rustmaps

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"

	"github.com/maintc/rustmaps-cli/pkg/common"
	"github.com/maintc/rustmaps-cli/pkg/types"
	"go.uber.org/zap"
)

func TestGenerator_OverrideDownloadsDir(t *testing.T) {
	type args struct {
		log *zap.Logger
		dir string
	}
	tests := []struct {
		name      string
		generator *Generator
		args      args
	}{
		{
			name: "Test OverrideDownloadsDir",
			generator: NewMockedGenerator(t, &Generator{
				config:       types.Config{},
				maps:         nil,
				target:       "",
				rmcli:        nil,
				configPath:   filepath.Join(t.TempDir(), "config.json"),
				importsDir:   t.TempDir(),
				downloadsDir: t.TempDir(),
				logPath:      filepath.Join(t.TempDir(), "generator.log"),
			}),
			args: args{
				log: zap.NewNop(),
				dir: "test",
			},
		},
		{
			name: "Test OverrideDownloadsDir invalid dir",
			generator: NewMockedGenerator(t, &Generator{
				config:       types.Config{},
				maps:         nil,
				target:       "",
				rmcli:        nil,
				configPath:   filepath.Join(t.TempDir(), "config.json"),
				importsDir:   t.TempDir(),
				downloadsDir: t.TempDir(),
				logPath:      filepath.Join(t.TempDir(), "generator.log"),
			}),
			args: args{
				log: zap.NewNop(),
				dir: "",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.generator.OverrideDownloadsDir(tt.args.log, tt.args.dir)
		})
	}
}

func TestGenerator_DownloadFile(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/valid" {
			// send file content for download
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			response := "test"
			w.Write([]byte(response))
		} else if r.URL.Path == "/error" {
			// send error response
			w.WriteHeader(http.StatusInternalServerError)
		}
	}))
	type args struct {
		log    *zap.Logger
		url    string
		target string
	}
	tests := []struct {
		name      string
		generator *Generator
		args      args
		wantErr   bool
	}{
		{
			name: "Test DownloadFile",
			generator: NewMockedGenerator(t, &Generator{
				config: types.Config{},
				maps: []*types.Map{
					{
						Seed:        "1",
						Size:        4000,
						SavedConfig: "default",
						Status:      common.StatusComplete,
					},
				},
				rmcli: &MockedRustMapsCLI{ApiUrl: mockServer.URL, MockedServerUrl: mockServer.URL},
			}),
			args: args{
				log:    zap.NewNop(),
				url:    fmt.Sprintf("%s/%s", mockServer.URL, "valid"),
				target: filepath.Join(t.TempDir(), "test.json"),
			},
		},
		{
			name: "Test DownloadFile error",
			generator: NewMockedGenerator(t, &Generator{
				config: types.Config{},
				maps: []*types.Map{
					{
						Seed:        "1",
						Size:        4000,
						SavedConfig: "default",
						Status:      common.StatusComplete,
					},
				},
				rmcli: &MockedRustMapsCLI{ApiUrl: mockServer.URL, MockedServerUrl: mockServer.URL},
			}),
			args: args{
				log:    zap.NewNop(),
				url:    fmt.Sprintf("%s/%s", mockServer.URL, "error"),
				target: filepath.Join(t.TempDir(), "test.json"),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.generator.DownloadFile(tt.args.log, tt.args.url, tt.args.target); (err != nil) != tt.wantErr {
				t.Errorf("Generator.DownloadFile() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGenerator_Download(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// send file content for download
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		response := "test"
		w.Write([]byte(response))
	}))
	defer mockServer.Close()
	type args struct {
		log     *zap.Logger
		version string
	}
	tests := []struct {
		name      string
		generator *Generator
		args      args
		wantErr   bool
	}{
		{
			name: "Test Download",
			generator: NewMockedGenerator(t, &Generator{
				config: types.Config{},
				maps: []*types.Map{
					{
						Seed:        "1",
						Size:        4000,
						SavedConfig: "default",
						Status:      common.StatusComplete,
					},
					{
						Seed:        "1",
						Size:        4000,
						SavedConfig: "default",
						Status:      common.StatusGenerating,
					},
				},
				rmcli: &MockedRustMapsCLI{ApiUrl: mockServer.URL, MockedServerUrl: mockServer.URL},
			}),
			args:    args{log: zap.NewNop(), version: "test"},
			wantErr: false,
		},
		{
			name: "Test Download no maps",
			generator: NewMockedGenerator(t, &Generator{
				config: types.Config{},
				maps:   []*types.Map{},
				rmcli:  &MockedRustMapsCLI{ApiUrl: mockServer.URL, MockedServerUrl: mockServer.URL},
			}),
			args:    args{log: zap.NewNop(), version: "test"},
			wantErr: true,
		},
		{
			name: "Test Download GetStatus error",
			generator: NewMockedGenerator(t, &Generator{
				config: types.Config{},
				maps: []*types.Map{
					{
						Seed:        "0",
						Size:        4000,
						SavedConfig: "default",
						Status:      common.StatusComplete,
					},
				},
				rmcli: &MockedRustMapsCLI{ApiUrl: mockServer.URL, MockedServerUrl: mockServer.URL},
			}),
			args:    args{log: zap.NewNop(), version: "test"},
			wantErr: true,
		},
		{
			name: "Test Download Can't download",
			generator: NewMockedGenerator(t, &Generator{
				config: types.Config{},
				maps: []*types.Map{
					{
						Seed:        "2",
						Size:        4000,
						SavedConfig: "default",
						Status:      common.StatusComplete,
						Staging:     true,
					},
				},
				rmcli: &MockedRustMapsCLI{ApiUrl: mockServer.URL, MockedServerUrl: mockServer.URL},
			}),
			args:    args{log: zap.NewNop(), version: "test"},
			wantErr: false,
		},
		{
			name: "Test Download failed download",
			generator: NewMockedGenerator(t, &Generator{
				config: types.Config{},
				maps: []*types.Map{
					{
						Seed:        "1",
						Size:        4000,
						SavedConfig: "default",
						Status:      common.StatusComplete,
						Staging:     true,
					},
				},
				rmcli: &MockedRustMapsCLI{ApiUrl: mockServer.URL},
			}),
			args:    args{log: zap.NewNop(), version: "test"},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.generator.Download(tt.args.log, tt.args.version); (err != nil) != tt.wantErr {
				t.Errorf("Generator.Download() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
