package main

import (
	protocol "github.com/tliron/glsp/protocol_3_16"
)

type QatSession struct {
	handler    protocol.Handler
	workspaces []protocol.WorkspaceFolder
	files      []File
}

type File struct {
	uri     string
	version protocol.Integer
	content string
	isOpen  bool
}
