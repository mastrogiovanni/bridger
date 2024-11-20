package executor

import (
	"bufio"
	"fmt"
	"log"
	"os/exec"
	"runtime"
	"strings"
)

func reader(scanner *bufio.Scanner, out chan string) {
	for scanner.Scan() {
		out <- scanner.Text()
	}
}

func ExecCmd(cmd string) string {
	out := ""
	err := error(nil)
	var outBytes []byte
	if runtime.GOOS == "windows" {
		// use command for windows
		outBytes, err = exec.Command("powershell", "-Command", cmd).Output()
	} else {
		// use command for linux/unix
		// outBytes, err = exec.Command("bash", "-c", cmd).Output()
		outBytes, err = exec.Command("sh", "-c", cmd).Output()
	}
	out = strings.TrimSpace(string(outBytes))
	if err != nil {
		errStr := fmt.Sprintf("Error running command: %s\n%s", cmd, err)
		log.Fatal(errStr)
	}
	return out
}

func ExecuteCommandAsync(cmd string, args []string, outputCh chan string, errorCh chan string, interruptCh chan struct{}) {
	defer close(outputCh)

	command := exec.Command(cmd, args...)
	stdout, err := command.StdoutPipe()
	if err != nil {
		outputCh <- fmt.Sprintf("Error creating stdout pipe: %v", err)
		errorCh <- fmt.Sprintf("Error creating stdout pipe: %v", err)
		return
	}

	stderr, err := command.StderrPipe()
	if err != nil {
		outputCh <- fmt.Sprintf("Error creating stderr pipe: %v", err)
		errorCh <- fmt.Sprintf("Error creating stderr pipe: %v", err)
		return
	}

	if err := command.Start(); err != nil {
		outputCh <- fmt.Sprintf("Error starting command: %v", err)
		errorCh <- fmt.Sprintf("Error starting command: %v", err)
		return
	}

	// Goroutine to handle process output
	scannerOutput := bufio.NewScanner(stdout)
	go func() {
		for scannerOutput.Scan() {
			outputCh <- scannerOutput.Text()
		}
	}()

	scannerError := bufio.NewScanner(stderr)
	go func() {
		for scannerError.Scan() {
			errorCh <- scannerError.Text()
		}
	}()

	// Wait for either interrupt signal or process completion
	done := make(chan error)
	go func() {
		done <- command.Wait()
	}()

	select {
	case <-interruptCh:
		if err := command.Process.Kill(); err != nil {
			outputCh <- fmt.Sprintf("Error killing process: %v", err)
			errorCh <- fmt.Sprintf("Error killing process: %v", err)
		} else {
			outputCh <- "Process interrupted."
			errorCh <- "Process interrupted."
		}
	case err := <-done:
		if err != nil {
			outputCh <- fmt.Sprintf("Process finished with error: %v", err)
			errorCh <- fmt.Sprintf("Process finished with error: %v", err)
		} else {
			outputCh <- "Process completed successfully."
			errorCh <- "Process completed successfully."
		}
	}
}
