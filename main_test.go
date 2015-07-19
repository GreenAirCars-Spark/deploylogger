package main

import (
	"fmt"
	"os"
	"os/exec"
)

func init() {
	// Create the good test data if it doesn't already exist
	if _, err := os.Stat(goodTestData); os.IsNotExist(err) {

		fmt.Println("Creating test data")

		//create good test data
		oldWD, _ := os.Getwd()
		err := exec.Command("mkdir", goodTestData).Run()
		if err != nil {
			panic("could not create dir")
		}

		err = exec.Command("cp", "testdata/app.yaml", "testdata/good/app.yaml").Run()
		if err != nil {
			panic("could not move app.yaml")
		}

		setWD(goodTestData)
		err = exec.Command("git", "init").Run()
		if err != nil {
			panic("could not init git dir")
		}

		err = exec.Command("git", "add", "-A").Run()
		if err != nil {
			panic("could not add")
		}

		err = exec.Command("git", "commit", "-m", "blah blah").Run()
		if err != nil {
			panic("could not commit")
		}
		setWD(oldWD)
	}
}
