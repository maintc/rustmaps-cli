package rustmaps

import (
	"reflect"
	"testing"

	"github.com/maintc/rustmaps-cli/pkg/types"
)

func TestGenerator_GetDownloadsDir(t *testing.T) {
	tests := []struct {
		name      string
		generator *Generator
		want      string
	}{
		{
			name:      "Get downloads directory",
			generator: NewMockedGenerator(t, &Generator{}),
			want:      "test",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.generator.GetDownloadsDir(); len(got) == 0 {
				t.Errorf("Generator.GetDownloadsDir() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGenerator_GetImportDir(t *testing.T) {
	tests := []struct {
		name      string
		generator *Generator
		want      string
	}{
		{
			name:      "Get imports directory",
			generator: NewMockedGenerator(t, &Generator{}),
			want:      "test",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.generator.GetImportDir(); len(got) == 0 {
				t.Errorf("Generator.GetImportDir() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGenerator_GetLogPath(t *testing.T) {
	tests := []struct {
		name      string
		generator *Generator
		want      string
	}{
		{
			name:      "Get log path",
			generator: NewMockedGenerator(t, &Generator{}),
			want:      "test",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.generator.GetLogPath(); len(got) == 0 {
				t.Errorf("Generator.GetLogPath() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGenerator_GetConfigPath(t *testing.T) {
	tests := []struct {
		name      string
		generator *Generator
		want      string
	}{
		{
			name:      "Get config path",
			generator: NewMockedGenerator(t, &Generator{}),
			want:      "test",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.generator.GetConfigPath(); len(got) == 0 {
				t.Errorf("Generator.GetConfigPath() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGenerator_GetMaps(t *testing.T) {
	tests := []struct {
		name      string
		generator *Generator
		want      []*types.Map
	}{
		{
			name: "Get maps",
			generator: NewMockedGenerator(t, &Generator{
				maps: []*types.Map{
					{
						Seed:        "test",
						Size:        1,
						SavedConfig: "test",
						Staging:     true,
						MapID:       "test",
						Status:      "test",
						LastSync:    "test",
						Filename:    "test",
					},
				},
			}),
			want: []*types.Map{
				{
					Seed:        "test",
					Size:        1,
					SavedConfig: "test",
					Staging:     true,
					MapID:       "test",
					Status:      "test",
					LastSync:    "test",
					Filename:    "test",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.generator.GetMaps(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Generator.GetMaps() = %v, want %v", got, tt.want)
			}
		})
	}
}
