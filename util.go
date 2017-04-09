package gopg

import "os"

func exists(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	return file.Close()
}
