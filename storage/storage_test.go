package storage

import (
	"fmt"
	"log"
	"os"
	"testing"
)

func TestCreatStorageIfNotExists(t *testing.T) {
	commonPath := "C:\\projects\\msa-nd-box-go\\file_storages"
	storage := "test_storage"
	if _, err := CreateStorage(commonPath, storage, CreateStringLines); err != nil {
		t.Fatalf("don't expected error : %s \n", err)
	}
}

func TestStorage_Put(t *testing.T) {
	commonPath := "C:\\projects\\msa-nd-box-go\\file_storages"
	storage := "test_storage"

	s, err := CreateStorage(commonPath, storage, CreateStringLines)
	if err != nil {
		t.Fatalf("don't expected error : %s \n", err)
	}
	testKey := "test_key"
	err = s.Put(testKey, StringLine{"test3"})
	err = s.Put(testKey, StringLine{"test1"})
	err = s.Put(testKey, StringLine{"test2"})
	if err != nil {
		t.Fatalf("don't expected error : %s \n", err)
	}

	for k, v := range s.memory {
		fmt.Printf("service : %s\n", k)
		PrintLines(v)
	}
}
func TestStorage_RemoveKey(t *testing.T) {
	commonPath := "C:\\projects\\msa-nd-box-go\\file_storages"
	storage := "test_storage"

	s, err := CreateStorage(commonPath, storage, CreateStringLines)
	if err != nil {
		t.Fatalf("don't expected error : %s \n", err)
	}
	testKey := "test_key"
	err = s.Put(testKey, StringLine{"test"})
	if err != nil {
		t.Fatalf("don't expected error : %s \n", err)
	}
	err = s.RemoveKey(testKey)
	if err != nil {
		t.Fatalf("don't expected error : %s \n", err)
	}
	if _, ok := s.memory[testKey]; ok {
		t.Fatalf("don't expected error : %s \n", err)
	}
	if _, err := os.Stat(s.storagePathKey(testKey)); !os.IsNotExist(err) {
		t.Fatalf("don't expected error : %s \n", err)
	}
}

func TestStorage_RemoveValue(t *testing.T) {
	commonPath := "C:\\projects\\msa-nd-box-go\\file_storages"
	storage := "test_storage"

	s, err := CreateStorage(commonPath, storage, CreateStringLines)
	if err != nil {
		t.Fatalf("don't expected error : %s \n", err)
	}
	testKey := "test_key"
	testVal := "test_Val"
	err = s.Put(testKey, StringLine{testVal})
	err = s.Put(testKey, StringLine{"test1"})
	err = s.Put(testKey, StringLine{"test2"})
	if err != nil {
		t.Fatalf("don't expected error : %s \n", err)
	}

	err = s.RemoveValue(testKey, StringLine{testVal})

	if err != nil {
		t.Fatalf("don't expected error : %s \n", err)
	}

	PrintLines(s.memory[testKey])

}
func TestStorage_RemoveValueAndFile(t *testing.T) {
	commonPath := "C:\\projects\\msa-nd-box-go\\file_storages"
	storage := "test_storage"

	s, err := CreateStorage(commonPath, storage, CreateStringLines)
	if err != nil {
		t.Fatalf("don't expected error : %s \n", err)
	}
	testKey := "test_key"
	testVal := "test_Val"
	err = s.Put(testKey, StringLine{testVal})
	if err != nil {
		t.Fatalf("don't expected error : %s \n", err)
	}

	err = s.RemoveValue(testKey, StringLine{testVal})

	if err != nil {
		t.Fatalf("don't expected error : %s \n", err)
	}

	PrintLines(s.memory[testKey])
}
func TestStorage_Clean(t *testing.T) {
	commonPath := "C:\\projects\\msa-nd-box-go\\file_storages"
	storage := "test_storage"

	s, err := CreateStorage(commonPath, storage, CreateStringLines)
	if err != nil {
		t.Fatalf("don't expected error : %s \n", err)
	}
	testKey := "test_key"
	testVal := "test_Val"
	err = s.Put(testKey, StringLine{testVal})
	err = s.Put(testKey, StringLine{"test2"})
	if err != nil {
		t.Fatalf("don't expected error : %s \n", err)
	}

	err = s.Clean()

	if err != nil {
		t.Fatalf("don't expected error : %s \n", err)
	}
}
func TestStorage_Handler(t *testing.T) {
	commonPath := "C:\\projects\\msa-nd-box-go\\file_storages"
	str := "test_storage"

	s, err := CreateStorage(commonPath, str, CreateStringLines)
	if err != nil {
		t.Fatalf("don't expected error : %s \n", err)
	}
	testKey := "test_key"
	testVal := "test_Val"
	s.AddListener(func(event StorageEvent, storage StorageName, key string, value Line) {
		if event == Put && StorageName(str) == storage{
			log.Printf(" listener : key: %s, value: %s",key,value)

			line := StringLine{testVal}
			if key != testKey && value != line {
				t.Fatalf("don't expected error : %s \n", err)
			}
		}
	})
	err = s.Put(testKey, StringLine{testVal})
	if err != nil {
		t.Fatalf("don't expected error : %s \n", err)
	}

	err = s.Clean()

	if err != nil {
		t.Fatalf("don't expected error : %s \n", err)
	}
}
