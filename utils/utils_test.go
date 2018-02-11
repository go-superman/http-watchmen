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
				url:       "http://47.91.255.0/tes5.jpg",
				retryCnt:  3,
				retryTime: 0,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if data, err := HealthCheck(tt.args.url, tt.args.retryCnt, []int{}, 5*time.Second, tt.args.retryTime); (err != nil) != tt.wantErr {
				t.Errorf("HealthCheck() error = %v, wantErr %v %v", err, tt.wantErr, data)
			}else{
				t.Log(data)
			}
		})
	}
}
