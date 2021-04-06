package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
)

func main() {
	/*
		read tokens first:
		MASTER_KEY - used for administrative actions only, can't be used for file uploads etc
		UPLOAD_KEY - used for temporary tokens generation. can't be used directly.
	*/

	keyMaster, ok := os.LookupEnv("MASTER_KEY")
	if !ok || len(keyMaster) == 0 {
		panic(fmt.Errorf("MASTER_KEY env var wasn't set"))
	}

	keyUpload, ok := os.LookupEnv("UPLOAD_KEY")
	if !ok || len(keyUpload) == 0 {
		panic(fmt.Errorf("UPLOAD_KEY env var wasn't set"))
	}

	var storage Storage
	var err error

	/*
		By default port 8080 is used
	*/
	strPort := GetEnvOrDefault("STATIKA_PORT", "9191")
	port, err := strconv.Atoi(strPort)
	if err != nil {
		panic(err)
	}

	if _, ok := os.LookupEnv("S3_BUCKET"); ok {
		bucket := GetEnvOrPanic("S3_BUCKET")
		region := GetEnvOrPanic("S3_REGION")
		endpoint := GetEnvOrPanic("S3_ENDPOINT")
		_ = GetEnvOrPanic("S3_KEY")
		_ = GetEnvOrPanic("S3_SECRET")
		storage, err = NewS3Storage(bucket, endpoint, region)
		if err != nil {
			panic(err)
		}

		log.Printf("Starting Statika server at port [%v], serving %v/%v\n", port, endpoint, bucket)
	} else {
		/*
			read root folder that will be used for storage
		*/
		rootFolder, ok := os.LookupEnv("ROOT_DIR")
		if !ok {
			panic(fmt.Errorf("ROOT_DIR env var wasn't set"))
		}

		log.Printf("Starting Statika server at port [%v], serving %v folder\n", port, rootFolder)
		storage = NewLocalStorage(rootFolder)
	}

	/*
		Create server instance and start it
	*/
	engine, err := CreateEngine(keyMaster, keyUpload, &storage, port)
	if err != nil {
		fmt.Printf("CreateEngine failed: %v\n", err.Error())
		panic(err)
	}

	err = engine.Start()
	if err != nil {
		fmt.Printf("Engine::Start failed: %v\n", err.Error())
		panic(err)
	}

	fmt.Printf("Gracefully exiting...\n")
}
