package storage

import (
	"fmt"
	"log"
	"os"
	"testing"
)

func TestCreatStorageIfNotExists(t *testing.T) {
	commonPath := "C:\\projects\\msa-nd-box-go\\file_storages_test"
	storage := "test_storage"
	if _, err := CreateStorageOnly(commonPath, storage, CreateStringLines); err != nil {
		t.Fatalf("don't expected error : %s \n", err)
	}
}

func TestStorage_Put(t *testing.T) {
	commonPath := "C:\\projects\\msa-nd-box-go\\file_storages_test"
	storage := "test_storage"

	s, err := CreateStorageOnly(commonPath, storage, CreateStringLines)
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
		fmt.Printf("Service : %s\n", k)
		PrintLines(v)
	}
}
func TestStorage_PutDubl(t *testing.T) {
	commonPath := "C:\\projects\\msa-nd-box-go\\file_storages_test"
	storage := "test_storage_dub"

	s, err := CreateStorageOnly(commonPath, storage, CreateStringLines)
	if err != nil {
		t.Fatalf("don't expected error : %s \n", err)
	}
	testKey := "test_key"
	err = s.Put(testKey, StringLine{"test1"})
	err = s.Put(testKey, StringLine{"test1"})
	err = s.Put(testKey, StringLine{"test1"})
	if err != nil {
		t.Fatalf("don't expected error : %s \n", err)
	}

	if s.memory[testKey].Size() != 1 {
		t.Fatalf("don't expected error : %s \n", err)
	}
}

func TestStorage_Put_CB(t *testing.T) {
	commonPath := "C:\\projects\\msa-nd-box-go\\file_storages_test"
	storage := "test_storage_cb"

	s, err := CreateStorageOnly(commonPath, storage, CreateCBLines)
	if err != nil {
		t.Fatalf("don't expected error : %s \n", err)
	}
	testKey := "test_key"
	err = s.Put(testKey, CBLine{Address: "1", Active: true})
	err = s.Put(testKey, CBLine{Address: "2", Active: true})
	err = s.Put(testKey, CBLine{Address: "3", Active: true})
	if err != nil {
		t.Fatalf("don't expected error : %s \n", err)
	}

	for k, v := range s.memory {
		fmt.Printf("Service : %s\n", k)
		PrintLines(v)
	}
}

func TestStorage_RemoveKey(t *testing.T) {
	commonPath := "C:\\projects\\msa-nd-box-go\\file_storages_test"
	storage := "test_storage"

	s, err := CreateStorageOnly(commonPath, storage, CreateStringLines)
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
func TestStorage_RemoveKey_CB(t *testing.T) {
	commonPath := "C:\\projects\\msa-nd-box-go\\file_storages_test"
	storage := "test_storage_cb"

	s, err := CreateStorageOnly(commonPath, storage, CreateCBLines)
	if err != nil {
		t.Fatalf("don't expected error : %s \n", err)
	}
	testKey := "test_key"
	err = s.Put(testKey, CBLine{Address: "2", Active: true})
	err = s.Put(testKey, CBLine{Address: "3", Active: true})

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
	commonPath := "C:\\projects\\msa-nd-box-go\\file_storages_test"
	storage := "test_storage"

	s, err := CreateStorageOnly(commonPath, storage, CreateStringLines)
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
func TestStorage_RemoveValue_CB(t *testing.T) {
	commonPath := "C:\\projects\\msa-nd-box-go\\file_storages_test"
	storage := "test_storage_cb"

	s, err := CreateStorageOnly(commonPath, storage, CreateCBLines)
	if err != nil {
		t.Fatalf("don't expected error : %s \n", err)
	}
	testKey := "test_key"
	testVal := "test_Val"
	err = s.Put(testKey, CBLine{Address: "2", Active: true})
	err = s.Put(testKey, CBLine{Address: "3", Active: true})
	err = s.Put(testKey, CBLine{Address: testVal, Active: true})
	if err != nil {
		t.Fatalf("don't expected error : %s \n", err)
	}

	err = s.RemoveValue(testKey, CBLine{Address: testVal, Active: false})

	if err != nil {
		t.Fatalf("don't expected error : %s \n", err)
	}

	PrintLines(s.memory[testKey])

}
func TestStorage_RemoveValueAndFile(t *testing.T) {
	commonPath := "C:\\projects\\msa-nd-box-go\\file_storages_test"
	storage := "test_storage"

	s, err := CreateStorageOnly(commonPath, storage, CreateStringLines)
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
func TestStorage_RemoveValueAndFile_CB(t *testing.T) {
	commonPath := "C:\\projects\\msa-nd-box-go\\file_storages_test"
	storage := "test_storage_cb"

	s, err := CreateStorageOnly(commonPath, storage, CreateCBLines)
	if err != nil {
		t.Fatalf("don't expected error : %s \n", err)
	}
	testKey := "test_key"
	testVal := "test_Val"
	err = s.Put(testKey, CBLine{Address: testVal, Active: true})
	if err != nil {
		t.Fatalf("don't expected error : %s \n", err)
	}

	err = s.RemoveValue(testKey, CBLine{Address: testVal, Active: false})

	if err != nil {
		t.Fatalf("don't expected error : %s \n", err)
	}

	PrintLines(s.memory[testKey])
}
func TestStorage_Clean(t *testing.T) {
	commonPath := "C:\\projects\\msa-nd-box-go\\file_storages_test"
	storage := "test_storage"

	s, err := CreateStorageOnly(commonPath, storage, CreateStringLines)
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
	commonPath := "C:\\projects\\msa-nd-box-go\\file_storages_test"
	str := "test_storage"

	s, err := CreateStorageOnly(commonPath, str, CreateStringLines)
	if err != nil {
		t.Fatalf("don't expected error : %s \n", err)
	}
	testKey := "test_key"
	testVal := "test_Val"
	s.AddListener(func(event Event, storage Name, key string, value Line) {
		if event == Put && Name(str) == storage {
			log.Printf(" listener : key: %s, value: %s", key, value)

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
