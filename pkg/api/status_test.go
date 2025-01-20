package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/maintc/rustmaps-cli/pkg/common"
	"github.com/maintc/rustmaps-cli/pkg/types"
	"go.uber.org/zap"
)

func TestRustMapsClient_GetStatus(t *testing.T) {

	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			if strings.HasPrefix(r.URL.Path, "/maps/4000/1") {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				response := &RustMapsStatusResponse{
					Meta: RustMapsStatusResponseMeta{
						Status:     common.StatusComplete,
						StatusCode: 200,
						Errors:     []string{},
					},
				}
				json.NewEncoder(w).Encode(response)
				return
			} else if strings.HasPrefix(r.URL.Path, "/maps/4000/2") {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusUnauthorized)
				response := &RustMapsStatusResponse{
					Meta: RustMapsStatusResponseMeta{
						Status:     common.StatusUnauthorized,
						StatusCode: 401,
						Errors:     []string{},
					},
				}
				json.NewEncoder(w).Encode(response)
				return
			} else if strings.HasPrefix(r.URL.Path, "/maps/4000/3") {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusForbidden)
				response := &RustMapsStatusResponse{
					Meta: RustMapsStatusResponseMeta{
						Status:     common.StatusForbidden,
						StatusCode: 403,
						Errors:     []string{},
					},
				}
				json.NewEncoder(w).Encode(response)
				return
			} else if strings.HasPrefix(r.URL.Path, "/maps/4000/4") {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusNotFound)
				response := &RustMapsStatusResponse{
					Meta: RustMapsStatusResponseMeta{
						Status:     common.StatusNotFound,
						StatusCode: 404,
						Errors:     []string{},
					},
				}
				json.NewEncoder(w).Encode(response)
				return
			} else if strings.HasPrefix(r.URL.Path, "/maps/4000/5") {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusConflict)
				response := &RustMapsStatusResponse{
					Meta: RustMapsStatusResponseMeta{
						Status:     common.StatusGenerating,
						StatusCode: 409,
						Errors:     []string{},
					},
				}
				json.NewEncoder(w).Encode(response)
				return
			} else if strings.HasPrefix(r.URL.Path, "/maps/4000/6") {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusInternalServerError)
				response := &RustMapsStatusResponse{
					Meta: RustMapsStatusResponseMeta{
						Status:     "",
						StatusCode: 500,
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
		m   *types.Map
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *RustMapsStatusResponse
		wantErr bool
	}{
		{
			name: "Test GetStatus 200",
			fields: fields{
				apiURL:      mockServer.URL,
				apiKey:      "test",
				rateLimiter: &RateLimiter{},
			},
			args: args{
				log: zap.NewNop(),
				m: &types.Map{
					Seed:        "1",
					Size:        4000,
					SavedConfig: "default",
				},
			},
			want: &RustMapsStatusResponse{
				Meta: RustMapsStatusResponseMeta{
					Status:     common.StatusComplete,
					StatusCode: 200,
					Errors:     []string{},
				},
			},
		},
		{
			name: "Test GetStatus 401",
			fields: fields{
				apiURL:      mockServer.URL,
				apiKey:      "test",
				rateLimiter: &RateLimiter{},
			},
			args: args{
				log: zap.NewNop(),
				m: &types.Map{
					Seed:        "2",
					Size:        4000,
					SavedConfig: "default",
				},
			},
			want: &RustMapsStatusResponse{
				Meta: RustMapsStatusResponseMeta{
					Status:     common.StatusUnauthorized,
					StatusCode: 401,
					Errors:     []string{},
				},
			},
		},
		{
			name: "Test GetStatus 403",
			fields: fields{
				apiURL:      mockServer.URL,
				apiKey:      "test",
				rateLimiter: &RateLimiter{},
			},
			args: args{
				log: zap.NewNop(),
				m: &types.Map{
					Seed:        "3",
					Size:        4000,
					SavedConfig: "default",
				},
			},
			want: &RustMapsStatusResponse{
				Meta: RustMapsStatusResponseMeta{
					Status:     common.StatusForbidden,
					StatusCode: 403,
					Errors:     []string{},
				},
			},
		},
		{
			name: "Test GetStatus 404",
			fields: fields{
				apiURL:      mockServer.URL,
				apiKey:      "test",
				rateLimiter: &RateLimiter{},
			},
			args: args{
				log: zap.NewNop(),
				m: &types.Map{
					Seed:        "4",
					Size:        4000,
					SavedConfig: "default",
				},
			},
			want: &RustMapsStatusResponse{
				Meta: RustMapsStatusResponseMeta{
					Status:     common.StatusNotFound,
					StatusCode: 404,
					Errors:     []string{},
				},
			},
		},
		{
			name: "Test GetStatus 409",
			fields: fields{
				apiURL:      mockServer.URL,
				apiKey:      "test",
				rateLimiter: &RateLimiter{},
			},
			args: args{
				log: zap.NewNop(),
				m: &types.Map{
					Seed:        "5",
					Size:        4000,
					SavedConfig: "default",
				},
			},
			want: &RustMapsStatusResponse{
				Meta: RustMapsStatusResponseMeta{
					Status:     common.StatusGenerating,
					StatusCode: 409,
					Errors:     []string{},
				},
			},
		},
		{
			name: "Test GetStatus 500",
			fields: fields{
				apiURL:      mockServer.URL,
				apiKey:      "test",
				rateLimiter: &RateLimiter{},
			},
			args: args{
				log: zap.NewNop(),
				m: &types.Map{
					Seed:        "6",
					Size:        4000,
					SavedConfig: "default",
				},
			},
			want: &RustMapsStatusResponse{
				Meta: RustMapsStatusResponseMeta{
					Status:     "",
					StatusCode: 500,
					Errors:     []string{},
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &RustMapsClient{
				ApiUrl:      tt.fields.apiURL,
				apiKey:      tt.fields.apiKey,
				rateLimiter: tt.fields.rateLimiter,
			}
			got, err := c.GetStatus(tt.args.log, tt.args.m)
			if (err != nil) != tt.wantErr {
				t.Errorf("RustMapsClient.GetStatus() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if got != nil && got.Meta.StatusCode != tt.want.Meta.StatusCode {
				t.Errorf("RustMapsClient.GetStatus() = %v, want %v", got, tt.want)
			}
		})
	}
}
