package item

import (
	"crypto/md5"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
	"strings"

	"github.com/gabriel-vasile/mimetype"
)

type FileItem struct {
	BasePath string `json:"basepath"` // path from config
	Checksum string `json:"checksum"`
	FileName string `json:"filename"` // rooted at basepath
	FileSize int64  `json:"filesize"`
	FullPath string `json:"fullpath"` // rooted at filesystem root
	MimeType string `json:"mime_type"`
	ModTime  string `json:"mod_time"`
}

func (f *FileItem) addChecksum() error {
	// Hash the first and last 2MB of the file to determine uniqueness
	chunkSize := int64(2 * 1024 * 1024) // 2MB
	hasher := md5.New()

	// Open file
	data, err := os.Open(f.FullPath)
	if err != nil {
		return err
	}
	defer data.Close()

	// Read at most 4MB of every file to compare uniqueness
	if f.FileSize > 2*chunkSize {
		// Only read first and last chunks
		n, err := io.CopyN(hasher, data, chunkSize)
		if err != nil {
			log.Printf("Error reading first chunk: %v", err)
			log.Printf("Bytes read: %v", n)
			return err
		}

		data.Seek(0-(chunkSize+1), io.SeekEnd)
		n, err = io.CopyN(hasher, data, chunkSize)
		if err != nil {
			log.Printf("Error reading last chunk: %v", err)
			log.Printf("Bytes read: %v", n)
			return err
		}
	} else {
		// Read entire file, up to two chunks
		n, err := io.CopyN(hasher, data, 2*chunkSize)
		if int64(n) < f.FileSize {
			log.Printf("Bytes read: %v", n)
			if err != nil {
				log.Printf("Error reading file: %v", err)
				return err
			}
		}
	}

	// Add file name and size for extra uniqueness
	name := []byte(f.FileName)
	size := []byte(fmt.Sprintf("%d", f.FileSize))
	hasher.Write(name)
	hasher.Write(size)

	// Calculate hash
	sum := hasher.Sum(nil)

	f.Checksum = fmt.Sprintf("%x", sum)
	//log.Printf("Checksum for %s: %s", f.FullPath, f.Checksum)

	return nil
}

func (f *FileItem) addMimeType() error {
	mType, err := mimetype.DetectFile(f.FullPath)
	if err != nil {
		return err
	}
	f.MimeType = mType.String()
	return nil
}

func (this *FileItem) SameFile(that *FileItem) (bool, error) {
	if that == nil {
		return false, nil
	}

	if this.Checksum == that.Checksum {
		if this.FileSize == that.FileSize {
			return true, nil
		} else {
			return false, errors.New("Error: checksum match; size mismatch")
		}
	}

	return false, nil
}

func StatFile(d fs.DirEntry, fPath string, cPath string) (*FileItem, error) {
	f := &FileItem{}

	dInfo, err := d.Info()
	if err != nil {
		return nil, err
	}

	fileName := strings.TrimPrefix(fPath, cPath)
	if strings.HasPrefix(fileName, string(os.PathSeparator)) {
		fileName = fileName[1:]
	}

	f.BasePath = cPath
	f.FileName = fileName
	f.FileSize = int64(dInfo.Size())
	f.FullPath = fPath
	f.ModTime = dInfo.ModTime().String()

	err = f.addChecksum()
	if err != nil {
		return nil, err
	}

	err = f.addMimeType()
	if err != nil {
		return nil, err
	}

	return f, nil
}
