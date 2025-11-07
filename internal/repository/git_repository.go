package repository

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"slices"
	"strings"
	"time"

	"github.com/go-git/go-git/v6"
	"github.com/go-git/go-git/v6/config"
	"github.com/go-git/go-git/v6/plumbing"
	"github.com/go-git/go-git/v6/plumbing/filemode"
	"github.com/go-git/go-git/v6/plumbing/object"
	"github.com/go-git/go-git/v6/storage"
	"github.com/go-git/go-git/v6/storage/memory"
)

type GitRepository struct {
	repo *git.Repository
}

func (g *GitRepository) Head() (*Ref, error) {
	reference, err := g.repo.Head()
	if err != nil {
		return nil, mapError(err)
	}
	ref := Ref{
		Path: reference.Name().String(),
		Name: path.Base(reference.Name().String()),
		Hash: reference.Hash().String(),
	}
	return &ref, nil
}

func (g *GitRepository) Refs(opts ...refsOption) ([]Ref, error) {
	options := newRefsOptions(opts...)

	references, err := g.repo.References()
	if err != nil {
		return nil, mapError(err)
	}
	refs := []Ref{}
	_ = references.ForEach(func(ref *plumbing.Reference) error {
		if ref.Type() != plumbing.HashReference {
			return nil
		}
		r := Ref{
			Path: ref.Name().String(),
			Name: path.Base(ref.Name().String()),
			Hash: ref.Hash().String(),
		}
		if options.prefix != nil && !strings.HasPrefix(r.Path, *options.prefix) {
			return nil
		}
		refs = append(refs, r)
		return nil
	})
	return refs, nil
}

func (g *GitRepository) RemoteRefs(remoteName string, opts ...refsOption) ([]Ref, error) {
	options := newRefsOptions(opts...)

	remote, err := g.repo.Remote(remoteName)
	if err != nil {
		return nil, mapError(err)
	}
	remoteRefs, err := remote.List(&git.ListOptions{})
	if err != nil {
		return nil, mapError(err)
	}
	refs := []Ref{}
	for _, ref := range remoteRefs {
		r := Ref{
			Path: ref.Name().String(),
			Name: path.Base(ref.Name().String()),
			Hash: ref.Hash().String(),
		}
		if options.prefix != nil && !strings.HasPrefix(r.Path, *options.prefix) {
			continue
		}
		refs = append(refs, r)
	}
	return refs, nil
}

func (g *GitRepository) Content(refName string, fileName string) (*[]byte, error) {
	ref, err := g.repo.Reference(plumbing.ReferenceName(refName), true)
	if err != nil {
		return nil, mapError(err)
	}
	commit, err := g.repo.CommitObject(ref.Hash())
	if err != nil {
		return nil, mapError(err)
	}
	tree, err := commit.Tree()
	if err != nil {
		return nil, mapError(err)
	}
	file, err := tree.File(fileName)
	if err != nil {
		return nil, mapError(err)
	}
	reader, err := file.Reader()
	if err != nil {
		return nil, mapError(err)
	}
	defer reader.Close()

	content, err := io.ReadAll(reader)
	if err != nil {
		return nil, mapError(err)
	}
	return &content, nil
}

func (g *GitRepository) Commits(refName string, opts ...commitsOption) ([]Commit, error) {
	options := newCommitsOptions(opts...)
	errStop := errors.New("stop iteration")

	ref, err := g.repo.Reference(plumbing.ReferenceName(refName), true)
	if err != nil {
		return nil, mapError(err)
	}
	commit, err := g.repo.CommitObject(ref.Hash())
	if err != nil {
		return nil, mapError(err)
	}
	commits := []Commit{}

	err = object.NewCommitPreorderIter(commit, nil, nil).ForEach(func(c *object.Commit) error {
		headers := []Header{}
		for _, header := range c.ExtraHeaders {
			headers = append(headers, Header{Key: header.Key, Value: header.Value})
		}
		commit := Commit{
			Hash:    c.Hash.String(),
			Author:  Author{Name: c.Author.Name, Email: c.Author.Email},
			When:    c.Author.When,
			Message: c.Message,
			Headers: headers,
		}
		commits = append(commits, commit)

		if options.headerKey != nil && slices.ContainsFunc(headers, func(header Header) bool { return header.Key == *options.headerKey }) {
			commits = []Commit{commit}
			return errStop
		}

		return nil
	})
	if err == errStop {
		return commits, nil
	} else if err != nil {
		return nil, mapError(err)
	}
	return commits, nil
}

