package storage

import (
	"bufio"
	"fmt"
	"os"
	"path"
	"sync"
)

type Lines []Line
type Line interface{}

type Storage struct {
	path   string
	name   string
	mutex  sync.Mutex
	memory map[string]Lines
}

func CreateStorage(p string, name string) (Storage, error) {
	storage, err := Storage{
		p, name, sync.Mutex{}, make(map[string]Lines),
	}.creatStorageIfNotExists()

	if err != nil {
		fmt.Printf(" error while creating path: %s , error: %s \n", name, err)
	}

	return storage,err
}

func (s Storage) storagePath() string {
	return path.Join(s.path, s.name)
}

func fromFile(m *sync.Mutex, p string) ([]string, error) {
	m.Lock()
	lines := make([]string, 0)
	file, e := os.Open(p)
	if e != nil {
		fmt.Printf(" error while reading file: %s , error: %s \n", p, e)
		return nil, e
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if e = scanner.Err(); e != nil {
		fmt.Printf(" error while reading file: %s , error: %s \n", p, e)
		return nil, e
	}

	m.Unlock()
	return lines, e
}

func (s Storage) creatStorageIfNotExists() (Storage, error) {
	err := createDir(s.path)
	err = createDir(s.storagePath())
	return s, err
}

func createDir(path string) error {
	if _, e := os.Stat(path); os.IsNotExist(e) {
		if err := os.MkdirAll(path, os.ModePerm); err != nil {
			fmt.Printf(" error while creating path: %s , error: %s \n", path, err)
			return err
		}
		fmt.Printf(" dir created  path: %s \n", path)
	}
	return nil
}
