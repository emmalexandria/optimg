package main

import (
	"errors"
	"fmt"
	"github.com/h2non/bimg"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"slices"
)

func getSubDirs(dir string, outputDir string) ([]string, error) {
	entries := []string{}
	err := filepath.WalkDir(dir, func(path string, info fs.DirEntry, err error) (e error) {
		if info.IsDir() && info.Name() != outputDir {
			entries = append(entries, path)
		}
		return nil
	})
	if err != nil {
		return nil, nil
	}

	return entries, nil
}

func getFilesInDir(dir string, extensions []string) ([]string, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	paths := []string{}
	for _, entry := range entries {
		if checkForFileExtensions(entry.Name(), extensions) && !entry.IsDir() {
			paths = append(paths, path.Join(dir, entry.Name()))
		}
	}

	return paths, nil
}

func checkForFileExtensions(path string, extensions []string) bool {
	pathExt := filepath.Ext(path)
	return slices.Contains(extensions, pathExt)
}

func saveImagesInOutput(images []bimg.Image, originalPath string, ext string, outputDir string) error {
	baseDir := path.Join(path.Dir(originalPath), outputDir)
	baseName := path.Base(originalPath)

	for _, img := range images {
		name, err := getNonConflictingName(img, baseName, baseDir, ext)
		if err != nil {
			return err
		}

		err = bimg.Write(path.Join(baseDir, name), img.Image())
		if err != nil {
			return err
		}
	}

	return nil
}

func getNonConflictingName(image bimg.Image, baseName string, baseDir string, ext string) (string, error) {
	genName, err := getImageName(baseName, ext, image, 0)
	if err != nil {
		return "", err
	}
	attempt := 0
	exists, err := doesFileExist(path.Join(baseDir, genName))

	for exists {
		attempt += 1
		genName, err = getImageName(baseName, ext, image, attempt)
		if err != nil {
			return "", err
		}
		exists, err = doesFileExist(path.Join(baseDir, genName))
		if err != nil {
			return "", err
		}
	}

	return genName, nil
}

func getImageName(original string, ext string, image bimg.Image, attempt int) (string, error) {
	size, err := image.Size()
	baseName := stripExtension(original)
	if err != nil {
		return "", err
	}
	width := size.Width
	attemptStr := ""
	if attempt > 0 {
		attemptStr = fmt.Sprint(attempt)
	}
	name := baseName + fmt.Sprint(width) + "w" + attemptStr + ext

	return name, nil
}

func stripExtension(baseName string) string {
	ext := path.Ext(baseName)
	return baseName[0 : len(baseName)-len(ext)]
}

func checkOutputClear(baseDir string, alreadyProcessed []string) bool {
	for _, img := range alreadyProcessed {
		//If we've already processed an image in this directory, the output directory
		//is from the current run and shouldn't be cleared
		if path.Dir(img) == baseDir {
			return false
		}
	}

	return true
}

func initOutputDir(baseDir string, outputDir string, clearOutput bool, alreadyProcessed []string) error {
	fullOutput := path.Join(baseDir, outputDir)
	_, err := os.Stat(fullOutput)
	exists := true
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			exists = false
		} else {
			return err
		}
	}

	if !clearOutput && exists {
		return nil
	}

	baseStat, err := os.Stat(baseDir)
	if err != nil {
		return err
	}
	if exists && clearOutput && checkOutputClear(baseDir, alreadyProcessed) {
		err = os.RemoveAll(fullOutput)
		os.Mkdir(fullOutput, baseStat.Mode())
	} else if !exists {
		os.Mkdir(fullOutput, baseStat.Mode())
	}

	return nil

}

func doesFileExist(path string) (bool, error) {
	stat, err := os.Stat(path)

	if err == nil {
		return !stat.IsDir(), nil
	} else {
		if errors.Is(err, os.ErrNotExist) {
			return false, nil
		}
		return false, err
	}
}

func isDir(path string) (bool, error) {
	stat, err := os.Stat(path)
	if err != nil {
		return false, err
	}
	return stat.IsDir(), nil
}
