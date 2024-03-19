package client

import (
	"grep_code_base/grep"
	"log"
	"net/rpc"
)

const (
	host = "0.0.0.0"
	port = ":1234"
)

type RequestCommand struct {
	Command  string
	Search   string
	LangType string
	IsDebug  bool
}

func CallAsCommand(request RequestCommand) *[]grep.Result {
	client, err := rpc.DialHTTP("tcp", host+port)

	if err != nil {
		log.Fatalln(err.Error())
	}

	var results []grep.Result

	err = client.Call("RPCServer.GrepCode", request, &results)
	defer client.Close()

	if err != nil {
		log.Fatalln(err.Error())
	}

	return &results
}
