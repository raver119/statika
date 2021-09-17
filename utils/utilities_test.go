package utils

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestEncodePath(t *testing.T) {
	type args struct {
		bucket   string
		fileName string
	}
	tests := []struct {
		name  string
		args  args
		wantB string
		wantF string
	}{
		{name: "test_0", args: args{bucket: "alpha", fileName: "filename.txt"}, wantB: "YWxwaGE=", wantF: "ZmlsZW5hbWU=.txt"},
		{name: "test_1", args: args{bucket: "alpha", fileName: "filename"}, wantB: "YWxwaGE=", wantF: "ZmlsZW5hbWU="},
		{name: "test_2", args: args{bucket: "alpha", fileName: "filename.jpg.txt"}, wantB: "YWxwaGE=", wantF: "ZmlsZW5hbWUuanBn.txt"},
		{name: "test_3", args: args{bucket: "alpha", fileName: "beta/filename.txt"}, wantB: "YWxwaGE=", wantF: "YmV0YS9maWxlbmFtZQ==.txt"},
		{name: "test_4", args: args{bucket: "alpha", fileName: "beta/gamma/filename.txt"}, wantB: "YWxwaGE=", wantF: "YmV0YS9nYW1tYS9maWxlbmFtZQ==.txt"},
		{name: "test_5", args: args{bucket: "alpha", fileName: "beta/../filename.txt"}, wantB: "YWxwaGE=", wantF: "ZmlsZW5hbWU=.txt"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotB, gotF := EncodePath(tt.args.bucket, tt.args.fileName)

			require.Equal(t, tt.wantB, gotB)
			require.Equal(t, tt.wantF, gotF)
		})
	}
}
