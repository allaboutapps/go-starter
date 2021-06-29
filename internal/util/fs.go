package util

import (
	"os"
	"time"
)

// TouchFile creates an empty file if the file doesn’t already exist.
// If the file already exists then TouchFile updates the modified time of the file.
// Returns the modification time of the created / updated file.
func TouchFile(absolutePathToFile string) (time.Time, error) {
	_, err := os.Stat(absolutePathToFile)

	if os.IsNotExist(err) {
		file, err := os.Create(absolutePathToFile)

		if err != nil {
			return time.Time{}, err
		}

		defer file.Close()

		stat, err := file.Stat()

		if err != nil {
			return time.Time{}, err
		}

		return stat.ModTime(), nil
	}

	currentTime := time.Now().Local()
	err = os.Chtimes(absolutePathToFile, currentTime, currentTime)
	return currentTime, err
}
