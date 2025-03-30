package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
)

func CheckUpdates(currentVersion string) (string, error) {
	resp, err := http.Get(fmt.Sprintf("http://yourserver.com/version?current=%s", currentVersion))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNoContent {
		return "", nil // Нет обновлений
	}

	version, _ := io.ReadAll(resp.Body)
	return string(version), nil
}

func DownloadUpdate(version string) error {
	resp, err := http.Get(fmt.Sprintf("http://yourserver.com/download/%s", version))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	out, err := os.CreateTemp("", "update-*.zip")
	if err != nil {
		return err
	}
	defer os.Remove(out.Name())

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	cmd := exec.Command("unzip", "-o", out.Name(), "-d", "/usr/local/bin")
	return cmd.Run()
}
