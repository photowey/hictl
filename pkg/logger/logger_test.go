package logger

import (
	"testing"
)

func TestInfo(t *testing.T) {
	type args struct {
		message string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "Test logger info level",
			args: args{
				message: "Hello, world",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Info(tt.args.message)
		})
	}
}

func TestInfof(t *testing.T) {
	type args struct {
		message string
		args    []any
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "Test logger info level",
			args: args{
				message: "Hello, %s",
				args:    []any{"world"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Infof(tt.args.message, tt.args.args...)
		})
	}
}

func TestDebugf(t *testing.T) {
	type args struct {
		message string
		file    string
		line    int
		args    []any
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "Test logger info level",
			args: args{
				message: "Hello, %s",
				file:    "logger_test",
				line:    76,
				args:    []any{"world"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Debugf(tt.args.message, tt.args.file, tt.args.line, tt.args.args...)
		})
	}
}
