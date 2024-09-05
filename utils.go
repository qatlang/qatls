package qatls

import (
	"github.com/tliron/glsp/server"
)

const LspName = "qatls"

var (
	Client     QatClient
	Server     *server.Server = nil
	LspVersion                = "0.1.0"
)

func FindFileByURI(uri string) *File {
	for _, file := range Client.files {
		if file.uri == uri {
			return &file
		}
	}
	return nil
}
