package qatls

import (
	"github.com/tliron/glsp"
	protocol "github.com/tliron/glsp/protocol_3_16"
)

type QatClient struct {
	handler    protocol.Handler
	workspaces []protocol.WorkspaceFolder
}

const lsName = "qatls"

var lsVersion string = "0.1.0"

var client QatClient

func main() {
	client.workspaces = []protocol.WorkspaceFolder{}
	client.handler = protocol.Handler{
		Initialize: initFn,
	}
}

func initFn(ctx *glsp.Context, params *protocol.InitializeParams) (any, error) {
	capabilities := client.handler.CreateServerCapabilities()
	if params.WorkspaceFolders != nil {
		client.workspaces = params.WorkspaceFolders
	}
	client.handler.LogTrace(ctx, &protocol.LogTraceParams{
		Message: "Intialised the QAT language server",
	})
	return protocol.InitializeResult{
		Capabilities: capabilities,
		ServerInfo: &protocol.InitializeResultServerInfo{
			Name:    lsName,
			Version: &lsVersion,
		},
	}, nil
}
