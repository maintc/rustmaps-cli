package rustmaps

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/maintc/rustmaps-cli/pkg/common"
	"github.com/maintc/rustmaps-cli/pkg/types"

	"go.uber.org/zap"
)

func TestGenerator_Generate(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// send generic json 200
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"status": "success"}`))
	}))
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
			name: "Test Generate",
			generator: NewMockedGenerator(t, &Generator{
				config: types.Config{
					APIKey: "test",
					Tier:   "Free",
				},
				maps: []*types.Map{
					{
						Status:      common.StatusPending,
						SavedConfig: "",
						Seed:        "test",
						Size:        4000,
						Staging:     false,
					},
					{
						Status:      common.StatusComplete,
						SavedConfig: "",
						Seed:        "test",
						Size:        4000,
						Staging:     false,
					},
					{
						Status:      common.StatusGenerating,
						SavedConfig: "",
						Seed:        "test",
						Size:        4000,
						Staging:     false,
					},
				},
				target: "",
				rmcli:  &MockedRustMapsCLI{},
			}),
			args: args{
				log: zap.NewNop(),
			},
			want: true,
		},
		{
			name: "Test Generate",
			generator: NewMockedGenerator(t, &Generator{
				config: types.Config{
					APIKey: "test",
					Tier:   "Free",
				},
				maps: []*types.Map{
					{
						Status:      common.StatusComplete,
						SavedConfig: "",
						Seed:        "test",
						Size:        4000,
						Staging:     false,
					},
				},
				target: "",
				rmcli:  &MockedRustMapsCLI{},
			}),
			args: args{
				log: zap.NewNop(),
			},
			want: false,
		},
		{
			name: "Test Generate fail auth",
			generator: NewMockedGenerator(t, &Generator{
				config: types.Config{},
				maps: []*types.Map{
					{
						Status:      common.StatusComplete,
						SavedConfig: "",
						Seed:        "test",
						Size:        4000,
						Staging:     false,
					},
				},
				target: "",
				rmcli:  &MockedRustMapsCLI{},
			}),
			args: args{
				log: zap.NewNop(),
			},
			want: false,
		},
		{
			name: "Test Generate genrate procedural",
			generator: NewMockedGenerator(t, &Generator{
				config: types.Config{
					APIKey: "test",
					Tier:   "Premium",
				},
				maps: []*types.Map{
					{
						Status:      common.StatusPending,
						SavedConfig: "",
						Seed:        "test",
						Size:        4000,
						Staging:     false,
					},
				},
				target: "",
				rmcli: &MockedRustMapsCLI{
					MockedServerUrl:   mockServer.URL,
					ConcurrentAllowed: 8,
					MonthlyAllowed:    800,
				},
			}),
			args: args{
				log: zap.NewNop(),
			},
			want: true,
		},
		// {
		// 	name: "Test Generate genrate custom",
		// 	generator: NewMockedGenerator(t, &Generator{
		// 		config: types.Config{
		// 			APIKey: "test",
		// 			Tier:   "Premium",
		// 		},
		// 		maps: []*types.Map{
		// 			{
		// 				Status:      common.StatusPending,
		// 				SavedConfig: "default",
		// 				Seed:        "test",
		// 				Size:        4000,
		// 				Staging:     false,
		// 			},
		// 		},
		// 		target: "",
		// 		rmcli: &MockedRustMapsCLI{
		// 			MockedServerUrl:   mockServer.URL,
		// 			ConcurrentAllowed: 8,
		// 			MonthlyAllowed:    800,
		// 		},
		// 	}),
		// 	args: args{
		// 		log: zap.NewNop(),
		// 	},
		// 	want: true,
		// },
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.generator.Generate(tt.args.log); got != tt.want {
				t.Errorf("Generator.Generate() = %v, want %v", got, tt.want)
			}
		})
	}
}
