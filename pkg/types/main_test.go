package types

import (
	"testing"
)

func TestNewMap(t *testing.T) {
	type args struct {
		seed        string
		size        int
		savedConfig string
		staging     bool
	}
	tests := []struct {
		name string
		args args
		want *Map
	}{
		{
			name: "TestNewMap",
			args: args{
				seed:        "test",
				size:        1,
				savedConfig: "test",
				staging:     true,
			},
			want: &Map{
				Seed:        "test",
				Size:        1,
				SavedConfig: "test",
				Staging:     true,
				Status:      "pending",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewMap(tt.args.seed, tt.args.size, tt.args.savedConfig, tt.args.staging); got.String() == tt.want.String() {
				t.Errorf("NewMap() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMap_SetFilename(t *testing.T) {
	type fields struct {
		Seed        string
		Size        int
		SavedConfig string
		Staging     bool
		MapID       string
		Status      string
		LastSync    string
		Filename    string
	}
	tests := []struct {
		name   string
		fields fields
	}{
		{
			name: "TestSetFilename",
			fields: fields{
				Seed:        "test",
				Size:        1,
				SavedConfig: "test",
				Staging:     true,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &Map{
				Seed:        tt.fields.Seed,
				Size:        tt.fields.Size,
				SavedConfig: tt.fields.SavedConfig,
				Staging:     tt.fields.Staging,
				MapID:       tt.fields.MapID,
				Status:      tt.fields.Status,
				LastSync:    tt.fields.LastSync,
				Filename:    tt.fields.Filename,
			}
			m.SetFilename()
		})
	}
}

func TestMap_ReportStatus(t *testing.T) {
	type fields struct {
		Seed        string
		Size        int
		SavedConfig string
		Staging     bool
		MapID       string
		Status      string
		LastSync    string
		Filename    string
	}
	type args struct {
		status string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "TestReportStatus",
			fields: fields{
				Seed:        "test",
				Size:        1,
				SavedConfig: "test",
				Staging:     true,
			},
			args: args{
				status: "test",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &Map{
				Seed:        tt.fields.Seed,
				Size:        tt.fields.Size,
				SavedConfig: tt.fields.SavedConfig,
				Staging:     tt.fields.Staging,
				MapID:       tt.fields.MapID,
				Status:      tt.fields.Status,
				LastSync:    tt.fields.LastSync,
				Filename:    tt.fields.Filename,
			}
			m.ReportStatus(tt.args.status)
		})
	}
}

func TestMap_String(t *testing.T) {
	type fields struct {
		Seed        string
		Size        int
		SavedConfig string
		Staging     bool
		MapID       string
		Status      string
		LastSync    string
		Filename    string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &Map{
				Seed:        tt.fields.Seed,
				Size:        tt.fields.Size,
				SavedConfig: tt.fields.SavedConfig,
				Staging:     tt.fields.Staging,
				MapID:       tt.fields.MapID,
				Status:      tt.fields.Status,
				LastSync:    tt.fields.LastSync,
				Filename:    tt.fields.Filename,
			}
			if got := m.String(); got != tt.want {
				t.Errorf("Map.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMap_ShouldSync(t *testing.T) {
	type fields struct {
		Seed        string
		Size        int
		SavedConfig string
		Staging     bool
		MapID       string
		Status      string
		LastSync    string
		Filename    string
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{
			name: "TestShouldSync",
			fields: fields{
				Seed:        "test",
				Size:        1,
				SavedConfig: "test",
				Staging:     true,
				LastSync:    "2021-01-01T00:00:00Z",
			},
			want: true,
		},
		{
			name: "TestShouldSync Fail to parse time",
			fields: fields{
				Seed:        "test",
				Size:        1,
				SavedConfig: "test",
				Staging:     true,
				LastSync:    "test",
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &Map{
				Seed:        tt.fields.Seed,
				Size:        tt.fields.Size,
				SavedConfig: tt.fields.SavedConfig,
				Staging:     tt.fields.Staging,
				MapID:       tt.fields.MapID,
				Status:      tt.fields.Status,
				LastSync:    tt.fields.LastSync,
				Filename:    tt.fields.Filename,
			}
			if got := m.ShouldSync(); got != tt.want {
				t.Errorf("Map.ShouldSync() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMap_MarkSynced(t *testing.T) {
	type fields struct {
		Seed        string
		Size        int
		SavedConfig string
		Staging     bool
		MapID       string
		Status      string
		LastSync    string
		Filename    string
	}
	tests := []struct {
		name   string
		fields fields
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &Map{
				Seed:        tt.fields.Seed,
				Size:        tt.fields.Size,
				SavedConfig: tt.fields.SavedConfig,
				Staging:     tt.fields.Staging,
				MapID:       tt.fields.MapID,
				Status:      tt.fields.Status,
				LastSync:    tt.fields.LastSync,
				Filename:    tt.fields.Filename,
			}
			m.MarkSynced()
		})
	}
}

func TestMap_MergeFrom(t *testing.T) {
	type fields struct {
		Seed        string
		Size        int
		SavedConfig string
		Staging     bool
		MapID       string
		Status      string
		LastSync    string
		Filename    string
	}
	type args struct {
		other Map
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "TestMergeFrom",
			fields: fields{
				Seed:        "test",
				Size:        1,
				SavedConfig: "test",
				Staging:     true,
			},
			args: args{
				other: Map{
					Seed:        "test",
					Size:        1,
					SavedConfig: "test",
					Staging:     true,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &Map{
				Seed:        tt.fields.Seed,
				Size:        tt.fields.Size,
				SavedConfig: tt.fields.SavedConfig,
				Staging:     tt.fields.Staging,
				MapID:       tt.fields.MapID,
				Status:      tt.fields.Status,
				LastSync:    tt.fields.LastSync,
				Filename:    tt.fields.Filename,
			}
			m.MergeFrom(tt.args.other)
		})
	}
}

func TestMap_SaveJSON(t *testing.T) {
	type fields struct {
		Seed        string
		Size        int
		SavedConfig string
		Staging     bool
		MapID       string
		Status      string
		LastSync    string
		Filename    string
	}
	type args struct {
		outputDir string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "TestSaveJSON",
			fields: fields{
				Seed:        "test",
				Size:        1,
				SavedConfig: "test",
				Staging:     true,
				Filename:    "test.json",
			},
			args: args{
				outputDir: t.TempDir(),
			},
			wantErr: false,
		},
		{
			name: "TestSaveJSON Empty filename",
			fields: fields{
				Seed:        "test",
				Size:        1,
				SavedConfig: "test",
				Staging:     true,
				Filename:    "",
			},
			args: args{
				outputDir: t.TempDir(),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &Map{
				Seed:        tt.fields.Seed,
				Size:        tt.fields.Size,
				SavedConfig: tt.fields.SavedConfig,
				Staging:     tt.fields.Staging,
				MapID:       tt.fields.MapID,
				Status:      tt.fields.Status,
				LastSync:    tt.fields.LastSync,
				Filename:    tt.fields.Filename,
			}
			if err := m.SaveJSON(tt.args.outputDir); (err != nil) != tt.wantErr {
				t.Errorf("Map.SaveJSON() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
