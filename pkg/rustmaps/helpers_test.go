package rustmaps

import (
	"testing"

	"github.com/maintc/rustmaps-cli/pkg/common"
	"github.com/maintc/rustmaps-cli/pkg/types"
)

func TestGenerator_IsApiKeySet(t *testing.T) {
	tests := []struct {
		name      string
		generator *Generator
		want      bool
	}{
		{
			name: "API key is set",
			generator: NewMockedGenerator(t, &Generator{
				config: types.Config{
					APIKey: "test",
				},
			}),
			want: true,
		},
		{
			name: "API key is not set",
			generator: NewMockedGenerator(t, &Generator{
				config: types.Config{},
			}),
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.generator.IsApiKeySet(); got != tt.want {
				t.Errorf("Generator.IsApiKeySet() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGenerator_Pending(t *testing.T) {
	tests := []struct {
		name      string
		generator *Generator
		want      bool
	}{
		{
			name: "Get pending maps",
			generator: NewMockedGenerator(t, &Generator{
				maps: []*types.Map{
					{
						Status: common.StatusPending,
					},
				},
			}),
			want: true,
		},
		{
			name: "Get no pending maps",
			generator: NewMockedGenerator(t, &Generator{
				maps: []*types.Map{
					{
						Status: common.StatusComplete,
					},
				},
			}),
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.generator.Pending(); got != tt.want {
				t.Errorf("Generator.Pending() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGenerator_Generating(t *testing.T) {
	tests := []struct {
		name      string
		generator *Generator
		want      bool
	}{
		{
			name: "Get generating maps",
			generator: NewMockedGenerator(t, &Generator{
				maps: []*types.Map{
					{
						Status: common.StatusGenerating,
					},
				},
			}),
			want: true,
		},
		{
			name: "Get no generating maps",
			generator: NewMockedGenerator(t, &Generator{
				maps: []*types.Map{
					{
						Status: common.StatusComplete,
					},
				},
			}),
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.generator.Generating(); got != tt.want {
				t.Errorf("Generator.Generating() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGenerator_ContainCustomMaps(t *testing.T) {
	tests := []struct {
		name      string
		generator *Generator
		want      bool
	}{
		{
			name: "Contain custom maps",
			generator: NewMockedGenerator(t, &Generator{
				maps: []*types.Map{
					{
						SavedConfig: "default",
					},
				},
			}),
			want: true,
		},
		{
			name: "Contain no custom maps",
			generator: NewMockedGenerator(t, &Generator{
				maps: []*types.Map{
					{
						SavedConfig: "",
					},
				},
			}),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.generator.ContainCustomMaps(); got != tt.want {
				t.Errorf("Generator.ContainCustomMaps() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGenerator_GetRandomSeed(t *testing.T) {
	tests := []struct {
		name      string
		generator *Generator
	}{
		{
			name:      "Get random seed",
			generator: NewMockedGenerator(t, &Generator{}),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.generator.GetRandomSeed(); parseInt(got) == 0 {
				t.Errorf("Generator.GetRandomSeed() = %v, want %v", got, "a number")
			}
		})
	}
}

func Test_parseInt(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "Parse 1",
			args: args{
				s: "1",
			},
			want: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := parseInt(tt.args.s); got != tt.want {
				t.Errorf("parseInt() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_parseBool(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "Parse true",
			args: args{
				s: "true",
			},
			want: true,
		},
		{
			name: "Parse false",
			args: args{
				s: "false",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := parseBool(tt.args.s); got != tt.want {
				t.Errorf("parseBool() = %v, want %v", got, tt.want)
			}
		})
	}
}
