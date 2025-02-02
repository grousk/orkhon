// Package file provides basic file utility functions
package file

import (
	"io"
	"os"
)

// Read reads the entire content of a file and returns it as a byte slice.
func Read(path string) ([]byte, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// TODO: maybe use io.Copy instead
	return io.ReadAll(file)
}

// Rename moves a file to a new destination.
func Rename(src, dst string) error {
	err := os.Rename(src, dst)
	if err != nil {
		return err
	}
	return nil
}

// Copy copies a file from src to dst.
func Copy(src, dst string) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	dstFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	_, err = io.Copy(dstFile, srcFile)
	return err
}