func (g *GitRepository) Commit(refName, fileName string, content []byte, msg string, opts ...commitOption) (*Ref, error) {
	options := newCommitOptions(opts...)
	store := g.repo.Storer

	parents := []plumbing.Hash{}
	ref, err := g.repo.Reference(plumbing.ReferenceName(refName), true)
	if err == nil {
		parents = append(parents, ref.Hash())
	} else if !errors.Is(err, plumbing.ErrReferenceNotFound) {
		return nil, err
	}

	blobHash, err := storeBlob(store, content)
	if err != nil {
		return nil, err
	}
	tree := &object.Tree{
		Entries: []object.TreeEntry{{
			Name: fileName,
			Mode: filemode.Regular,
			Hash: blobHash,
		}},
	}
	treeHash, err := storeObject(store, tree)
	if err != nil {
		return nil, err
	}
	extraHeaders := []object.ExtraHeader{}
	for _, header := range options.headers {
		extraHeaders = append(extraHeaders, object.ExtraHeader{Key: header.Key, Value: header.Value})
	}
	commit := &object.Commit{
		Author: object.Signature{
			Name:  options.author.Name,
			Email: options.author.Email,
			When:  time.Now(),
		},
		Message:      msg,
		TreeHash:     treeHash,
		ParentHashes: parents,
		ExtraHeaders: extraHeaders,
	}
	commitHash, err := storeObject(store, commit)
	if err != nil {
		return nil, err
	}
	newRef := plumbing.NewHashReference(plumbing.ReferenceName(refName), commitHash)
	if err := store.SetReference(newRef); err != nil {
		return nil, err
	}
	if options.setHead {
		headRef := plumbing.NewSymbolicReference(plumbing.HEAD, newRef.Name())
		if err := store.SetReference(headRef); err != nil {
			return nil, err
		}
	}
	return &Ref{
		Path: newRef.Name().String(),
		Name: path.Base(newRef.Name().String()),
		Hash: commitHash.String(),
	}, nil
}

func (g *GitRepository) Delete(refName string) error {
	ref := plumbing.ReferenceName(refName)

	_, err := g.repo.Reference(ref, false)
	if err != nil {
		return mapError(err)
	}
	err = g.repo.Storer.RemoveReference(ref)
	if err != nil {
		return mapError(err)
	}
	return nil
}

func (g *GitRepository) Fetch(refName string, remoteName string, force bool) error {
	spec := fmt.Sprintf("%s:%s", refName, refName)

	err := g.repo.Fetch(&git.FetchOptions{
		RemoteName: remoteName,
		RefSpecs: []config.RefSpec{
			config.RefSpec(spec),
		},
		Force: force,
	})
	if err != nil && err != git.NoErrAlreadyUpToDate {
		return mapError(err)
	}
	return nil
}

