package main

import (
	"fmt"
	"grep_code_base/client"
	"grep_code_base/gui"
	"grep_code_base/server"
	"os"
)

func PrintHelp() {
    help := `
  grep_code_base [-h|-t language] keyword
    -h: Help Menu
    
    -t|-type language: target language
    `
    fmt.Print(help)
}


func main () {
    var typeCommand string
    var searchKeyWord string
    if len(os.Args) <= 1 {
        PrintHelp()
        os.Exit(0)
    }
    for i, args := range os.Args { 
        if args == "server" {
            server.Start()
        }else if args == "sync" {
            server.Sync()
           
        }else if args == "-h" {
            PrintHelp()
           
        } else if args == "-t" || args == "-type" {
            if i + 2 > len(os.Args) - 1 {
                fmt.Println("  Requeired: -t|-type [language] [keyword] ")
                os.Exit(0)
            }
            
            typeCommand = os.Args[i + 1]
            searchKeyWord = os.Args[i + 2]
        }
    }

    if len(typeCommand) <= 0 || len(searchKeyWord) <= 0 {
        fmt.Println("Requeired: -t|-type [language] [keyword] ")
        os.Exit(0)

    }

    if len(searchKeyWord) < 3 {

        fmt.Println("Require search keywords length at lease 3")
        os.Exit(0)
    }

    c := client.RequestCommand {
        LangType: typeCommand,
        Search: searchKeyWord,
        IsDebug: true,
    }
    result := client.CallAsCommand(c)

    gui.Show(result)

}
