package downloader

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
)

func SaveFile(link string, fileName string) error {
	response, err := http.Get(link)

	if err != nil {
		return err
	}

	defer response.Body.Close()

	homeDir, err := os.UserHomeDir()

	if err != nil {
		return err
	}

	file, err := os.Create(path.Join(homeDir, "Downloads", fmt.Sprintf("%s.torrent", fileName)))

	if err != nil {
		return err
	}

	defer file.Close()

	_, err = io.Copy(file, response.Body)

	if err != nil {
		return err
	}

	return nil
}
