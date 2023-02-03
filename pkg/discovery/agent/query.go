// Copyright 2023 NJWS Inc.

package agent

import (
	"fmt"

	"git.fg-tech.ru/listware/cmdb/pkg/cmdb/documents"
	"git.fg-tech.ru/listware/cmdb/pkg/cmdb/qdsl"
)

func (a *Agent) getDocument(query string) (document *documents.Node, err error) {
	log.Infof("qdsl (%s)", query)
	documents, err := qdsl.Qdsl(a.ctx, query, qdsl.WithKey(), qdsl.WithId())
	if err != nil {
		return
	}
	for _, document = range documents {
		return
	}
	err = fmt.Errorf("document '%s' not found", query)
	return
}
