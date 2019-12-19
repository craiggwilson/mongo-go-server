package main

import (
	"bytes"
	"flag"
	"io"
	"os"

	"github.com/craiggwilson/mongo-go-server/mrpc"
)

func main() {
	flag.Parse()

	args := flag.Args()
	if len(args) == 0 {
		panic("must specify a rpc definition file")
	}

	tree, err := mrpc.ParseFile(args[0])
	if err != nil {
		panic(err)
	}

	var buf bytes.Buffer

	err = mrpc.Generate(&buf, tree)
	if err != nil {
		panic(err)
	}

	_, _ = io.Copy(os.Stdout, &buf)
}
