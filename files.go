package main

import (
	"fmt"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"slices"

	"github.com/h2non/bimg"
)

func getInputFiles(inputPath string, extensions []string, recursive bool) ([]string, error) {
	file, err := os.Stat(inputPath)
	if err != nil {
		return nil, err
	}

	if !file.IsDir() {
		if slices.Contains(extensions, filepath.Ext(inputPath)) {
			return []string{inputPath}, nil
		}
	}
	return walkImagePath(inputPath, extensions, recursive)
}

// This function assumes that the directory status of the path has already been checked
func walkImagePath(dir string, extensions []string, recursive bool) ([]string, error) {
	if !recursive {
		return getImagesInPath(dir, extensions)
	}
	return walkPathForImages(dir, extensions)
}

func getImagesInPath(dir string, extensions []string) ([]string, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	paths := []string{}
	for _, entry := range entries {
		paths = append(paths, path.Join(dir, entry.Name()))
	}
	return paths, nil

}

func walkPathForImages(dir string, extensions []string) ([]string, error) {
	entries := []string{}
	err := filepath.WalkDir(dir, func(path string, info fs.DirEntry, err error) (e error) {
		if !info.IsDir() && slices.Contains(extensions, filepath.Ext(info.Name())) {
			entries = append(entries, path)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return entries, nil
}

func getModifiedImageName(original string, outputExt string, image bimg.Image, attempt int) (string, error) {
	size, err := image.Size()
	if err != nil {
		return "", err
	}
	width := size.Width

	name := original + fmt.Sprint(width) + "w" + outputExt

	return name, nil
}
