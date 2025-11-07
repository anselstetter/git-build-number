package repository

import (
	"errors"
	"time"
)

var (
	ErrReferenceNotFound = errors.New("reference not found")
	ErrRemoteNotFound    = errors.New("remote not found")
)

type Commit struct {
	Hash    string
	Author  Author
	When    time.Time
	Message string
	Headers []Header
}

type Header struct {
	Key   string
	Value string
}

type Author struct {
	Name  string
	Email string
}

type Ref struct {
	Path string
	Name string
	Hash string
}

type Repository interface {
	Head() (*Ref, error)
	Refs(opts ...refsOption) ([]Ref, error)
	Content(refName string, fileName string) (*[]byte, error)
	Commit(refName string, fileName string, content []byte, msg string, opts ...commitOption) (*Ref, error)
	Commits(refName string, opts ...commitsOption) ([]Commit, error)
	Delete(refName string) error
	Fetch(refName string, remoteName string, force bool) error
	Push(refName string, remoteName string, force bool) error
	Mirror(refName string, remoteName string) error
	AddRemote(name string, urls ...string) error
}
