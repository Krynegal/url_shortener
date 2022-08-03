package storage

import (
	"bufio"
	"encoding/json"
	"os"
	"strconv"
)

type URLObj struct {
	Key string `json:"key"`
	URL string `json:"url"`
}

type FileStorage struct {
	storagePath string
	memStorage  *MemStorage
}

func NewFileStorage(filePath string) (Storager, error) {
	fs := &FileStorage{
		storagePath: filePath,
		memStorage:  NewMemStorage(),
	}
	if err := fs.ReadURLsFromFile(); err != nil {
		return nil, err
	}
	return fs, nil
}

func (fs *FileStorage) Shorten(uid string, u string) (int, error) {
	id, err := fs.memStorage.Shorten(uid, u)
	if err != nil {
		return -1, err
	}
	if err = fs.WriteURLInFile(strconv.Itoa(id), fs.memStorage.store[strconv.Itoa(id)]); err != nil {
		return -1, err
	}
	return id, nil
}

func (fs *FileStorage) Unshorten(id string) (string, error) {
	url, err := fs.memStorage.Unshorten(id)
	if err != nil {
		return "", err
	}
	return url, nil
}

func (fs *FileStorage) GetAllURLs(uid string) map[string]string {
	usersURLsIds := fs.memStorage.GetAllURLs(uid)
	return usersURLsIds
}

func (fs *FileStorage) ReadURLsFromFile() error {
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
		fs.memStorage.store[u.Key] = u.URL
		fs.memStorage.counter++
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
