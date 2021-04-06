package main

import (
	"fmt"
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

	/*
		read root folder that will be used for storage
	*/
	rootFolder, ok := os.LookupEnv("ROOT_DIR")
	if !ok {
		panic(fmt.Errorf("ROOT_DIR env var wasn't set"))
	}

	/*
		By default port 8080 is used
	*/
	strPort := GetEnvOrDefault("STATIKA_PORT", "9191")
	port, err := strconv.Atoi(strPort)
	if err != nil {
		panic(err)
	}

	var ls Storage = NewLocalStorage(rootFolder)

	/*
		Create server instance and start it
	*/
	fmt.Printf("Starting Statika server at port [%v], serving %v folder\n", port, rootFolder)
	engine, err := CreateEngine(keyMaster, keyUpload, &ls, port)
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