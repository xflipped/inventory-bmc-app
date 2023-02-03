// Copyright 2023 NJWS Inc.

package agent

func (a *Agent) entrypoint() (err error) {
	/*route := &pbtypes.FunctionRoute{
		Url:             fmt.Sprintf("http://inventory-bmc:%d/statefun", a.m.Port()),
		ExecuteOnCreate: true,
		ExecuteOnUpdate: true,
	}

	// exists node or create
	node, err := a.getNode()
	if err == nil {
		// search function
		documents, err := qdsl.Qdsl(a.ctx, fmt.Sprintf("%s.%s", node.Key, types.FunctionPath), qdsl.WithLinkId())
		if err != nil {
			return err
		}

		// func on node: exists
		for _, document := range documents {
			updateLink, err := system.UpdateAdvancedLink(document.LinkId.String(), route)
			if err != nil {
				return err
			}

			return a.executor.ExecSync(a.ctx, updateLink)
		}

		return a.createLink(route)
	}

	// TODO move to register
	nodes, err := a.getNodes()
	if err != nil {
		return
	}

	message := &pbtypes.FunctionMessage{
		// FunctionType: ,
		Route: route,
	}

	createNode, err := system.CreateChild(nodes.Id.String(), types.NodeID, a.hostname(), a.node, message)
	if err != nil {
		return
	}

	if err = a.executor.ExecSync(a.ctx, createNode); err != nil {
		return err
	}

	return a.createLink(route)*/
	return
}
