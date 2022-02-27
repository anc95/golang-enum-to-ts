package ast

type TypeKind string

const (
	String TypeKind = "string"
	Int             = "int"
)

type BaseDeclaration struct {
	LeadingComments  []Comment
	TrailingComments []Comment
}

type TypeDeclaration struct {
	BaseDeclaration
	Id   string
	Kind TypeKind
}

type ConstDeclaration struct {
	BaseDeclaration
	Declarators []ConstDeclarator
}

type Comment struct {
	Value string
}

type ConstDeclarator struct {
	BaseDeclaration
	Kind  string
	Id    string
	Value string
}

type File struct {
	Name string
	Body []interface{}
}
