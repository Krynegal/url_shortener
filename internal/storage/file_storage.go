package storage

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
)

type URLObj struct {
	Key string `json:"key"`
	URL string `json:"url"`
}

type FileStorage struct {
	storagePath string
}

func NewFileStorage(filePath string) (*FileStorage, error) {
	fs := &FileStorage{
		storagePath: filePath,
	}

	return fs, nil
}

func (fs *FileStorage) ReadURLsFromFile(memStorage *MemStorage) error {
	file, err := os.OpenFile(fs.storagePath, os.O_RDONLY|os.O_CREATE, 0777)
	if err != nil {
		return err
	}

	defer func() {
		cerr := file.Close()
		if cerr != nil {
			err = cerr
		}
	}()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		if scanner.Err() != nil {
			return scanner.Err()
		}
		line := scanner.Text()
		u := URLObj{}
		err := json.Unmarshal([]byte(line), &u)
		if err != nil {
			break
		}
		memStorage.store[u.Key] = u.URL
		memStorage.counter++
		fmt.Printf("URLObj: %v\n", u)
	}
	return nil
}

func (fs *FileStorage) WriteURLInFile(key, url string) error {
	file, err := os.OpenFile(fs.storagePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0777)
	if err != nil {
		return err
	}

	defer func() {
		cerr := file.Close()
		if cerr != nil {
			err = cerr
		}
	}()

	u := URLObj{Key: key, URL: url}
	b, err := json.Marshal(u)
	if err != nil {
		return err
	}

	writer := bufio.NewWriter(file)
	defer func() {
		_ = writer.Flush()
	}()

	if _, err = writer.WriteString(string(b) + "\n"); err != nil {
		return err
	}

	return nil
}
