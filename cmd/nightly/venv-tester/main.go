package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"time"

	"github.com/Netflix/go-expect"
)

const (
	expectedString        = "aws-node"
	defaultTimeoutSeconds = 30
)

func main() {
	args := os.Args[1:3]

	if len(args) != 2 {
		printUsage()

		os.Exit(1)
	}

	for _, path := range args {
		_, err := os.Stat(path)
		if err != nil {
			crashPrint(err)
		}
	}

	okctlBinaryPath := args[0]
	clusterManifestPath := args[1]

	err := testVenv(okctlBinaryPath, clusterManifestPath)
	if err != nil {
		crashPrint(err)
	}
}

func testVenv(okctlBinaryPath string, clusterManifestPath string) error {
	c, err := expect.NewConsole(
		expect.WithDefaultTimeout(defaultTimeoutSeconds*time.Second),
		expect.WithStdout(os.Stdout),
	)
	if err != nil {
		return fmt.Errorf("creating pseudo console: %w", err)
	}

	defer c.Close()

	cmd := exec.Command(okctlBinaryPath, "venv", "-c", clusterManifestPath)
	cmd.Stdin = c.Tty()
	cmd.Stdout = c.Tty()
	cmd.Stderr = c.Tty()

	log.Println("Starting venv")

	err = cmd.Start()
	if err != nil {
		return fmt.Errorf("starting command: %w", err)
	}

	log.Println("Running kubectl -n kube-system get pods")

	_, err = c.SendLine("kubectl -n kube-system get pods")
	if err != nil {
		return fmt.Errorf("sending kubectl command: %w", err)
	}

	log.Println("Running exit")

	_, err = c.SendLine("exit")
	if err != nil {
		return fmt.Errorf("sending exit command: %w", err)
	}

	log.Printf("Looking for a match for %s", expectedString)

	_, err = c.ExpectString(expectedString)
	if err != nil {
		return fmt.Errorf("finding expected string: %w", err)
	}

	log.Println("Found match")

	err = cmd.Wait()
	if err != nil {
		return fmt.Errorf("waiting for command to return")
	}

	return nil
}

func printUsage() {
	fmt.Println("venv-tester")
	fmt.Println()
	fmt.Println("Opens okctl venv, runs kubectl -n kube-system get pods and exits")
	fmt.Println("with a status of 0 or 1 depending on if it finds \"aws-node\"")
	fmt.Println()
	fmt.Println("USAGE:")
	fmt.Println("\tvenv-tester <okctl binary path> <cluster manifest path>")
	fmt.Println("Example:")
	fmt.Println("\tvenv-tester /usr/bin/okctl cluster.yaml")
	fmt.Println()
}

func crashPrint(err error) {
	fmt.Printf("Error: %s\n", err.Error())

	os.Exit(1)
}
