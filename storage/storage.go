package storage

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"
	"sync"
)

type Storage struct {
	path       string
	name       string
	mutex      sync.Mutex
	memory     map[string]Lines
	createFunc func() Lines
}

func CreateStorage(p string, name string, createType func() Lines) (Storage, error) {
	storage := Storage{p, name, sync.Mutex{}, make(map[string]Lines), createType}
	storage, err := storage.creatStorageIfNotExists()

	if err != nil {
		fmt.Printf(" error while creating path: %s , error: %s \n", name, err)
		return storage, err
	}

	err = storage.readAllFiles(createType)
	return storage, err
}

func (s *Storage) Get(key string) (*Lines, bool) {
	if lines, ok := s.memory[key]; ok {
		return &lines, ok
	}
	return nil, false
}

func (s *Storage) Contains(key string) bool {
	_, ok := s.memory[key]
	return ok
}

func (s *Storage) Put(key string, line Line) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	pathKey := s.storagePathKey(key)

	lines, ok := s.memory[key]
	if !ok {
		lines := s.createFunc()
		lines.Add(line)
		if e := createFile(pathKey); e != nil {
			fmt.Printf(" error while creating a file for key : %s", e)
		}
		s.memory[key] = lines
	} else {
		lines.Remove(line)
		lines.Add(line)
		s.memory[key] = lines
	}
	return rewriteFile(pathKey, s.memory[key])
}

func (s *Storage) readAllFiles(createType func() Lines) error {
	p := s.storagePath()
	err := filepath.Walk(p, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			records, err := readRawFromFile(path)
			if err != nil {
				fmt.Printf(" error while pasing path: %s , error: %s \n", path, err)
			}
			lines := createType()
			lines.fromString(records)
			s.memory[info.Name()] = lines
		}
		return err
	})

	if err != nil {
		fmt.Printf(" error while pasing path: %s , error: %s \n", s.storagePath(), err)
	}
	return err
}
func (s *Storage) storagePath() string {
	return path.Join(s.path, s.name)
}
func (s *Storage) storagePathKey(key string) string {
	return path.Join(s.path, s.name, key)
}

func (s *Storage) creatStorageIfNotExists() (Storage, error) {
	err := createDir(s.path)
	err = createDir(s.storagePath())
	return *s, err
}
func createDir(path string) error {
	if _, e := os.Stat(path); os.IsNotExist(e) {
		if err := os.MkdirAll(path, os.ModePerm); err != nil {
			log.Fatalf(" error while creating path: %s , error: %s \n", path, err)
			return err
		}
		log.Fatalf(" dir created  path: %s \n", path)
	}
	return nil
}
func readRawFromFile(p string) (Records, error) {
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

	return lines, e
}
func createFile(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		_, err := os.Create(path)
		return err
	}
	return nil
}
func rewriteFile(path string, lines Lines) error {
	err := os.Remove(path)
	err = createFile(path)
	if err != nil {
		return err
	}
	file, err := os.OpenFile(path, os.O_WRONLY, 0600)
	if err != nil {
		return err
	}

	defer file.Close()

	if lines == nil {
		return nil
	}

	records := lines.ToString()
	wr := bufio.NewWriter(file)
	for _, r := range records {
		if _, err := fmt.Fprintln(wr, r); err != nil {
			fmt.Printf("error while writting file: %s", err)
		}
	}
	return wr.Flush()
}
