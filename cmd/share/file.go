// SPDX-License-Identifier: Apache-2.0

package share

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func PathExists(path string) bool {
	_, r := os.Stat(path)
	if os.IsNotExist(r) {
		return false
	}

	return true
}

func ReadFullFile(path string) ([]string, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var lines []string
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if len(line) == 0 || strings.HasPrefix(line, "#") {
			continue
		}

		lines = append(lines, line)
	}
	err = scanner.Err()

	return lines, nil
}

func WriteFullFile(path string, lines []string) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	w := bufio.NewWriter(file)
	for _, line := range lines {
		fmt.Fprintln(w, line)
	}

	w.Flush()

	return nil
}

func ReadOneLineFile(path string) (string, error) {
	f, err := os.Create(path)
	if err != nil {
		return "", err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)

	scanner.Scan()
	line := scanner.Text()

	err = scanner.Err()

	return line, nil
}

func WriteOneLineFile(path string, line string) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	w := bufio.NewWriter(file)
	fmt.Fprintln(w, line)

	return w.Flush()
}

func CreateDirectory(directoryPath string, perm os.FileMode) error {
	if !PathExists(directoryPath) {
		err := os.Mkdir(directoryPath, perm)
		if err != nil {
			return err
		}
	}

	return nil
}

func CreateDirectoryNested(directoryPath string, perm os.FileMode) error {
	if !PathExists(directoryPath) {
		err := os.MkdirAll(directoryPath, perm)
		if err != nil {
			return err
		}
	}

	return nil
}