func (g *GitRepository) Mirror(refName string, remoteName string) error {
	spec := fmt.Sprintf("+%s*:%s*", refName, refName)
	localRefs := map[string]bool{}

	refs, err := g.Refs(WithPrefix(refName))
	if err != nil {
		return mapError(err)
	}
	for _, ref := range refs {
		localRefs[ref.Path] = true
	}
	err = g.repo.Push(&git.PushOptions{
		RemoteName: remoteName,
		RefSpecs: []config.RefSpec{
			config.RefSpec(spec),
		},
		Force: true,
	})
	if err != nil && err != git.NoErrAlreadyUpToDate {
		return mapError(err)
	}
	remoteRefs, err := g.RemoteRefs(remoteName, WithPrefix(refName))
	if err != nil {
		return mapError(err)
	}
	for _, ref := range remoteRefs {
		if !localRefs[ref.Path] {
			delSpec := fmt.Sprintf(":%s", ref.Path)

			err := g.repo.Push(&git.PushOptions{
				RemoteName: remoteName,
				RefSpecs: []config.RefSpec{
					config.RefSpec(delSpec),
				},
				Force: true,
			})
			if err != nil {
				return fmt.Errorf("delete %s: %w", ref.Path, err)
			}
		}
	}
	return nil
}

func (g *GitRepository) Push(refName string, remoteName string, force bool) error {
	spec := fmt.Sprintf("%s:%s", refName, refName)

	err := g.repo.Push(&git.PushOptions{
		RemoteName: remoteName,
		RefSpecs: []config.RefSpec{
			config.RefSpec(spec),
		},
		Force: force,
	})
	if err != nil && err != git.NoErrAlreadyUpToDate {
		return mapError(err)
	}
	return nil
}

func (g *GitRepository) AddRemote(name string, urls ...string) error {
	_, err := g.repo.CreateRemote(&config.RemoteConfig{
		Name: name,
		URLs: urls,
	})
	if err != nil {
		return err
	}
	return nil
}

func storeBlob(store storage.Storer, data []byte) (plumbing.Hash, error) {
	obj := &plumbing.MemoryObject{}
	obj.SetType(plumbing.BlobObject)
	obj.SetSize(int64(len(data)))

	if _, err := obj.Write(data); err != nil {
		return plumbing.ZeroHash, err
	}
	return store.SetEncodedObject(obj)
}

func storeObject(store storage.Storer, obj object.Object) (plumbing.Hash, error) {
	mem := &plumbing.MemoryObject{}
	if err := obj.Encode(mem); err != nil {
		return plumbing.ZeroHash, err
	}
	return store.SetEncodedObject(mem)
}

func mapError(err error) error {
	switch {
	case errors.Is(err, plumbing.ErrReferenceNotFound):
		return ErrReferenceNotFound
	case errors.Is(err, git.ErrRemoteNotFound):
		return ErrRemoteNotFound
	default:
		return err
	}
}

func randomRepositoryName() (string, error) {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return fmt.Sprintf("%s.git", hex.EncodeToString(b)), nil
}

func NewGitRepository(path string) (*GitRepository, error) {
	repo, err := git.PlainOpenWithOptions(path, &git.PlainOpenOptions{DetectDotGit: true})
	if err != nil {
		return nil, err
	}
	repository := GitRepository{
		repo: repo,
	}
	return &repository, nil
}

func NewGitTempBareRepository(initialCommit bool) (*GitRepository, *string, error) {
	name, err := randomRepositoryName()
	if err != nil {
		return nil, nil, err
	}
	tmpDir := os.TempDir()
	remotePath := filepath.Join(tmpDir, name)
	repo, err := git.PlainInit(remotePath, true)
	if err != nil {
		return nil, nil, err
	}
	repository := GitRepository{
		repo: repo,
	}
	if initialCommit {
		_, err = repository.Commit("refs/heads/main", "initial", []byte(""), "Initial commit", WithHead())
		if err != nil {
			return nil, nil, mapError(err)
		}
	}
	return &repository, &remotePath, nil
}

func NewGitInMemoryRepository(initialCommit bool) (*GitRepository, *Ref, error) {
	var ref *Ref

	repo, err := git.Init(memory.NewStorage(), nil)
	if err != nil {
		return nil, ref, err
	}
	repository := GitRepository{
		repo: repo,
	}
	if initialCommit {
		ref, err = repository.Commit("refs/heads/main", "initial", []byte(""), "Initial commit", WithHead())
		if err != nil {
			return nil, ref, mapError(err)
		}
	}
	return &repository, ref, nil
}

var _ Repository = &GitRepository{}
