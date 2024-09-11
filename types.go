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

type LocalValue struct {
	Info struct {
		Name       string `json:"name"`
		TypeID     string `json:"typeID"`
		TypeName   string `json:"type"`
		IsVariable bool   `json:"isVariable"`
		FunctionID string `json:"functionID"`
	} `json:"info"`
	Range    FileRange   `json:"origin"`
	Mentions []FileRange `json:"mentions"`
}

type Function struct {
	Info struct {
		Name            string            `json:"name"`
		FullName        string            `json:"fullName"`
		ID              string            `json:"functionID"`
		GenericArgs     []GenericArgument `json:"genericArguments"`
		ModuleID        string            `json:"moduleID"`
		Visibility      Visibility        `json:"visibility"`
		IsVariadic      bool              `json:"isVariadic"`
		LocalValues     []LocalValue      `json:"locals"`
		DefinitionRange *FileRange        `json:"definitionRange,omitempty"`
	} `json:"info"`
	Range           FileRange        `json:"origin"`
	Mentions        []FileRange      `json:"mentions"`
	BroughtMentions []BroughtMention `json:"broughtMentions"`
}

// Function - GenericEntity interface

func (self Function) IsGeneric() bool {
	return len(self.Info.GenericArgs) != 0
}

func (self Function) GenericArgumentCount() int {
	return len(self.Info.GenericArgs)
}

func (self Function) HasGenericArgument(name string) bool {
	for _, arg := range self.Info.GenericArgs {
		if arg.Name == name {
			return true
		}
	}
	return false
}

func (self Function) GetGenericArgument(name string) *GenericArgument {
	for _, arg := range self.Info.GenericArgs {
		if arg.Name == name {
			return &arg
		}
	}
	return nil
}

type GenericFunction struct {
	Info struct {
		Name              string                     `json:"name"`
		FullName          string                     `json:"fullName"`
		ID                string                     `json:"functionID"`
		GenericParameters []GenericAbstractParameter `json:"genericParameters"`
		ModuleID          string                     `json:"moduleID"`
		Visibility        Visibility                 `json:"visibility"`
	} `json:"info"`
	Range           FileRange        `json:"origin"`
	Mentions        []FileRange      `json:"mentions"`
	BroughtMentions []BroughtMention `json:"broughtMentions"`
}

type StructField struct {
	Info struct {
		Name       string     `json:"name"`
		TypeID     string     `json:"typeID"`
		TypeName   string     `json:"type"`
		IsVar      bool       `json:"isVariable"`
		Visibility Visibility `json:"visibility"`
	} `json:"info"`
	Range    FileRange   `json:"origin"`
	Mentions []FileRange `json:"mentions"`
}

type GenericAbstractParameter struct {
	Kind         string      `json:"genericKind"`
	Index        json.Number `json:"index"`
	Name         string      `json:"name"`
	HasDefault   bool        `json:"hasDefault"`
	DefaultValue *string     `json:"defaultValueString,omitempty"`
	Range        FileRange   `json:"range"`
}

type TypeField struct {
	Type Type
	Name string
}

type TypeStaticField struct {
	Type Type
	Name string
}

type Type interface {
	GetID() string
	GetFullName() string

	IsTriviallyCopyable() bool
	IsTriviallyMovable() bool
	IsCopyConstructible() bool
	IsCopyAssignable() bool
	IsMoveConstructible() bool
	IsMoveAssignable() bool

	HasFields() bool
	GetFields() []TypeField

	HasStaticFields() bool
	GetStaticFields() []TypeStaticField

	HasMethods() bool
	GetMethods() []MethodSignature

	HasStaticMethods() bool
	GetStaticMethods() []FunctionSignature
}

type GenericEntity interface {
	IsGeneric() bool
	GenericArgumentCount() uint32
	HasGenericArgument(name string) bool
	GetGenericArgument(name string) *GenericArgument
}

var allTypes map[string]Type
