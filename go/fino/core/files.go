package core

import "os"

// CreatDir create dir -r
func CreatDir(path string) error {
	if !IsExist(path) {
		return os.MkdirAll(path, os.ModePerm)
	}
	return nil
}

// IsExist check path is exist, return true if exist
func IsExist(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		return os.IsExist(err)
	}
	return true
}

func IsExecutable(mode os.FileMode) bool {
	return mode&0111 != 0
}

func IsOwnerExecutable(mode os.FileMode) bool {
	return mode&0100 != 0
}

func IsGroupExecutable(mode os.FileMode) bool {
	return mode&0010 != 0
}

func IsOtherExecutable(mode os.FileMode) bool {
	return mode&0001 != 0
}

func IsAllExecutables(mode os.FileMode) bool {
	return mode&0111 == 0111
}
