package main

import (
	"bytes"
	"os/exec"
	"strconv"
	"strings"
)

type StatusEntry struct {
	Program string
	Status  string
	Info    string
}

func NewStatusEntryFromString(row string) *StatusEntry {
	entry := &StatusEntry{}

	var stage uint8
	for _, char := range []rune(row) {
		if stage == 0 {
			if char == ' ' {
				stage = 1
			} else {
				entry.Program += string(char)
			}
		} else if stage == 1 {
			if char != ' ' {
				entry.Status = string(char)
				stage = 2
			}
		} else if stage == 2 {
			if char == ' ' {
				stage = 3
			} else {
				entry.Status += string(char)
			}
		} else if stage == 3 {
			if char != ' ' {
				entry.Info = string(char)
				stage = 4
			}
		} else {
			entry.Info += string(char)
		}
	}

	return entry
}

func getStatusAll() ([]*StatusEntry, error) {
	entries := []*StatusEntry{}

	stdout, _, err := getStatusRaw("all")
	if err != nil {
		return entries, err
	}

	rows := strings.Split(stdout, "\n")
	for _, row := range rows {
		if row == "" {
			continue
		}
		entries = append(entries, NewStatusEntryFromString(row))
	}
	return entries, nil
}

func getStatusRaw(program string) (string, string, error) {
	cmd := exec.Command("supervisorctl", "status", program)
	var stdout bytes.Buffer
	cmd.Stdout = &stdout
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	err := cmd.Run()

	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			if exitErr.ExitCode() != 3 {
				return stdout.String(), stderr.String(), exitErr
			}
		} else {
			return stdout.String(), stderr.String(), err
		}
	}

	return stdout.String(), stderr.String(), nil
}

func getTailRaw(program string, pipe string, numBytes uint64) (string, string, error) {
	numBytesStr := strconv.FormatUint(numBytes, 10)

	cmd := exec.Command("supervisorctl", "tail", "-"+numBytesStr, program, pipe)
	var stdout bytes.Buffer
	cmd.Stdout = &stdout
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	err := cmd.Run()

	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			if exitErr.ExitCode() != 3 {
				return stdout.String(), stderr.String(), exitErr
			}
		} else {
			return stdout.String(), stderr.String(), err
		}
	}

	return stdout.String(), stderr.String(), nil
}

func manageProcess(program string, action string) (string, string, error) {
	cmd := exec.Command("supervisorctl", action, program)
	var stdout bytes.Buffer
	cmd.Stdout = &stdout
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	err := cmd.Run()

	return stdout.String(), stderr.String(), err
}
