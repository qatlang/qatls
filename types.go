package main

import (
	"encoding/json"
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

// CODE INFO TYPES

type Visibility struct {
	Nature    string `json:"nature"`
	HasModule bool   `json:"hasModule"`
	ModuleID  string `json:"moduleID"`
	HasType   bool   `json:"hasType"`
	TypeID    string `json:"typeID"`
}

type BroughtMention struct {
	ModuleID string    `json:"module"`
	Range    FileRange `json:"range"`
}

type Mod struct {
	Info struct {
		ID                        string        `json:"moduleID"`
		FullName                  string        `json:"fullName"`
		IsFilesystemLib           bool          `json:"isFilesystemLib"`
		ModuleType                string        `json:"moduleType"`
		Visibility                Visibility    `json:"visibility"`
		HasModuleInitialiser      bool          `json:"hasModuleInitialiser"`
		IntegerBitwidths          []json.Number `json:"integerBitwidths"`
		UnsignedBitwidths         []json.Number `json:"unsignedBitwidths"`
		FilesystemBroughtMentions []FileRange   `json:"filesystemBroughtMentions"`
	} `json:"info"`
	Range           FileRange        `json:"origin"`
	Mentions        []FileRange      `json:"mentions"`
	BroughtMentions []BroughtMention `json:"broughtMentions"`
}
