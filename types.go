package main

import (
	"os"

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

type Position struct {
	Line      protocol.UInteger `json:"line"`
	Character protocol.UInteger `json:"character"`
}

type FileRange struct {
	Path  string   `json:"path"`
	Start Position `json:"start"`
	End   Position `json:"end"`
}

func (self FileRange) IsSamePathAs(pathVal string) bool {
	selfPathInfo, err := os.Lstat(self.Path)
	if err != nil {
		Server.Log.Error("Error getting FileInfo of path: " + self.Path)
		return false
	}
	pathValInfo, err := os.Lstat(pathVal)
	if err != nil {
		Server.Log.Error("Error getting FileInfo of path: " + pathVal)
		return false
	}
	return os.SameFile(selfPathInfo, pathValInfo)
}

func (self FileRange) IsWithin(pathVal string, pos Position) bool {
	return self.IsSamePathAs(pathVal) && ((self.Start.Line < pos.Line) && (self.End.Line > pos.Line)) || ((self.Start.Line == pos.Line) && (self.Start.Character <= pos.Character) && (self.End.Line == pos.Line) && (self.End.Character >= pos.Character)) || ((self.Start.Line == pos.Line) && (self.Start.Character <= pos.Character)) || ((self.End.Line == pos.Line) && (self.End.Character >= pos.Character))
}

func (self FileRange) IsRightAfter(pathVal string, pos Position) bool {
	return self.IsSamePathAs(pathVal) && (self.End.Line == pos.Line) && ((self.End.Character + 1) == pos.Character)
}
