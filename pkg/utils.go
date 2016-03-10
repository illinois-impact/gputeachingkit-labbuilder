package pandoc

import "os"

func isFileOk(file string) (bool, error) {
	sf, err := os.Open(file)
	if err != nil {
		return false, err
	}
	defer sf.Close()

	fi, err := sf.Stat()
	if err != nil {
		return false, err
	}

	if fi.IsDir() {
		return false, nil
	}

	return true, nil
}

func isFile(file string) bool {
	ok, err := isFileOk(file)
	if err != nil {
		return false
	}

	return ok
}

func isDir(file string) bool {
	ok, err := isFileOk(file)
	if err != nil {
		return false
	}
	return !ok
}
