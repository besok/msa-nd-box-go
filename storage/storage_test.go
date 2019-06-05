package storage

import (
	path2 "path"
	"testing"
)

func TestCreatStorageIfNotExists(t *testing.T){
	commonPath := "C:\\projects\\msa-nd-box-go\\file_storage"
	storage := "test_storage"
	if _, err := CreateStorage(commonPath, storage); err != nil {
		t.Fatalf("don't expected error : %s \n",err)
	}
}

func TestFromFile(t *testing.T){
	commonPath := "C:\\projects\\msa-nd-box-go\\file_storage"
	storage := "test_storage"

	s, err := CreateStorage(commonPath, storage)
	if err != nil {
		t.Fatalf("don't expected error : %s \n",err)
	}

	path := path2.Join(s.storagePath(),"1.txt")
	lines, err := fromFile(&s.mutex, path)
	if err != nil {
		t.Fatalf("don't expected error : %s \n",err)
	}

	for _,l := range lines{
		println(l)
	}

}
