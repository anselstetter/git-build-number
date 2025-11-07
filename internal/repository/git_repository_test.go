package repository_test

import (
	"os"
	"testing"

	"github.com/anselstetter/git-build-number/internal/repository"
	"github.com/stretchr/testify/assert"
)

func addRemote(t *testing.T, remoteName string, initialCommit bool, local repository.Repository) repository.Repository {
	t.Helper()

	remote, remotePath, _ := repository.NewGitTempBareRepository(initialCommit)
	t.Cleanup(func() {
		_ = os.RemoveAll(*remotePath)
	})
	_ = local.AddRemote(remoteName, *remotePath)
	return remote
}

func TestHead(t *testing.T) {
	t.Run("no head", func(t *testing.T) {
		t.Parallel()

		repo, _, err := repository.NewGitInMemoryRepository(false)
		assert.NoError(t, err)

		ref, err := repo.Head()
		assert.Nil(t, ref)
		assert.ErrorIs(t, err, repository.ErrReferenceNotFound)
	})

	t.Run("head", func(t *testing.T) {
		t.Parallel()

		repo, ref, err := repository.NewGitInMemoryRepository(true)
		assert.NoError(t, err)

		head, err := repo.Head()
		assert.NoError(t, err)

		assert.Equal(t, head, ref)
	})
}

func TestCommit(t *testing.T) {
	repo, _, err := repository.NewGitInMemoryRepository(false)
	assert.NoError(t, err)

	_, err = repo.Commit("refs/heads/main", "test", []byte(""), "commit", repository.WithHead())
	assert.NoError(t, err)

	ref, err := repo.Commit("refs/custom/main", "test", []byte(""), "commit", repository.WithHead())
	assert.NoError(t, err)

	head, err := repo.Head()
	assert.NoError(t, err)

	assert.Equal(t, head, ref)
}

func TestCommits(t *testing.T) {
	t.Run("without header keys", func(t *testing.T) {
		t.Parallel()

		repo, _, err := repository.NewGitInMemoryRepository(false)
		assert.NoError(t, err)

		ref1, err := repo.Commit("refs/custom/test", "test", []byte(""), "commit",
			repository.WithAuthor(repository.Author{Name: "name", Email: "mail"}),
			repository.WithHeaders([]repository.Header{{Key: "key", Value: "value"}}),
		)
		assert.NoError(t, err)

		ref2, err := repo.Commit("refs/custom/test", "test", []byte(""), "commit",
			repository.WithAuthor(repository.Author{Name: "name", Email: "mail"}),
			repository.WithHeaders([]repository.Header{{Key: "key2", Value: "value2"}}),
		)
		assert.NoError(t, err)

		commits, err := repo.Commits("refs/custom/test")
		assert.NoError(t, err)

		assert.Equal(t, len(commits), 2)

		assert.Equal(t, commits[0].Hash, ref2.Hash)
		assert.Equal(t, commits[0].Author, repository.Author{Name: "name", Email: "mail"})
		assert.Equal(t, commits[0].Message, "commit")
		assert.Equal(t, commits[0].Headers, []repository.Header{{Key: "key2", Value: "value2"}})

		assert.Equal(t, commits[1].Hash, ref1.Hash)
		assert.Equal(t, commits[1].Author, repository.Author{Name: "name", Email: "mail"})
		assert.Equal(t, commits[1].Message, "commit")
		assert.Equal(t, commits[1].Headers, []repository.Header{{Key: "key", Value: "value"}})
	})
	t.Run("with header keys", func(t *testing.T) {
		t.Parallel()

		repo, _, err := repository.NewGitInMemoryRepository(false)
		assert.NoError(t, err)

		_, err = repo.Commit("refs/custom/test", "test", []byte(""), "commit",
			repository.WithAuthor(repository.Author{Name: "name", Email: "mail"}),
			repository.WithHeaders([]repository.Header{{Key: "key", Value: "value"}}),
		)
		assert.NoError(t, err)

		ref2, err := repo.Commit("refs/custom/test", "test", []byte(""), "commit",
			repository.WithAuthor(repository.Author{Name: "name", Email: "mail"}),
			repository.WithHeaders([]repository.Header{{Key: "key2", Value: "value2"}}),
		)
		assert.NoError(t, err)

		commits, err := repo.Commits("refs/custom/test", repository.WithHeaderKey("key2"))
		assert.NoError(t, err)

		assert.Equal(t, len(commits), 1)

		assert.Equal(t, commits[0].Hash, ref2.Hash)
		assert.Equal(t, commits[0].Author, repository.Author{Name: "name", Email: "mail"})
		assert.Equal(t, commits[0].Message, "commit")
		assert.Equal(t, commits[0].Headers, []repository.Header{{Key: "key2", Value: "value2"}})
	})
}

