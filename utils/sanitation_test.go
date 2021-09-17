package utils

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSanitizeFileName(t *testing.T) {
	tests := []struct {
		name      string
		inputName string
		saneName  string
		wantErr   bool
	}{
		{"test_0", "filename.txt", "filename.txt", false},
		{"test_1", "/filename.txt", "filename.txt", false},
		{"test_2", "../filename.txt", "filename.txt", false},
		{"test_3", "../../filename.txt", "filename.txt", false},
		{"test_4", "../foo/../filename.txt", "filename.txt", false},
		{"test_5", `\filename.txt`, "filename.txt", false},
		{"test_6", "//filename.txt", "filename.txt", false},
		{"test_7", `\\filename.txt`, "filename.txt", false},
		{"test_8", `/alpha/filename.txt`, "alpha/filename.txt", false},
		{"test_9", `alpha/filename.txt`, "alpha/filename.txt", false},
		{"test_10", `/alpha/../filename.txt`, "filename.txt", false},
		{"test_11", `alpha\filename.txt`, "alpha/filename.txt", false},
		{"test_12", `/./../alpha/filename.txt`, "alpha/filename.txt", false},
		{"test_13", `./alpha/filename.txt`, "alpha/filename.txt", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotFname, err := SanitizeFileName(tt.inputName)
			if (err != nil) != tt.wantErr {
				t.Errorf("SanitizeFileName() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			require.Equal(t, tt.saneName, gotFname)
		})
	}
}

func TestExtractPath(t *testing.T) {
	tests := []struct {
		name     string
		fileName string
		wantPath string
		wantErr  bool
	}{
		{"test_0", "filename.txt", "", false},
		{"test_1", "alpha/filename.txt", "alpha", false},
		{"test_2", `alpha\filename.txt`, "alpha", false},
		{"test_3", "alpha/../../filename.txt", "", false},
		{"test_4", "alpha/beta/filename.txt", "alpha/beta", false},
		{"test_5", `\alpha\filename.txt`, "alpha", false},
		{"test_6", `/alpha/filename.txt`, "alpha", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotPath, err := ExtractPath(tt.fileName)
			if (err != nil) != tt.wantErr {
				t.Errorf("ExtractPath() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			require.Equal(t, tt.wantPath, gotPath)
		})
	}
}

func TestExtractFileName(t *testing.T) {
	tests := []struct {
		name         string
		fullPath     string
		wantFileName string
		wantErr      bool
	}{
		{"test_0", "filename.txt", "filename.txt", false},
		{"test_1", "filename", "filename", false},
		{"test_2", "alpha/filename", "filename", false},
		{"test_3", "alpha/beta/filename", "filename", false},
		{"test_4", "alpha/filename.txt", "filename.txt", false},
		{"test_5", "ALPHA/filename.txt", "filename.txt", false},
		{"test_6", "/filename.txt", "filename.txt", false},
		{"test_7", "/filename", "filename", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotFileName, err := ExtractFileName(tt.fullPath)
			if (err != nil) != tt.wantErr {
				t.Errorf("ExtractFileName() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			require.Equal(t, tt.wantFileName, gotFileName)
		})
	}
}

func TestExtractExtension(t *testing.T) {
	tests := []struct {
		name     string
		fullPath string
		wantExt  string
		wantErr  bool
	}{
		{"test_0", "abs/filename.txt", "txt", false},
		{"test_1", "abs/filename", "", false},
		{"test_2", "abs/filename.txt.exe", "exe", false},
		{"test_3", "abs/filename.txt..exe", "exe", false},
		{"test_4", "abs.exe/filename.txt", "txt", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotExt, err := ExtractExtension(tt.fullPath)
			if (err != nil) != tt.wantErr {
				t.Errorf("ExtractExtension() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			require.Equal(t, tt.wantExt, gotExt)
		})
	}
}

func TestSplitPath(t *testing.T) {
	tests := []struct {
		name         string
		path         string
		wantBucket   string
		wantFileName string
		wantErr      bool
	}{
		{"test_0", "images/alpha.png", "images", "alpha.png", false},
		{"test_1", "images/sub/alpha.png", "images", "sub/alpha.png", false},
		{"test_2", "images/../alpha.png", "images", "../alpha.png", false},
		{"test_3", "/images/alpha.png", "images", "alpha.png", false},
		{"test_4", "/images/sub/alpha.png", "images", "sub/alpha.png", false},
		{"test_5", "/images/../alpha.png", "images", "../alpha.png", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotBucket, gotFileName, err := SplitPath(tt.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("SplitPath() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			require.Equal(t, tt.wantBucket, gotBucket)
			require.Equal(t, tt.wantFileName, gotFileName)
		})
	}
}
