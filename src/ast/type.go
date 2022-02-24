package ast

type TypeKind string

const (
	String TypeKind = "string"
	Int             = "int"
)

type TypeDeclaration struct {
	Id   string
	Kind TypeKind
}

type ConstDeclaration struct {
	Declarators []ConstDeclarator
}

type ConstDeclarator struct {
	Kind  string
	Id    string
	Value string
}

type File struct {
	Name string
	Body []interface{}
}
