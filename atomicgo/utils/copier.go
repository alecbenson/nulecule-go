package utils

import (
	"errors"
	"io"
	"os"
	"path/filepath"

	"github.com/Sirupsen/logrus"
)

//CopyFile copies a file from source to dest
func CopyFile(source, dest string) error {
	if !PathExists(source) {
		logrus.Errorf("Cannot copy file at %s: file does not exist", source)
		return errors.New("Invalid file")
	}

	sourcefile, err := os.Open(source)
	if err != nil {
		return err
	}
	defer sourcefile.Close()
	destfile, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer destfile.Close()

	//Copy the file from source to dest
	_, err = io.Copy(destfile, sourcefile)
	if err == nil {
		sourceinfo, err := os.Stat(source)
		if err != nil {
			err = os.Chmod(dest, sourceinfo.Mode())
		}
	}
	return nil
}

//CopyDirectory copies the contents of the directory at source to dest
func CopyDirectory(source, dest string) error {
	var err error
	if !PathExists(source) {
		logrus.Errorf("Cannot copy directory at %s: file does not exist", source)
		return errors.New("Invalid directory")
	}

	info, err := os.Stat(source)
	if err != nil {
		return err
	}

	//Create the destination directory
	err = os.MkdirAll(dest, info.Mode())
	if err != nil {
		return err
	}
	//Open the source and read all of its objects
	directory, _ := os.Open(source)
	defer directory.Close()
	children, err := directory.Readdir(-1)

	//Iterate over all children of the directory
	for _, child := range children {
		sourcePath := filepath.Join(source, child.Name())
		destPath := filepath.Join(dest, child.Name())

		//if the child is a directory...
		if child.IsDir() {
			//...Recurse over subdirectories
			err = CopyDirectory(sourcePath, destPath)
			if err != nil {
				logrus.Errorf("Failed to copy directory from %s to %s", sourcePath, destPath)
			}
			continue
		}

		//If the child is a file, copy it
		err = CopyFile(sourcePath, destPath)
		if err != nil {
			logrus.Errorf("Failed to copy directory from %s to %s", sourcePath, destPath)
		}
	}
	return nil
}
