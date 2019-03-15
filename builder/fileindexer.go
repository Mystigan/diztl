package builder

import (
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/gravetii/diztl/diztl"
	"github.com/gravetii/diztl/util"
)

const (
	sharePathSuffix = "/Documents/diztl/share"
)

// FileIndexer : The struct type that represents a file indexer on a node which indexes all the shared files.
type FileIndexer struct {
	files []diztl.FileMetadata
}

// Index : Indexes all the files present in the default share directory thus making them available for discovery by peers.
func (f *FileIndexer) Index() {
	log.Println("Indexing files...")
	dir := shareDir()
	files := fileWalk(dir)
	f.files = files
	log.Println("Finished indexing...")
}

// Search : Searches for a given pattern in the names of the indexed files and returns those files.
func (f *FileIndexer) Search(pattern string) []*diztl.FileMetadata {
	result := []*diztl.FileMetadata{}
	for _, file := range f.files {
		if strings.Contains(file.Name, pattern) {
			result = append(result, &file)
		}
	}

	return result
}

// Returns the path to the default share directory.
func shareDir() string {
	rootdir, err := os.UserHomeDir()
	if err != nil {
		log.Printf("Unable to fetch user's home directory: %v", err)
		panic(err)
	}

	dir := rootdir + sharePathSuffix
	return dir
}

// Performs a recursive file walk of the given directory path.
func fileWalk(dir string) []diztl.FileMetadata {
	files := []diztl.FileMetadata{}
	var start int32
	counter := util.Counter{Count: &start}
	log.Println("File walk --->")
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}

		counter.IncrBy1()
		metadata := diztl.FileMetadata{Id: counter.Value(), Name: info.Name()}
		files = append(files, metadata)
		log.Printf("%d: %s", metadata.Id, metadata.Name)
		return nil
	})

	if err != nil {
		panic(err)
	}

	return files
}
