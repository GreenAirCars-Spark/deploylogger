package main

import (
	"os"
	"reflect"
	"sort"
	"testing"
)

const (
	goodTestData = "testdata/good"
	badTestData  = "testdata/bad"
)

func setWD(dir string) {
	err := os.Chdir(dir)
	if err != nil {
		panic("uh oh")
	}
}

func TestIsDeployable(t *testing.T) {
	oldWD, _ := os.Getwd()
	setWD(badTestData)

	if isDeployable() {
		t.Fail()
	}

	setWD(oldWD)
	setWD(goodTestData)
	if !isDeployable() {
		t.Fail()
	}
	setWD(oldWD)
}

func TestIsEverythingCommitted(t *testing.T) {
	oldWD, _ := os.Getwd()
	setWD(badTestData)

	if isEverythingCommitted() {
		t.Fail()
	}

	setWD(oldWD)
	setWD(goodTestData)
	if !isEverythingCommitted() {
		t.Fatal("aa")
		t.Fail()
	}
	setWD(oldWD)
}

func TestGetCommitHash(t *testing.T) {
	oldWD, _ := os.Getwd()
	setWD(goodTestData)

	//get commit hash
	hash, err := getCommitHash()
	if err != nil {
		t.Fail()
	}

	if hash != "07a8a47" {
		t.Fatal("could not get commit hash")
	}
	setWD(oldWD)
}

func TestGetAppID(t *testing.T) {
	oldWD, _ := os.Getwd()
	setWD(goodTestData)

	//get current application from app.yaml
	appId, err := getApplicationId()
	if err != nil {
		t.Fail()
	}

	if appId != "fake-app-id" {
		t.Fatalf("could not get app id: %s", appId)
	}

	setWD(oldWD)
}

func TestLogDeployments(t *testing.T) {

	oldWD, _ := os.Getwd()
	setWD(goodTestData)

	deployments := make([]Deployment, 0, 10)
	for i := 0; i < 10; i++ {
		d, _ := newDeployment()
		deployments = append(deployments, d)
	}
	sort.Sort(ByDate(deployments))

	err := setDeployments(deployments)
	if err != nil {
		t.Fatal("could not write deployments")
	}

	writtenDeployments, err := getDeployments()
	if err != nil {
		t.Fatal("could not get deployments")
	}

	if !reflect.DeepEqual(deployments, writtenDeployments) {
		t.Fatal("what was written is not what was read")
	}

	os.Remove(DeploymentsFileName)

	setWD(oldWD)
}
