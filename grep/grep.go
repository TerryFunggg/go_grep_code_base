package grep

import (
	"bufio"
	"log"
	"os"
	"path/filepath"
	"strings"
)

type Result struct {
    Path    string
    FileName string
    Body    []byte
}

type Entry struct {
    Path string
}

type EntryList struct {
    jobs chan Entry
}

func (el *EntryList) Add(entry Entry) {
    el.jobs <- entry
}

func (el *EntryList) Next() Entry {
    job := <- el.jobs
    return job
}

func NewEntryList(size int) EntryList {
    return EntryList{make(chan Entry, size)}
}

func NewEntry(path string) Entry {
    return Entry{ path }
}

func (el *EntryList) Finish(numWorkers int) {
    for i := 0; i < numWorkers; i++ {
        el.Add(Entry{""})
    }
}

func IsIgnoreDir(name string) bool {
    ignoreDir := []string {".vscode"}

    for _, dir := range ignoreDir {
        if dir == name {
            return true
        }
    }
    return false
}

func IsIgnoreExt(name string) bool {
    ignoreExt := []string {".lock", ".backup", ".toml", ".min.js"}

    fileExt := filepath.Ext(name)

    for _, ext := range ignoreExt {
        if ext == fileExt {
            return true
        }
    }
    return false
}

func GetAllFiles(fileList *EntryList, path string) {
    enteries, err := os.ReadDir(path)
    if err != nil {
        panic(err.Error())

    }
    for _, entry := range enteries {
        if entry.IsDir() && !IsIgnoreDir(entry.Name()) {
            next := filepath.Join(path, entry.Name())
            GetAllFiles(fileList, next)
        } else {
            if !IsIgnoreExt(entry.Name()) {
                //*fileList = append(*fileList, filepath.Join(path, entry.Name()) )
                fileList.Add(NewEntry(filepath.Join(path, entry.Name())))
            }
        }
    }
}

func ScanText(path string, targetString string) *Result {
    file, err := os.Open(path)

    if err != nil {
        log.Println("Error:", err.Error())
        return nil
    }

    scanner := bufio.NewScanner(file)
    match := false

    for scanner.Scan() {
        if strings.Contains(scanner.Text(), targetString) {
            match = true
            break
        }

    }

    if match {
        buf, err := os.ReadFile(path)
        if err != nil {
            log.Println("Error: ", err.Error())
        }

        r := Result{Path: path, Body: buf }
        return &r
    } else {
        return nil
    }
}
