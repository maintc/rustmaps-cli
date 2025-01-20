package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"go.uber.org/zap"
)

func TestRustMapsClient_GetLimits(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			if strings.HasPrefix(r.URL.Path, "/maps/limits") {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				response := &RustMapsLimitsResponse{
					Meta: RustMapsLimitsResponseMeta{
						Status:     "complete",
						StatusCode: 200,
						Errors:     []string{},
					},
				}
				json.NewEncoder(w).Encode(response)
				return
			}
		}
		http.NotFound(w, r)
	}))
	defer mockServer.Close()
	type fields struct {
		apiURL      string
		apiKey      string
		rateLimiter *RateLimiter
	}
	type args struct {
		log *zap.Logger
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *RustMapsLimitsResponse
		wantErr bool
	}{
		{
			name: "GetLimits 200",
			fields: fields{
				apiURL:      mockServer.URL,
				apiKey:      "test",
				rateLimiter: &RateLimiter{},
			},
			args: args{
				log: zap.NewNop(),
			},
			want: &RustMapsLimitsResponse{
				Meta: RustMapsLimitsResponseMeta{
					Status:     "complete",
					StatusCode: 200,
					Errors:     []string{},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &RustMapsClient{
				ApiUrl:      tt.fields.apiURL,
				apiKey:      tt.fields.apiKey,
				rateLimiter: tt.fields.rateLimiter,
			}
			got, err := c.GetLimits(tt.args.log)
			if (err != nil) != tt.wantErr {
				t.Errorf("RustMapsClient.GetLimits() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != nil && got.Meta.StatusCode != tt.want.Meta.StatusCode {
				t.Errorf("RustMapsClient.GetLimits() = %v, want %v", got, tt.want)
			}
		})
	}
}
