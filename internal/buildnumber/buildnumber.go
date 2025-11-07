package buildnumber

import (
	"bytes"
	"errors"
	"fmt"
	"strconv"

	"github.com/anselstetter/git-build-number/internal/repository"
)

var (
	ErrBuildNumberNotFound = errors.New("could not find build number")
	ErrNoHead              = errors.New("could not find head")
	ErrZeroBuildNumber     = errors.New("build number can't be zero")
	ErrInvalidBuildNumber  = errors.New("build number is invalid")
	ErrInvalidHash         = errors.New("hash is invalid")
	ErrInvalidFormat       = errors.New("format is invalid")
)

type Namespace struct {
	Name  string
	Entry Entry
}

type Entry struct {
	Number int64
	Hash   string
}

type BuildNumber struct {
	repository repository.Repository
	fileName   string
	refName    string
}

func (bn *BuildNumber) Hash(namespace string, number int64) (*Entry, error) {
	commits, err := bn.repository.Commits(bn.ref(namespace), repository.WithHeaderKey(strconv.FormatInt(number, 10)))
	if err != nil && errors.Is(err, repository.ErrReferenceNotFound) {
		return nil, errors.Join(err, ErrBuildNumberNotFound)
	} else if err != nil {
		return nil, err
	}
	if len(commits) != 1 {
		return nil, ErrBuildNumberNotFound
	}
	commit := commits[0]
	entry := Entry{
		Number: number,
		Hash:   commit.Headers[0].Value,
	}
	return &entry, nil
}

func (bn *BuildNumber) Get(namespace string, user string, email string, create bool) (*Entry, error) {
	content, err := bn.repository.Content(bn.ref(namespace), bn.fileName)
	if err != nil && errors.Is(err, repository.ErrReferenceNotFound) {
		if !create {
			return nil, errors.Join(err, ErrBuildNumberNotFound)
		}
		return bn.Set(namespace, user, email, 1)
	} else if err != nil {
		return nil, err
	}
	entry, err := Unmarshal(*content)
	if err != nil {
		return nil, err
	}
	return entry, nil
}

func (bn *BuildNumber) Inc(namespace string, user string, email string, force bool) (*Entry, bool, error) {
	entry, err := bn.Get(namespace, user, email, true)
	if err != nil {
		return nil, false, err
	}
	head, err := bn.repository.Head()
	if err != nil && errors.Is(err, repository.ErrReferenceNotFound) {
		return nil, false, errors.Join(err, ErrNoHead)
	} else if err != nil {
		return nil, false, err
	}
	if head.Hash == entry.Hash && !force {
		return entry, false, nil
	}
	entry, err = bn.Set(namespace, user, email, entry.Number+1)
	if err != nil {
		return nil, false, err
	}
	return entry, true, nil
}

func (bn *BuildNumber) Set(namespace string, user string, email string, number int64) (*Entry, error) {
	head, err := bn.repository.Head()
	if err != nil && errors.Is(err, repository.ErrReferenceNotFound) {
		return nil, errors.Join(err, ErrNoHead)
	} else if err != nil {
		return nil, err
	}
	entry := Entry{
		Number: number,
		Hash:   head.Hash,
	}
	content, err := Marshal(entry)
	if err != nil {
		return nil, err
	}
	msg := fmt.Sprintf("Set build number to %d for %s\n", number, head.Hash)

	_, err = bn.repository.Commit(bn.ref(namespace), bn.fileName, content, msg,
		repository.WithAuthor(repository.Author{Name: user, Email: email}),
		repository.WithHeaders([]repository.Header{{Key: strconv.FormatInt(number, 10), Value: head.Hash}}),
	)
	if err != nil {
		return nil, err
	}
	return &entry, nil
}

func (bn *BuildNumber) Delete(namespaces ...string) error {
	for _, namespace := range namespaces {
		err := bn.repository.Delete(bn.ref(namespace))
		if err != nil {
			return err
		}
	}
	return nil
}

func (bn *BuildNumber) Clear() error {
	namespaces, err := bn.Namespaces()
	if err != nil {
		return err
	}
	for _, namespace := range namespaces {
		err := bn.Delete(namespace.Name)
		if err != nil {
			return err
		}
	}
	return nil
}

func (bn *BuildNumber) Namespaces() ([]Namespace, error) {
	refs, err := bn.repository.Refs(repository.WithPrefix(bn.refName))
	if err != nil {
		return nil, err
	}
	namespaces := make([]Namespace, 0, len(refs))

	for _, ref := range refs {
		entry, err := bn.Get(ref.Name, "", "", false)
		if err != nil {
			return nil, err
		}
		namespaces = append(namespaces, Namespace{Name: ref.Name, Entry: *entry})
	}
	return namespaces, nil
}

func (bn *BuildNumber) Mirror(remoteName string) error {
	return bn.repository.Mirror(fmt.Sprintf("%s/", bn.refName), remoteName)
}

func (bn *BuildNumber) Push(remoteName string) error {
	return bn.repository.Push(fmt.Sprintf("%s/*", bn.refName), remoteName, true)
}

func (bn *BuildNumber) Fetch(remoteName string) error {
	return bn.repository.Fetch(fmt.Sprintf("%s/*", bn.refName), remoteName, true)
}

func (bn *BuildNumber) ref(namespace string) string {
	return fmt.Sprintf("%s/%s", bn.refName, namespace)
}

func Marshal(entry Entry) ([]byte, error) {
	if entry.Number == int64(0) {
		return nil, ErrZeroBuildNumber
	}
	if entry.Hash == "" {
		return nil, ErrInvalidHash
	}
	return fmt.Appendf(nil, "%d %s", entry.Number, entry.Hash), nil
}

func Unmarshal(data []byte) (*Entry, error) {
	fields := bytes.Fields(data)

	if len(fields) < 2 {
		return nil, fmt.Errorf("%w: %q", ErrInvalidFormat, data)
	}
	number, err := strconv.ParseInt(string(fields[0]), 10, 64)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrInvalidBuildNumber, err)
	}
	hash := string(fields[1])

	return &Entry{
		Number: number,
		Hash:   hash,
	}, nil
}

func New(repository repository.Repository) BuildNumber {
	return BuildNumber{
		repository: repository,
		fileName:   "build-number",
		refName:    "refs/build-number",
	}
}
