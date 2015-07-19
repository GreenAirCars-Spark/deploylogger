package main

import (
	"encoding/json"
	"fmt"
	"github.com/codegangsta/cli"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"sort"
	"strconv"
)

const (
	DeploymentsFileName = "deployments.json"
)

var (
	ErrDeploymentsFileNotPresent = fmt.Errorf("Deployments file does not exist")
)

func main() {
	app := cli.NewApp()
	app.Name = "deploylogger"
	app.Usage = "log what we deploy"
	app.Commands = []cli.Command{
		{
			Name:    "deploy",
			Aliases: []string{"d"},
			Usage:   "add a task to the list",
			Action: func(c *cli.Context) {
				if !isDeployable() {
					log.Fatal("Directory is not deployable")
				}

				deployment, err := newDeployment()
				if err != nil {
					return
				}

				err = goappDeploy()
				if err != nil {
					return
				}

				deployments, err := getDeployments()
				if err != nil && err != ErrDeploymentsFileNotPresent {
					log.Fatalf("Could not get deployments %v", err)
				}

				deployments = append(deployments, deployment)
				sort.Sort(ByDate(deployments))

				err = setDeployments(deployments)
				if err != nil {
					log.Fatalf("Could not save deployments %v", err)
				}

			},
		},
		{
			Name:    "list",
			Aliases: []string{"l"},
			Usage:   "add a task to the list",
			Action: func(c *cli.Context) {
				if !isDeployable() {
					log.Fatal("Directory is not deployable")
				}

				deployments, err := getDeployments()
				if err != nil {
					log.Fatalf("Could not get deployments %v", err)
				}

				if len(deployments) == 0 {
					fmt.Println("No deployments to list :'(")
					return
				}

				if c.Args().Present() {
					n, err := strconv.Atoi(c.Args().First())
					if err != nil && n > 0 {
						deployments = deployments[:n]
					}
				}

				fmt.Printf("Time\t\t\t\t   Hash\t\t   App ID\n")
				fmt.Printf("----------------------------------------------------------\n")

				for _, val := range deployments {
					fmt.Printf("%s\t   %s\t   %s\n", val.DeployedOn.Format("Jan 2, 2006 at 3:04pm (MST)"), val.Commit, val.AppID)
				}
			},
		},
	}

	app.Run(os.Args)
}

func getDeployments() ([]Deployment, error) {

	deployments := make([]Deployment, 0)
	if _, err := os.Stat(DeploymentsFileName); os.IsNotExist(err) {
		return deployments, ErrDeploymentsFileNotPresent
	}

	data, err := ioutil.ReadFile(DeploymentsFileName)
	if err != nil {
		log.Printf("Could not read deployments file ")
		return deployments, fmt.Errorf("Could not read deployments file ")
	}

	err = json.Unmarshal(data, &deployments)
	if err != nil {
		log.Printf("Could not Unmarshal deployments file err %v", err)
		return deployments, fmt.Errorf("Could not Unmarshal deployments file %s", string(data))
	}

	return deployments, nil
}

func setDeployments(deployments []Deployment) error {
	data, err := json.MarshalIndent(deployments, "", "  ")
	if err != nil {
		log.Printf("Could not marshal deployments %v", err)
		return err
	}

	err = ioutil.WriteFile(DeploymentsFileName, data, 0644)
	if err != nil {
		log.Printf("Could not write to deployments file %v", err)
		return err
	}

	return nil
}

func goappDeploy() error {
	cmd := exec.Command("goapp", "deploy")
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	err := cmd.Run()

	if err != nil {
		log.Printf("Could not run goapp deploy: %v", err)
		return err
	}

	return nil
}
