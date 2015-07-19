package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"
)

type Deployment struct {
	DeployedOn time.Time
	Commit     string
	AppID      string
}

// For sorting our deployments by date
type ByDate []Deployment

func (a ByDate) Len() int           { return len(a) }
func (a ByDate) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByDate) Less(i, j int) bool { return a[j].DeployedOn.Before(a[i].DeployedOn) }

// Peform the verification steps and gather the info before creating a new Deployment
func newDeployment() (Deployment, error) {
	//check nothing is staged
	if !isEverythingCommitted() {
		return Deployment{}, fmt.Errorf("commit your stuff first!")
	}

	//get commit hash
	hash, err := getCommitHash()
	if err != nil {
		return Deployment{}, err
	}

	//get current application from app.yaml
	appId, err := getApplicationId()
	if err != nil {
		return Deployment{}, err
	}

	d := Deployment{
		DeployedOn: time.Now(),
		Commit:     hash,
		AppID:      appId,
	}

	return d, nil
}

func getCommitHash() (string, error) {
	//git log --pretty=format:'%h' -n 1
	hash, err := exec.Command("git", "log", "--pretty=format:%h", "-n 1").Output()
	if err != nil {
		log.Printf("getCommitHash() : Could not get hash of last commit: %v", err)
		return "", err
	}

	return string(hash), nil
}

func isDeployable() bool {

	//is git dir?
	_, err := exec.Command("git", "rev-parse").Output()
	if err != nil {
		log.Printf("isDeployable() : is not a git directory: %v", err)
		return false
	}

	//has app.yaml?
	if _, err := os.Stat("app.yaml"); os.IsNotExist(err) {
		log.Printf("Directory does not contain an app.yaml")
		return false
	}

	return true
}

func isEverythingCommitted() bool {
	out, err := exec.Command("git", "status", "--porcelain").Output()
	if err != nil {
		log.Printf("isEverythingCommitted() : could not call git status: %v", err)
		return false
	}

	if string(out) != "" {
		log.Printf("isEverythingCommitted() : uncommitted changes")
		return false
	}
	return true
}

func getApplicationId() (string, error) {
	// Open an input file, exit on error.
	inputFile, err := os.Open("app.yaml")
	if err != nil {
		log.Fatal("Error opening app.yaml:", err)
	}
	defer inputFile.Close()

	r := bufio.NewReader(inputFile)

	line, isPrefix, err := r.ReadLine()
	if err != nil || isPrefix {
		log.Printf("Could not readline from the app.yaml: %v", err)
		return "", err
	}

	parts := strings.Split(string(line), ": ")
	if len(parts) != 2 || parts[0] != "application" {
		log.Printf("Invalid app.yaml, could not parse the application line : %s", string(line))
		return "", fmt.Errorf("Invalid app.yaml, could not parse the application line : %s", string(line))
	}

	return parts[1], nil
}
