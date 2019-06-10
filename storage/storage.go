package storage

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"
	"reflect"
	"sync"
)

type Storage struct {
	path       string
	name       string
	mutex      sync.Mutex
	memory     map[string]Lines
	createFunc func() Lines
	handler    ListenerHandler
}

func Snapshot(s *Storage) string {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	var b bytes.Buffer
	b.WriteString(fmt.Sprintf("storage[%s] snapshot:\n", s.name))
	for k, v := range s.memory {
		records := v.ToString()
		b.WriteString(fmt.Sprintf("key: %s\n", k))
		for i, l := range records {
			b.WriteString(fmt.Sprintf("value[%d]: %s\n", i, l))
		}
	}
	return b.String()
}

func CreateStorageOnly(p string, name string, createType func() Lines) (Storage, error) {
	return CreateStorage(p, name, createType, make([]Listener, 0))
}

func CreateStorage(p string, name string, createType func() Lines, listeners []Listener) (Storage, error) {
	log.Printf("init storage, path: %s, storage: %s, type: %s\n", p, name, reflect.TypeOf(createType()))
	handler := CreateListenerHandler()
	if len(listeners) > 0 {
		handler.listeners = listeners
	}
	storage := Storage{p,
		name, sync.Mutex{}, make(map[string]Lines),
		createType,
		handler}
	storage, err := storage.creatStorage()

	if err != nil {
		log.Fatalf(" error while creating path: %s , error: %s \n", name, err)
		return storage, err
	}

	err = storage.readAllFiles(createType)
	storage.handler.Handle(Init, StorageName(name), "", nil)
	return storage, err
}

func (s *Storage) Get(key string) (Lines, bool) {
	if lines, ok := s.memory[key]; ok {
		return lines, ok
	}
	s.handler.Handle(Get, StorageName(s.name), key, nil)
	return nil, false
}

func (s *Storage) Contains(key string) bool {
	_, ok := s.memory[key]
	return ok
}

func (s *Storage) AddListener(listener Listener) bool {
	s.handler.AddListener(listener)
	return true
}
func (s *Storage) Put(key string, line Line) error {
	s.mutex.Lock()
	defer func() {
		log.Printf("put value to a storage: %s, key: %s, value: %s \n", s.name, key, line)
		s.handler.Handle(Put, StorageName(s.name), key, line)
		s.mutex.Unlock()
	}()
	pathKey := s.storagePathKey(key)
	lines, ok := s.memory[key]
	if !ok {
		lines := s.createFunc()
		lines.Add(line)
		if e := createFile(pathKey); e != nil {
			log.Fatalf(" error while creating a file for key : %s", e)
		}
		s.memory[key] = lines
	} else {
		lines.Remove(line)
		lines.Add(line)
		s.memory[key] = lines
	}
	return rewriteFile(pathKey, s.memory[key])
}

func (s *Storage) RemoveKey(key string) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	pathKey := s.storagePathKey(key)
	delete(s.memory, key)
	log.Printf("remove key at a storage: %s, key: %s\n", s.name, key)
	s.handler.Handle(RemoveKey, StorageName(s.name), key, nil)
	return os.Remove(pathKey)
}

func (s *Storage) RemoveValue(key string, line Line) error {
	s.mutex.Lock()
	needUnlock := true
	defer func() {
		if needUnlock {
			s.mutex.Unlock()
		}

		s.handler.Handle(RemoveVal, StorageName(s.name), key, line)
		log.Printf("remove value at a storage: %s, key: %s, value: %s \n", s.name, key, line)
	}()

	lines, ok := s.memory[key]
	if !ok {
		return nil
	}
	ok = lines.Remove(line)
	if !ok {
		return nil
	}

	if lines.Size() == 0 {
		s.mutex.Unlock()
		needUnlock = false
		return s.RemoveKey(key)
	}

	s.memory[key] = lines
	return rewriteFile(s.storagePathKey(key), lines)
}
func (s *Storage) Clean() error {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.memory = make(map[string]Lines)
	strPath := s.storagePath()
	err := os.RemoveAll(strPath)
	err = createDir(strPath)
	log.Printf("clean storage: %s \n", s.name)
	s.handler.Handle(Clean, StorageName(s.name), "", nil)
	return err
}

func (s *Storage) readAllFiles(createType func() Lines) error {
	p := s.storagePath()
	err := filepath.Walk(p, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			records, err := readFile(path)
			if err != nil {
				log.Fatalf(" error while pasing path: %s , error: %s \n", path, err)
			}
			lines := createType()
			lines.fromString(records)
			s.memory[info.Name()] = lines
		}
		return err
	})

	if err != nil {
		log.Fatalf(" error while pasing path: %s , error: %s \n", s.storagePath(), err)
	}
	return err
}
func (s *Storage) storagePath() string {
	return path.Join(s.path, s.name)
}
func (s *Storage) storagePathKey(key string) string {
	return path.Join(s.path, s.name, key)
}
func (s *Storage) creatStorage() (Storage, error) {
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
		log.Printf(" dir created  path: %s \n", path)
	}
	return nil
}
func readFile(p string) (Records, error) {
	lines := make([]string, 0)
	file, e := os.Open(p)
	if e != nil {
		log.Fatalf(" error while reading file: %s , error: %s \n", p, e)
		return nil, e
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if e = scanner.Err(); e != nil {
		log.Fatalf(" error while reading file: %s , error: %s \n", p, e)
		return nil, e
	}

	return lines, e
}
func createFile(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		f, err := os.Create(path)
		if err != nil {
			return err
		}
		defer f.Close()
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

	records := lines.ToString()
	wr := bufio.NewWriter(file)
	for _, r := range records {
		if _, err := fmt.Fprintln(wr, r); err != nil {
			log.Fatalf("error while writting file: %s", err)
		}
	}
	err = wr.Flush()
	err = file.Sync()
	return err
}
