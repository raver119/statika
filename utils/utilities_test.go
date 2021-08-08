package utils

import "testing"

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
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotB, gotF := EncodePath(tt.args.bucket, tt.args.fileName)
			if gotB != tt.wantB {
				t.Errorf("EncodePath() gotB = %v, want %v", gotB, tt.wantB)
			}
			if gotF != tt.wantF {
				t.Errorf("EncodePath() gotF = %v, want %v", gotF, tt.wantF)
			}
		})
	}
}
