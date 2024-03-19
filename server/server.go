package server

import (
	"errors"
	"fmt"
	"grep_code_base/database"
	db "grep_code_base/database"
	"grep_code_base/grep"
	"io/fs"
	"log"
	"net"
	"net/http"
	"net/rpc"
	"os"
	"path/filepath"
	_ "runtime"
	"strings"
	"sync"
	"time"
)

const (
	mainDir      = "code"
	port         = ":1234"
	numOfWorkers = 20 // the workers are for scan file content.
)

type RPCServer int64

type RequestCommand struct {
	Command  string
	Search   string
	LangType string
	IsDebug  bool
}

func (server *RPCServer) GrepCode(args *RequestCommand, reply *[]grep.Result) error {
	var err error
	var start time.Time
	//var m runtime.MemStats
	if args.IsDebug {
		start = time.Now()
	}
	*reply, err = Grep(args.LangType, args.Search)
	if err != nil {
		return err
	}

	if args.IsDebug {
		log.Println(time.Since(start))
		// runtime.ReadMemStats(&m)
		// fmt.Printf("Alloc = %v MiB", m.Alloc/1024/1024)
		// fmt.Printf("\tTotalAlloc = %v MiB", m.TotalAlloc/1024/1024)
		// fmt.Printf("\tSys = %v MiB", m.Sys/1024/1024)
		// fmt.Printf("\tNumGC = %v\n", m.NumGC)

	}
	return nil
}

func insertToDirList(slice []string, name string) []string {
	existed := false

	for _, path := range slice {
		if name == path {
			existed = true
			break
		}
	}

	if !existed {
		slice = append(slice, name)
	}

	return slice
}

func Sync() {
	var codebaseDir string
	var codebaseDirStrLength int = 0
	var codeBases []db.CodeBaseFolder
	var dirList []string

	usrDir, err := os.UserHomeDir()
	if err != nil {
		panic(err.Error())
	}

	codebaseDir = fmt.Sprintf("%s/%s", usrDir, mainDir)
	codebaseDirStrLength = len(codebaseDir)

	err = filepath.WalkDir(codebaseDir, func(path string, info fs.DirEntry, err error) error {

		if err != nil {
			return err
		}
		if info.IsDir() && !grep.IsIgnoreDir(info.Name()) {
			if codebaseDirStrLength < len(path) {

				pp := path[codebaseDirStrLength+1:]
				count := strings.Count(pp, string(os.PathSeparator))

				if count == 0 {
					// language Type
				}

				// Domain
				if count == 1 {
					ss := strings.Split(pp, string(os.PathSeparator))
					inner_path := filepath.Join(ss[0], ss[1])
					dirList = insertToDirList(dirList, inner_path)
				}
				// subdomain or more depth folder
				if count > 1 {
					ss := strings.Split(pp, string(os.PathSeparator))
					inner_path := filepath.Join(ss[0], ss[1], ss[2])
					dirList = insertToDirList(dirList, inner_path)
				}
			}
		}

		return nil
	})

	if err != nil {
		panic(err.Error())
	}

	if len(dirList) > 0 {
		fmt.Println("Sync codebase to db...")
	} else {
		fmt.Println("No code folder need to sync..")
		os.Exit(0)
	}

	for _, dirs := range dirList {
		var cb = &db.CodeBaseFolder{}
		ss := strings.Split(dirs, string(os.PathSeparator))
		// The Path on include Domain
		if len(ss) == 2 {
			cb = &db.CodeBaseFolder{}
			cb.Lang = ss[0]
			cb.Domain = ss[1]
			cb.Path = filepath.Join(ss[0], ss[1])
			codeBases = append(codeBases, *cb)
		}
		// The Path include domain, subdomain
		if len(ss) == 3 {
			cb.Lang = ss[0]
			cb.Domain = ss[1]
			cb.Subdomain = ss[2]
			cb.Path = filepath.Join(ss[0], ss[1], ss[2])
			codeBases = append(codeBases, *cb)
		}

	}

	for _, v := range codeBases {
		db.InsertCodeBaseFolder(v)
	}

	fmt.Println("Sync codebase finish!")
}

func scanFile(workersWg *sync.WaitGroup, entryList *grep.EntryList, results *[]grep.Result, searchText string) {
	defer workersWg.Done()
	for {
		entry := entryList.Next()

		if entry.Path == "" {
			return
		}

		result := grep.ScanText(entry.Path, searchText)
		if result != nil {
			*results = append(*results, *result)
		}
	}
}

func Grep(typeCommand string, searchText string) ([]grep.Result, error) {

	dbResults := database.GetCodeBaseFoldersByLangDistinctDomain(typeCommand)

	if len(dbResults) == 0 {
		err := errors.New(fmt.Sprintf("No target language %s in code base found.\n", typeCommand))
		return nil, err
	}

	usrDir, err := os.UserHomeDir()
	if err != nil {
		panic(err.Error())
	}

	var workersWg sync.WaitGroup

	entryList := grep.NewEntryList(len(dbResults))

	workersWg.Add(1)
	// create a goruntinue for pipe all registered paths into channel
	go func() {
		defer workersWg.Done()

		for _, dbResult := range dbResults {
			grep.GetAllFiles(&entryList, filepath.Join(usrDir, mainDir, dbResult.Path))
		}
		entryList.Finish(numOfWorkers)
	}()

	results := new([]grep.Result)

	// another goruntinue to create workers
	for i := 0; i < numOfWorkers; i++ {
		workersWg.Add(1)
		go scanFile(&workersWg, &entryList, results, searchText)
	}

	workersWg.Wait()

	// Max result
	if len(*results) >= 51 {
		*results = (*results)[:51]
	}

	return *results, nil
}

func Start() {
	server := new(RPCServer)

	rpc.Register(server)
	rpc.HandleHTTP()
	l, err := net.Listen("tcp", "0.0.0.0:1234")
	if err != nil {
		log.Fatal("Listen error: ", err)
	}

	http.Serve(l, nil)
}
