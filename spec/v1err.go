package spec

import "fmt"

type ErrSpec1PropInvalidInput struct {
	Path string
	Key  string
}

func (err ErrSpec1PropInvalidInput) Error() string {
	return fmt.Sprintf("Invalid input defined for %s/%s", err.Path, err.Key)
}

type ErrSpec1PropInvalidType struct {
	Path string
	Key  string
	Type string
}

func (err ErrSpec1PropInvalidType) Error() string {
	return fmt.Sprintf("Invalid type defined for %s/%s: %s", err.Path, err.Key, err.Type)
}

type ErrSpec1PropInvalidDefault struct {
	Path    string
	Key     string
	Default interface{}
}

func (err ErrSpec1PropInvalidDefault) Error() string {
	return fmt.Sprintf("Invalid default defined for %s/%s: %s", err.Path, err.Key, err.Default)
}
