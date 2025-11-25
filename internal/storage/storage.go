// Package storage
package storage

import (
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

var ErrInvalidPath = errors.New("invalid filepath")

type Storage interface {
	Save(file io.Reader, path string) error
	Read(path string) ([]byte, error)
	Open(path string) (io.ReadSeeker, fs.FileInfo, error)
	Delete(path string) error
	List() ([]FileList, error)
}

type FileList struct {
	Path string `json:"path"`
	Name string `json:"name"`
	Size int64  `json:"size"`
}

type LocalStorage struct {
	BasePath string
}

func NewLocalStorage(basepath string) *LocalStorage {
	return &LocalStorage{BasePath: filepath.Clean(basepath)}
}

func (l *LocalStorage) WithBucket(bucket string) *LocalStorage {
	clean := filepath.Clean(bucket)
	if clean == "." || clean == "" {
		clean = ""
	}

	return &LocalStorage{
		BasePath: filepath.Join(l.BasePath, clean),
	}
}

func (l *LocalStorage) Save(file io.Reader, path string) error {
	fullPath, err := l.safePath(path)
	if err != nil {
		return err
	}

	if err := os.MkdirAll(filepath.Dir(fullPath), 0o755); err != nil {
		return fmt.Errorf("create dirs: %w", err)
	}

	dst, err := os.Create(fullPath)
	if err != nil {
		return fmt.Errorf("create file: %w", err)
	}

	defer dst.Close()

	if _, err := io.Copy(dst, file); err != nil {
		return fmt.Errorf("write file: %w", err)
	}

	return nil
}

func (l *LocalStorage) Read(path string) ([]byte, error) {
	fullPath, err := l.safePath(path)
	if err != nil {
		return nil, err
	}

	content, err := os.ReadFile(fullPath)
	if err != nil {
		return nil, fmt.Errorf("read file: %w", err)
	}

	return content, nil
}

func (l *LocalStorage) Open(path string) (io.ReadSeeker, fs.FileInfo, error) {
	fullPath, err := l.safePath(path)
	if err != nil {
		return nil, nil, err
	}

	file, err := os.Open(fullPath)
	if err != nil {
		return nil, nil, err
	}

	stat, err := file.Stat()
	if err != nil {
		return nil, nil, err
	}

	return file, stat, nil
}

func (l *LocalStorage) Delete(path string) error {
	fullPath, err := l.safePath(path)
	if err != nil {
		return err
	}

	if err := os.Remove(fullPath); err != nil {
		return fmt.Errorf("delete file: %w", err)
	}

	return nil
}

func (l *LocalStorage) List() ([]FileList, error) {
	var files []FileList

	err := filepath.Walk(l.BasePath, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info != nil && !info.IsDir() {
			files = append(files, FileList{
				Path: path,
				Name: info.Name(),
				Size: info.Size(),
			})
		}
		return nil
	})
	if err != nil {
		return []FileList{}, err
	}

	return files, nil
}

func (l *LocalStorage) safePath(path string) (string, error) {
	cleanName := filepath.Clean(path)

	if cleanName == "." || cleanName == "/" || cleanName == "" {
		return "", fmt.Errorf("%w", ErrInvalidPath)
	}

	if filepath.IsAbs(cleanName) || strings.HasPrefix(cleanName, "..") {
		return "", fmt.Errorf("%w", ErrInvalidPath)
	}

	fullPath := filepath.Join(l.BasePath, cleanName)
	base := l.BasePath

	if base != "." {
		base = fmt.Sprintf("%s%s", filepath.Clean(base), string(os.PathSeparator))
	}

	if !strings.HasPrefix(fullPath, base) && filepath.Clean(fullPath) != filepath.Clean(l.BasePath) {
		return "", fmt.Errorf("%w", ErrInvalidPath)
	}

	return fullPath, nil
}
