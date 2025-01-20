package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/maintc/rustmaps-cli/pkg/types"
	"go.uber.org/zap"
)

func TestRustMapsClient_GenerateCustom(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			if strings.HasPrefix(r.URL.Path, "/maps/custom/saved-config") {
				w.Header().Set("Content-Type", "application/json")
				// Serialize the request body to RustMapsGenerateCustomRequest
				var req RustMapsGenerateCustomRequest
				err := json.NewDecoder(r.Body).Decode(&req)
				if err != nil {
					http.Error(w, err.Error(), http.StatusBadRequest)
					return
				}

				var response *RustMapsGenerateResponse
				if req.MapParameters.Seed == "1" {
					w.WriteHeader(http.StatusOK)
					response = &RustMapsGenerateResponse{
						Meta: RustMapsGenerateResponseMeta{
							Status:     "complete",
							StatusCode: 200,
							Errors:     []string{},
						},
					}
				} else if req.MapParameters.Seed == "2" {
					w.WriteHeader(http.StatusUnauthorized)
					response = &RustMapsGenerateResponse{
						Meta: RustMapsGenerateResponseMeta{
							Status:     "complete",
							StatusCode: 401,
							Errors:     []string{},
						},
					}
				} else if req.MapParameters.Seed == "3" {
					w.WriteHeader(http.StatusForbidden)
					response = &RustMapsGenerateResponse{
						Meta: RustMapsGenerateResponseMeta{
							Status:     "complete",
							StatusCode: 403,
							Errors:     []string{},
						},
					}
				} else if req.MapParameters.Seed == "4" {
					w.WriteHeader(http.StatusConflict)
					response = &RustMapsGenerateResponse{
						Meta: RustMapsGenerateResponseMeta{
							Status:     "complete",
							StatusCode: 409,
							Errors:     []string{},
						},
					}
				} else {
					w.WriteHeader(http.StatusInternalServerError)
					response = &RustMapsGenerateResponse{
						Meta: RustMapsGenerateResponseMeta{
							Status:     "error",
							StatusCode: 500,
							Errors:     []string{"Invalid seed"},
						},
					}
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
		want    *RustMapsGenerateResponse
		wantErr bool
	}{
		{
			name: "GenerateCustom 200",
			fields: fields{
				apiURL:      mockServer.URL,
				apiKey:      "test",
				rateLimiter: &RateLimiter{},
			},
			args: args{
				log: zap.NewNop(),
				m: &types.Map{
					Size:        3500,
					Seed:        "1",
					Staging:     false,
					SavedConfig: "default",
				},
			},
			want: &RustMapsGenerateResponse{
				Meta: RustMapsGenerateResponseMeta{
					Status:     "complete",
					StatusCode: 200,
					Errors:     []string{},
				},
			},
		},
		{
			name: "GenerateCustom 401",
			fields: fields{
				apiURL:      mockServer.URL,
				apiKey:      "test",
				rateLimiter: &RateLimiter{},
			},
			args: args{
				log: zap.NewNop(),
				m: &types.Map{
					Size:        3500,
					Seed:        "2",
					Staging:     false,
					SavedConfig: "default",
				},
			},
			want: &RustMapsGenerateResponse{
				Meta: RustMapsGenerateResponseMeta{
					Status:     "complete",
					StatusCode: 401,
					Errors:     []string{},
				},
			},
			wantErr: true,
		},
		{
			name: "GenerateCustom 403",
			fields: fields{
				apiURL:      mockServer.URL,
				apiKey:      "test",
				rateLimiter: &RateLimiter{},
			},
			args: args{
				log: zap.NewNop(),
				m: &types.Map{
					Size:        3500,
					Seed:        "3",
					Staging:     false,
					SavedConfig: "default",
				},
			},
			want: &RustMapsGenerateResponse{
				Meta: RustMapsGenerateResponseMeta{
					Status:     "complete",
					StatusCode: 403,
					Errors:     []string{},
				},
			},
			wantErr: true,
		},
		{
			name: "GenerateCustom 409",
			fields: fields{
				apiURL:      mockServer.URL,
				apiKey:      "test",
				rateLimiter: &RateLimiter{},
			},
			args: args{
				log: zap.NewNop(),
				m: &types.Map{
					Size:        3500,
					Seed:        "4",
					Staging:     false,
					SavedConfig: "default",
				},
			},
			want:    nil,
			wantErr: false,
		},
		{
			name: "GenerateCustom 500",
			fields: fields{
				apiURL:      "http://localhost",
				apiKey:      "test",
				rateLimiter: &RateLimiter{},
			},
			args: args{
				log: zap.NewNop(),
				m: &types.Map{
					Size:        3500,
					Seed:        "100",
					Staging:     false,
					SavedConfig: "default",
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
			got, err := c.GenerateCustom(tt.args.log, tt.args.m)
			if (err != nil) != tt.wantErr {
				t.Errorf("RustMapsClient.GenerateCustom() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != nil && got.Meta.StatusCode != tt.want.Meta.StatusCode {
				t.Errorf("RustMapsClient.GenerateCustom() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRustMapsClient_GenerateProcedural(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			if strings.HasPrefix(r.URL.Path, "/maps") {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				response := &RustMapsGenerateResponse{
					Meta: RustMapsGenerateResponseMeta{
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
		m   *types.Map
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *RustMapsGenerateResponse
		wantErr bool
	}{
		{
			name: "GenerateProcedural 200",
			fields: fields{
				apiURL:      mockServer.URL,
				apiKey:      "test",
				rateLimiter: &RateLimiter{},
			},
			args: args{
				log: zap.NewNop(),
				m: &types.Map{
					Size:    3500,
					Seed:    "1",
					Staging: false,
				},
			},
			want: &RustMapsGenerateResponse{
				Meta: RustMapsGenerateResponseMeta{
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
			got, err := c.GenerateProcedural(tt.args.log, tt.args.m)
			if (err != nil) != tt.wantErr {
				t.Errorf("RustMapsClient.GenerateProcedural() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != nil && got.Meta.StatusCode != tt.want.Meta.StatusCode {
				t.Errorf("RustMapsClient.GenerateCustom() = %v, want %v", got, tt.want)
			}
		})
	}
}
