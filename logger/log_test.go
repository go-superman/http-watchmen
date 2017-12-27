package logger

import (
	"github.com/astaxie/beego/logs"
	"testing"
)

func TestResetLevel(t *testing.T) {
	type args struct {
		level int
	}
	tests := []struct {
		name string
		args args
	}{
	// TODO: Add test cases.
		{
			name: "LevelDebug",
			args: args{
				level: logs.LevelDebug,
			},
		},
		{
			name: "LevelInfo",
			args: args{
				level: logs.LevelInfo,
			},
		},
		{
			name: "LevelWarn",
			args: args{
				level: logs.LevelWarn,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ResetLevel(tt.args.level)
			Debugf("name:%v level:%v", tt.name, tt.args)
			Infof("name:%v level:%v", tt.name, tt.args)
			Warnf("name:%v level:%v", tt.name, tt.args)
		})
	}
}
