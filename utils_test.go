//go:build unit

package flickr

import (
	"errors"
	"fmt"
	"io"
	"os"
	"testing"
)

func Test_getApiKey(t *testing.T) {
	// Create .env file
	const envKey = "xyzpdq"
	f, err := os.CreateTemp(".", "")
	if err != nil {
		t.Fatalf("Unable to create temp env file: %s", err)
	}
	defer f.Close()
	defer os.Remove(f.Name())
	io.WriteString(f, fmt.Sprintf("%s=%s\n", ApiKeyEnvVar, envKey))
	fmt.Printf("env file is %s\n", f.Name())

	type args struct {
		key         string
		envFileName string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "Provided key",
			args: args{
				key: "abc123", envFileName: "",
			},
			want:    "abc123",
			wantErr: false,
		},
		{
			name: "Missing specified env file",
			args: args{
				key: "", envFileName: "XXXXX",
			},
			want:    "",
			wantErr: true,
		},
		{
			name: "Specified env file",
			args: args{
				key:         "",
				envFileName: f.Name(),
			},
			want:    envKey,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getApiKey(tt.args.key, tt.args.envFileName)
			if (err != nil) != tt.wantErr {
				t.Errorf("getApiKey() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("getApiKey() = '%v', want '%v'", got, tt.want)
			}
		})
	}
}

func Test_getApiKey_EnvironmentKeySet(t *testing.T) {
	envVal := "xyzpdq"
	t.Setenv(ApiKeyEnvVar, envVal)
	type args struct {
		key         string
		envFileName string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "Provided key",
			args: args{
				key: "abc123", envFileName: "",
			},
			want:    "abc123",
			wantErr: false,
		},
		{
			name: "No provided key",
			args: args{
				key: "", envFileName: "",
			},
			want:    envVal,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getApiKey(tt.args.key, tt.args.envFileName)
			if (err != nil) != tt.wantErr {
				t.Errorf("getApiKey() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("getApiKey() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getApiKey_NoEnvironmentKeySet(t *testing.T) {
	os.Unsetenv(ApiKeyEnvVar)
	_, err := getApiKey("", "")
	if err == nil || !errors.Is(err, ErrKeyNotInEnv) {
		t.Errorf("Expecting error '%s', but got '%s'", ErrKeyNotInEnv, err)
	}
}
