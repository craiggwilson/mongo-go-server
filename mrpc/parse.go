package mrpc

import (
	"github.com/craiggwilson/mongo-go-server/mrpc/internal"
	"github.com/craiggwilson/mongo-go-server/mrpc/tree"
)

func ParseFile(filename string) (*tree.Tree, error) {
	t, err := internal.ParseFile(filename)
	if err != nil {
		return nil, err
	}

	tr := t.(*tree.Tree)
	tr.Filename = filename

	return tr, nil
}
