package client

import (
	"grep_code_base/grep"
	"log"
	"net/rpc"
)

type RequestCommand struct {
	Command  string
	Search   string
	LangType string
	IsDebug  bool
}

func CallAsCommand(request RequestCommand) *[]grep.Result {
	client, err := rpc.DialHTTP("tcp", "0.0.0.0"+":1234")

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
