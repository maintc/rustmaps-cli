package api

import (
	"reflect"
	"sync"
	"testing"
	"time"
)

func TestNewRustMapsClient(t *testing.T) {
	type args struct {
		apiKey string
	}
	tests := []struct {
		name string
		args args
		want *RustMapsClient
	}{
		{
			name: "NewRustMapsClient",
			args: args{
				apiKey: "test",
			},
			want: &RustMapsClient{
				ApiUrl: "https://api.rustmaps.com/v4",
				apiKey: "test",
				rateLimiter: &RateLimiter{
					callsPerMinute: 60,
					interval:       time.Minute / time.Duration(60),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewRustMapsClient(tt.args.apiKey); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewRustMapsClient() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRustMapsClient_SetApiKey(t *testing.T) {
	type fields struct {
		apiKey      string
		rateLimiter *RateLimiter
	}
	type args struct {
		apiKey string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "SetApiKey",
			args: args{
				apiKey: "test",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &RustMapsClient{
				apiKey:      tt.fields.apiKey,
				rateLimiter: tt.fields.rateLimiter,
			}
			r.SetApiKey(tt.args.apiKey)
		})
	}
}

func TestNewRateLimiter(t *testing.T) {
	type args struct {
		callsPerMinute int
	}
	tests := []struct {
		name string
		args args
		want *RateLimiter
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewRateLimiter(tt.args.callsPerMinute); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewRateLimiter() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRateLimiter_Wait(t *testing.T) {
	type fields struct {
		callsPerMinute int
		interval       time.Duration
		lastCall       time.Time
		mu             sync.Mutex
	}
	tests := []struct {
		name   string
		fields fields
	}{
		{
			name: "Wait",
			fields: fields{
				lastCall: time.Now(),
				interval: 100000000,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &RateLimiter{
				callsPerMinute: tt.fields.callsPerMinute,
				interval:       tt.fields.interval,
				lastCall:       tt.fields.lastCall,
				mu:             tt.fields.mu,
			}
			r.Wait()
		})
	}
}
