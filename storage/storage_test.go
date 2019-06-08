package storage

import (
	"fmt"
	path2 "path"
	"testing"
)

func TestCreatStorageIfNotExists(t *testing.T) {
	commonPath := "C:\\projects\\msa-nd-box-go\\file_storages"
	storage := "test_storage"
	if _, err := CreateStorage(commonPath, storage, CreateStrLines); err != nil {
		t.Fatalf("don't expected error : %s \n", err)
	}
}

func TestFromFile(t *testing.T) {
	commonPath := "C:\\projects\\msa-nd-box-go\\file_storage"
	storage := "test_storage"

	s, err := CreateStorage(commonPath, storage, CreateStrLines)
	if err != nil {
		t.Fatalf("don't expected error : %s \n", err)
	}

	path := path2.Join(s.storagePath(), "1.txt")
	lines, err := readRawFromFile(path)
	if err != nil {
		t.Fatalf("don't expected error : %s \n", err)
	}

	for _, l := range lines {
		println(l)
	}

}

func TestInterface(t *testing.T) {
	commonPath := "C:\\projects\\msa-nd-box-go\\file_storages"
	storage := "test_storage"

	s, err := CreateStorage(commonPath, storage, CreateStrLines)
	if err != nil {
		t.Fatalf("don't expected error : %s \n", err)
	}

	path := path2.Join(s.storagePath(), "1.txt")
	lines, err := readRawFromFile(path)

	strLines := new(StrLines)

	strLines.fromString(lines)

	PrintLines(strLines)
	if len(strLines.lines) == 0 {
		t.Fatalf("should be values from file")
	}
}
func TestReadAllFiles(t *testing.T) {
	commonPath := "C:\\projects\\msa-nd-box-go\\file_storages"
	storage := "test_storage"

	s, err := CreateStorage(commonPath, storage, CreateStrLines)
	if err != nil {
		t.Fatalf("don't expected error : %s \n", err)
	}

	for k, v := range s.memory {
		fmt.Printf("service : %s\n", k)
		PrintLines(v)
	}
}
func TestStorage_Put(t *testing.T) {
	commonPath := "C:\\projects\\msa-nd-box-go\\file_storages"
	storage := "test_storage"

	s, err := CreateStorage(commonPath, storage, CreateStrLines)
	if err != nil {
		t.Fatalf("don't expected error : %s \n", err)
	}
	err = s.Put("n", StrLine{"test3"})
	err = s.Put("n", StrLine{"test1"})
	err = s.Put("n", StrLine{"test2"})
	if  err != nil {
		t.Fatalf("don't expected error : %s \n", err)
	}

	for k, v := range s.memory {
		fmt.Printf("service : %s\n", k)
		PrintLines(v)
	}
}
