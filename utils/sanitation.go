package utils

import (
	"fmt"
	"path/filepath"
	"regexp"
	"strings"
)

var re = regexp.MustCompile(`[a-zA-Z0-9-]+`)
var ref = regexp.MustCompile(`^/{1,}`)

func SplitPath(path string) (bucket string, fileName string, err error) {
	path = strings.ReplaceAll(path, `\`, "/")
	path = ref.ReplaceAllString(path, "")
	chunks := strings.Split(path, "/")

	bucket = chunks[0]
	if !re.MatchString(bucket) {
		return "", "", fmt.Errorf("bucket looks bad: [%v]", bucket)
	}

	if len(chunks) > 1 {
		fileName = chunks[1]

		for i := 2; i < len(chunks); i++ {
			fileName = fmt.Sprintf("%v/%v", fileName, chunks[i])
		}
	}

	return
}

// SanitizeFileName function does exactly what it says
func SanitizeFileName(fileName string) (fname string, err error) {
	// convert all slashes to the single format
	fname = strings.ReplaceAll(fileName, `\`, "/")

	// treat file as if it's mapped to the root folder
	fname = filepath.Clean("/" + fname)

	// remove
	fname = strings.Replace(fname, "/", "", 1)
	return
}

// ExtractPath function returns a path component out of given filename.
// if filename has no path component, it returns an empty string
func ExtractPath(fileName string) (path string, err error) {
	path, err = SanitizeFileName(fileName)
	if err != nil {
		return "", err
	}

	path = filepath.Dir(path)

	if path == "." {
		path = ""
	}

	return
}

// ExtractFileName function returns a filename component out of a given path
func ExtractFileName(fullPath string) (fileName string, err error) {
	path, err := ExtractPath(fullPath)
	if err != nil {
		return "", err
	}

	fileName = strings.Replace(fullPath, path, "", 1)
	if strings.HasPrefix(fileName, "/") {
		fileName = strings.Replace(fileName, "/", "", 1)
	}

	return
}

// ExtractExtension function returns an extension (if any) from the given file name
func ExtractExtension(fullPath string) (ext string, err error) {
	fileName, err := ExtractFileName(fullPath)
	if err != nil {
		return "", err
	}

	if !strings.Contains(fileName, ".") {
		return "", err
	} else {
		re := regexp.MustCompile(`^.*\.`)
		ext = re.ReplaceAllString(fileName, "")
	}

	return
}

func Reverse(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}
