package utils

import (
	"encoding/base64"
	"fmt"
	. "github.com/raver119/statika/classes"
	"io"
	"net/http"
	"os"
	"regexp"
	"time"
)

/*
	This function does exactly what it says
*/
func GetEnvOrDefault(variable string, defaultValue string) string {
	if val, ok := os.LookupEnv(variable); ok {
		return val
	} else {
		return defaultValue
	}
}

func GetEnvOrPanic(variable string) string {
	if val, ok := os.LookupEnv(variable); ok {
		return val
	} else {
		panic(fmt.Errorf("Environment variable wasn't found: [%v]", variable))
	}
}

func IsTimingEnabled() bool {
	if _, ok := os.LookupEnv("TIMING"); ok {
		return true
	} else {
		return false
	}
}

func FileExists(fileName string, fileOnly bool) bool {
	fi, err := os.Stat(fileName)
	if err != nil {
		return false
	}

	if !fileOnly {
		return true
	} else {
		return !fi.IsDir()
	}
}

func TransferBytes(r io.Reader, w io.Writer) error {
	buf := make([]byte, 16384)
	for {
		n, err := r.Read(buf)
		if err == io.EOF {
			if n > 0 {
				// write tail bytes
				_, err = w.Write(buf[:n])
				if err != nil {
					return err
				}
			}

			break
		}

		_, err = w.Write(buf[:n])
		if err != nil {
			return err
		}
	}

	return nil
}

func OptionallyReport(message string, w http.ResponseWriter, err error) (ok bool) {
	ok = err == nil
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(message + ": " + err.Error()))
	}

	return
}

/*
	This function sets all required headers in HTTP response
*/
func SetupCorsHeaders(w *http.ResponseWriter, req *http.Request) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	(*w).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
}

func SetupCacheHeaders(w *http.ResponseWriter, req *http.Request) {
	// no-cache part
	var epoch = time.Unix(0, 0).Format(time.RFC1123)
	(*w).Header().Set("Expires", epoch)
	(*w).Header().Set("Cache-Control", "no-cache, no-store, no-transform, must-revalidate, private, max-age=0")
	(*w).Header().Set("Pragma", "no-cache")
	(*w).Header().Set("X-Accel-Expires", "0")
}

func EncodePath(bucket string, fileName string) (b string, f string) {
	b = base64.StdEncoding.EncodeToString([]byte(bucket))
	// extract file extension
	re := regexp.MustCompile("(\\.[a-zA-Z0-9]+)$")
	match := re.FindStringSubmatch(fileName)

	if len(match) == 0 {
		f = base64.StdEncoding.EncodeToString([]byte(fileName))
	} else {
		// encode filename
		rep := re.ReplaceAllString(fileName, "")
		f = base64.StdEncoding.EncodeToString([]byte(rep)) + match[0]
	}

	return b, f
}

func CorsHandler(hf http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		SetupCorsHeaders(&w, r)
		if r.Method == http.MethodOptions {

			// do nothing else
		} else {
			SetupCacheHeaders(&w, r)
			// pass request to the actual handler
			hf.ServeHTTP(w, r)
		}
	}
}

func MasterFileName(bucket string, fileName string) string {
	b, f := EncodePath(bucket, fileName)
	return fmt.Sprintf("%v/%v", b, f)
}

func MasterMetaName(bucket string, fileName string) string {
	b := base64.StdEncoding.EncodeToString([]byte(bucket))
	f := base64.StdEncoding.EncodeToString([]byte(fileName)) + META_EXTENSION
	return fmt.Sprintf("%v/%v", b, f)
}

func TimeIt(f func()) int64 {
	start := time.Now().UnixNano() / 1000000
	f()
	return (time.Now().UnixNano() / 1000000) - start
}
