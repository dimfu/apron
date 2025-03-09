package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/dimfu/apron/parser"
	"github.com/dimfu/apron/scanner"
)

var (
	path string
	args = make([]string, 0)
)

func init() {
	flag.Parse()
	for i := len(os.Args) - len(flag.Args()) + 1; i < len(os.Args); {
		if i > 1 && os.Args[i-2] == "--" {
			break
		}
		args = append(args, flag.Arg(0))
		if err := flag.CommandLine.Parse(os.Args[i:]); err != nil {
			log.Fatal("error while parsing arguments")
		}

		i += 1 + len(os.Args[i:]) - len(flag.Args())
	}
	args = append(args, flag.Args()...)

	if len(args) < 1 {
		fmt.Println("error: file path should be specified as an argument")
		os.Exit(0)
	}
	path = args[0]
}

func main() {
	source, err := os.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}

	scanner, err := scanner.New(source)
	if err != nil {
		log.Fatal(err)
	}

	p, err := parser.New(scanner.Tokens)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(p)
}
