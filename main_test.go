package main

import (
	"testing"
)

func Test_httpPostJson(t *testing.T) {
	type args struct {
		jsonstr []byte
		url     string
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			httpPostJson(tt.args.jsonstr, tt.args.url)
		})
	}
}

func Test_main(t *testing.T) {
	tests := []struct {
		name string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			main()
		})
	}
}

func Test_getPdfPath(t *testing.T) {
	tests := []struct {
		name string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			getPdfPath()
		})
	}
}

func Test_readConfigFile(t *testing.T) {
	type args struct {
		dir string
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			readConfigFile(tt.args.dir)
		})
	}
}

func Test_postFile(t *testing.T) {
	type args struct {
		filename  string
		targetUrl string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := postFile(tt.args.filename, tt.args.targetUrl); (err != nil) != tt.wantErr {
				t.Errorf("postFile() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