func TestRefs(t *testing.T) {
	t.Run("without filter", func(t *testing.T) {
		t.Parallel()

		repo, ref, err := repository.NewGitInMemoryRepository(true)
		assert.NoError(t, err)

		refs, err := repo.Refs()
		assert.NoError(t, err)

		expected := []repository.Ref{{
			Path: "refs/heads/main",
			Name: "main",
			Hash: ref.Hash,
		}}
		assert.Equal(t, expected, refs)
	})
	t.Run("with filter", func(t *testing.T) {
		t.Parallel()

		repo, _, err := repository.NewGitInMemoryRepository(true)
		assert.NoError(t, err)

		refs, err := repo.Refs(repository.WithPrefix("refs/custom"))
		assert.NoError(t, err)

		expected := []repository.Ref{}
		assert.Equal(t, expected, refs)
	})
}

func TestRemoteRefs(t *testing.T) {
	repo, _, err := repository.NewGitInMemoryRepository(true)
	assert.NoError(t, err)

	refs, err := repo.RemoteRefs("origin")
	assert.ErrorIs(t, err, repository.ErrRemoteNotFound)

	expected := []repository.Ref(nil)
	assert.Equal(t, expected, refs)
}

func TestContent(t *testing.T) {
	repo, _, err := repository.NewGitInMemoryRepository(true)
	assert.NoError(t, err)

	_, err = repo.Commit("refs/custom/main", "test", []byte("test"), "commit")
	assert.NoError(t, err)

	content, err := repo.Content("refs/custom/main", "test")
	assert.NoError(t, err)

	assert.Equal(t, []byte("test"), *content)
}

func TestDelete(t *testing.T) {
	repo, _, err := repository.NewGitInMemoryRepository(true)
	assert.NoError(t, err)

	_, err = repo.Commit("refs/custom/main", "test", []byte("test"), "commit")
	assert.NoError(t, err)

	refs, err := repo.Refs(repository.WithPrefix("refs/custom"))
	assert.NoError(t, err)
	assert.Equal(t, len(refs), 1)

	err = repo.Delete("refs/custom/main")
	assert.NoError(t, err)

	refs2, err := repo.Refs(repository.WithPrefix("refs/custom"))
	assert.NoError(t, err)
	assert.Equal(t, len(refs2), 0)
}

func TestAddRemote(t *testing.T) {
	repo, _, err := repository.NewGitInMemoryRepository(false)
	assert.NoError(t, err)

	err = repo.AddRemote("origin", "ssh://domain.tld")
	assert.NoError(t, err)
}

func TestPush(t *testing.T) {
	t.Run("without remote", func(t *testing.T) {
		t.Parallel()

		repo, _, err := repository.NewGitInMemoryRepository(true)
		assert.NoError(t, err)

		err = repo.Push("refs/heads/main", "origin", false)
		assert.ErrorIs(t, err, repository.ErrRemoteNotFound)
	})
	t.Run("with remote", func(t *testing.T) {
		t.Parallel()

		repo, ref, err := repository.NewGitInMemoryRepository(true)
		remote := addRemote(t, "origin", false, repo)
		assert.NoError(t, err)

		err = repo.Push("refs/heads/main", "origin", false)
		assert.NoError(t, err)

		refs, _ := remote.Refs()
		assert.Equal(t, *ref, refs[0])
	})
}

func TestMirror(t *testing.T) {
	t.Run("without remote", func(t *testing.T) {
		t.Parallel()

		repo, _, err := repository.NewGitInMemoryRepository(true)
		assert.NoError(t, err)

		err = repo.Mirror("refs/heads/main", "origin")
		assert.ErrorIs(t, err, repository.ErrRemoteNotFound)
	})
	t.Run("with remote", func(t *testing.T) {
		t.Parallel()

		repo, ref, err := repository.NewGitInMemoryRepository(true)
		remote := addRemote(t, "origin", false, repo)
		assert.NoError(t, err)

		err = repo.Mirror("refs/heads/main", "origin")
		assert.NoError(t, err)

		refs, _ := remote.Refs()
		assert.Equal(t, *ref, refs[0])
	})
}

func TestFetch(t *testing.T) {
	t.Run("without remote", func(t *testing.T) {
		repo, _, err := repository.NewGitInMemoryRepository(true)
		assert.NoError(t, err)

		err = repo.Fetch("refs/heads/main", "origin", false)
		assert.ErrorIs(t, err, repository.ErrRemoteNotFound)
	})
	t.Run("with remote", func(t *testing.T) {
		t.Parallel()

		repo, _, err := repository.NewGitInMemoryRepository(false)
		remote := addRemote(t, "origin", true, repo)
		assert.NoError(t, err)

		err = repo.Fetch("refs/heads/main", "origin", false)
		assert.NoError(t, err)

		remoteRefs, _ := remote.Refs()
		localRefs, _ := repo.Refs()

		assert.Equal(t, remoteRefs, localRefs)
	})
}
