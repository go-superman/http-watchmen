package utils

import (
	"testing"
	"time"
)

func TestHealthCheck(t *testing.T) {
	type args struct {
		url       string
		retryCnt  int
		retryTime time.Duration
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "",
			args: args{
				url:       "http://47.91.255.0/api/index",
				retryCnt:  3,
				retryTime: 5 * time.Second,
			},
			wantErr: false,
		},
		{
			name: "img",
			args: args{
				url:       "http://47.91.255.0/test.jpg",
				retryCnt:  3,
				retryTime: 5 * time.Second,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := HealthCheck(tt.args.url, tt.args.retryCnt, tt.args.retryTime); (err != nil) != tt.wantErr {
				t.Errorf("HealthCheck() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
