package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"syscall"
	"time"

	"github.com/Netflix/go-expect" //nolint:typecheck
)

func main() {
	if len(os.Args) != 4 {
		printUsage()

		os.Exit(1)
	}

	rawArgs := os.Args[1:]

	args := applicationOpts{
		OkctlBinaryPath:     rawArgs[0],
		ClusterManifestPath: rawArgs[1],
		DatabaseName:        rawArgs[2],
	}

	err := args.Validate()
	if err != nil {
		crashPrint(err)
	}

	err = testForwardPostgres(args)
	if err != nil {
		crashPrint(err)
	}
}

func testForwardPostgres(opts applicationOpts) error {
	c, err := expect.NewConsole(
		expect.WithDefaultTimeout(defaultTimeoutSeconds*time.Second),
		expect.WithStdout(os.Stdout),
	)
	if err != nil {
		return fmt.Errorf("creating pseudo console: %w", err)
	}

	defer c.Close()

	passPath, err := generatePassfile()
	if err != nil {
		return fmt.Errorf("generating password file: %w", err)
	}

	args := []string{
		"forward", "postgres",
		"-c", opts.ClusterManifestPath,
		"-n", opts.DatabaseName,
		"-u", "testuser",
		"-p", passPath,
	}

	cmd := exec.Command(opts.OkctlBinaryPath, args...)
	cmd.Stdin = c.Tty()
	cmd.Stdout = c.Tty()
	cmd.Stderr = c.Tty()

	log.Println("Starting forward postgres test")

	err = cmd.Start()
	if err != nil {
		return fmt.Errorf("starting command: %w", err)
	}

	log.Printf("Looking for a match for %s", expectedString)

	_, err = c.ExpectString(expectedString)
	if err != nil {
		return fmt.Errorf("finding expected string: %w", err)
	}

	log.Println("Found match")

	log.Println("Sending SIGINT")

	err = cmd.Process.Signal(syscall.SIGINT)
	if err != nil {
		return fmt.Errorf("sending SIGINT")
	}

	err = cmd.Wait()
	// Ideally, since SIGINT is the expected way to shutdown the process, the exit code should be 0. Due to a bug in the
	// teardown process, forward postgres returns 1 and fails to properly clean up created resources. See
	// https://trello.com/c/4q0jprDy
	if err != nil && cmd.ProcessState.ExitCode() != expectedExitCode {
		return fmt.Errorf("waiting for command to return: %w", err)
	}

	return nil
}

func printUsage() {
	fmt.Println("forward-postgres-tester")
	fmt.Println()
	fmt.Printf("Runs okctl forward postgres and asserts the text \"%s\" printed", expectedString)
	fmt.Println()
	fmt.Println("USAGE:")
	fmt.Println("\tforward-postgres-tester <okctl binary path> <cluster manifest path>")
	fmt.Println("Example:")
	fmt.Println("\tforward-postgres-tester /usr/bin/okctl cluster.yaml")
	fmt.Println()
}
