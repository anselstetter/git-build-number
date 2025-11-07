package buildnumber_test

import (
	"os"
	"testing"

	"github.com/anselstetter/git-build-number/internal/buildnumber"
	"github.com/anselstetter/git-build-number/internal/repository"
	"github.com/stretchr/testify/assert"
)

var (
	user  = "user"
	email = "email@domain.tld"
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

func TestGet(t *testing.T) {
	t.Run("create", func(t *testing.T) {
		t.Parallel()

		repo, _, _ := repository.NewGitInMemoryRepository(true)
		bn := buildnumber.New(repo)

		entry, err := bn.Get("test", user, email, true)
		assert.NoError(t, err)
		assert.Equal(t, int64(1), entry.Number)
	})
	t.Run("no create", func(t *testing.T) {
		t.Parallel()

		repo, _, _ := repository.NewGitInMemoryRepository(true)
		bn := buildnumber.New(repo)

		entry, err := bn.Get("test", user, email, false)
		assert.ErrorIs(t, err, buildnumber.ErrBuildNumberNotFound)
		assert.Nil(t, entry)
	})
}

func TestSet(t *testing.T) {
	t.Run("head", func(t *testing.T) {
		t.Parallel()

		repo, _, _ := repository.NewGitInMemoryRepository(true)
		bn := buildnumber.New(repo)

		entry, err := bn.Set("test", user, email, 123)
		assert.NoError(t, err)
		assert.Equal(t, int64(123), entry.Number)
	})
	t.Run("no head", func(t *testing.T) {
		t.Parallel()

		repo, _, _ := repository.NewGitInMemoryRepository(false)
		bn := buildnumber.New(repo)

		entry, err := bn.Set("test", user, email, 123)
		assert.Nil(t, entry)
		assert.ErrorIs(t, err, buildnumber.ErrNoHead)
	})
}

func TestInc(t *testing.T) {
	t.Run("no force", func(t *testing.T) {
		t.Parallel()

		repo, _, _ := repository.NewGitInMemoryRepository(true)
		bn := buildnumber.New(repo)

		entry, updated, err := bn.Inc("test", user, email, false)
		assert.NoError(t, err)
		assert.Equal(t, int64(1), entry.Number)
		assert.False(t, updated)
	})
	t.Run("force", func(t *testing.T) {
		t.Parallel()

		repo, _, _ := repository.NewGitInMemoryRepository(true)
		bn := buildnumber.New(repo)

		entry, updated, err := bn.Inc("test", user, email, true)
		assert.NoError(t, err)
		assert.Equal(t, int64(2), entry.Number)
		assert.True(t, updated)
	})
	t.Run("no head", func(t *testing.T) {
		t.Parallel()

		repo, _, _ := repository.NewGitInMemoryRepository(false)
		bn := buildnumber.New(repo)

		entry, updated, err := bn.Inc("test", user, email, true)
		assert.Nil(t, entry)
		assert.ErrorIs(t, err, buildnumber.ErrNoHead)
		assert.False(t, updated)
	})
	t.Run("from remote", func(t *testing.T) {
		t.Parallel()

		repo, _, _ := repository.NewGitInMemoryRepository(true)
		remote := addRemote(t, "origin", true, repo)

		bnLocal := buildnumber.New(repo)
		bnRemote := buildnumber.New(remote)

		_, _ = bnRemote.Set("test", user, email, 123)
		_ = bnLocal.Fetch("origin")

		entry, updated, err := bnLocal.Inc("test", user, email, false)
		assert.NoError(t, err)
		assert.Equal(t, int64(123), entry.Number)
		assert.False(t, updated)
	})
}

func TestDelete(t *testing.T) {
	repo, _, _ := repository.NewGitInMemoryRepository(true)
	bn := buildnumber.New(repo)

	_, _ = bn.Set("test", user, email, 123)
	entry, _ := bn.Set("test2", user, email, 123)

	err := bn.Delete("test")
	assert.NoError(t, err)

	ns, _ := bn.Namespaces()
	assert.Equal(t, []buildnumber.Namespace{{Name: "test2", Entry: *entry}}, ns)
}

func TestClear(t *testing.T) {
	repo, _, _ := repository.NewGitInMemoryRepository(true)
	bn := buildnumber.New(repo)

	_, _ = bn.Set("test", user, email, 123)
	_, _ = bn.Set("test2", user, email, 123)

	err := bn.Clear()
	assert.NoError(t, err)

	ns, _ := bn.Namespaces()
	assert.Equal(t, []buildnumber.Namespace{}, ns)
}

func TestNamespaces(t *testing.T) {
	repo, ref, _ := repository.NewGitInMemoryRepository(true)
	bn := buildnumber.New(repo)

	_, _ = bn.Set("test", user, email, 123)
	_, _ = bn.Set("other", user, email, 321)

	ns, err := bn.Namespaces()
	assert.NoError(t, err)

	expected := []buildnumber.Namespace{
		{
			Name: "test",
			Entry: buildnumber.Entry{
				Number: 123,
				Hash:   ref.Hash,
			},
		},
		{
			Name: "other",
			Entry: buildnumber.Entry{
				Number: 321,
				Hash:   ref.Hash,
			},
		},
	}
	assert.Equal(t, expected, ns)
}

func TestHash(t *testing.T) {
	t.Run("found", func(t *testing.T) {
		t.Parallel()

		repo, _, _ := repository.NewGitInMemoryRepository(true)
		bn := buildnumber.New(repo)

		entry, _ := bn.Set("test", user, email, 123)

		hash, err := bn.Hash("test", 123)
		assert.NoError(t, err)

		assert.Equal(t, entry.Hash, hash.Hash)
	})
	t.Run("not found", func(t *testing.T) {
		t.Parallel()

		repo, _, _ := repository.NewGitInMemoryRepository(true)
		bn := buildnumber.New(repo)

		_, _ = bn.Set("test", user, email, 123)

		hash, err := bn.Hash("missing", 321)
		assert.ErrorIs(t, err, buildnumber.ErrBuildNumberNotFound)
		assert.Nil(t, hash)
	})
	t.Run("no ref found", func(t *testing.T) {
		t.Parallel()

		repo, _, _ := repository.NewGitInMemoryRepository(true)
		bn := buildnumber.New(repo)

		hash, err := bn.Hash("test", 321)
		assert.ErrorIs(t, err, buildnumber.ErrBuildNumberNotFound)
		assert.Nil(t, hash)
	})
}

func TestPush(t *testing.T) {
	t.Run("without remote", func(t *testing.T) {
		t.Parallel()

		repo, _, _ := repository.NewGitInMemoryRepository(true)
		bn := buildnumber.New(repo)

		err := bn.Push("origin")
		assert.ErrorIs(t, err, repository.ErrRemoteNotFound)
	})
	t.Run("with remote", func(t *testing.T) {
		t.Parallel()

		repo, _, _ := repository.NewGitInMemoryRepository(true)
		remote := addRemote(t, "origin", false, repo)

		bnLocal := buildnumber.New(repo)
		bnRemote := buildnumber.New(remote)

		_, _ = bnLocal.Set("test", user, email, 123)

		err := bnLocal.Push("origin")
		assert.NoError(t, err)

		entry, _ := bnRemote.Get("test", user, email, false)
		assert.Equal(t, entry.Number, int64(123))
	})
}

func TestFetch(t *testing.T) {
	t.Run("without remote", func(t *testing.T) {
		t.Parallel()

		repo, _, _ := repository.NewGitInMemoryRepository(true)
		bn := buildnumber.New(repo)

		err := bn.Fetch("origin")
		assert.ErrorIs(t, err, repository.ErrRemoteNotFound)
	})
	t.Run("with remote", func(t *testing.T) {
		t.Parallel()

		repo, _, _ := repository.NewGitInMemoryRepository(true)
		remote := addRemote(t, "origin", true, repo)

		bnLocal := buildnumber.New(repo)
		bnRemote := buildnumber.New(remote)

		_, _ = bnRemote.Set("test", user, email, 123)

		err := bnLocal.Fetch("origin")
		assert.NoError(t, err)

		entry, _ := bnLocal.Get("test", user, email, false)
		assert.Equal(t, entry.Number, int64(123))
	})
}

func TestMirror(t *testing.T) {
	t.Run("without remote", func(t *testing.T) {
		repo, _, _ := repository.NewGitInMemoryRepository(true)
		bn := buildnumber.New(repo)

		err := bn.Mirror("origin")
		assert.ErrorIs(t, err, repository.ErrRemoteNotFound)
	})
	t.Run("with remote", func(t *testing.T) {
		repo, _, _ := repository.NewGitInMemoryRepository(true)
		remote := addRemote(t, "origin", true, repo)

		bnLocal := buildnumber.New(repo)
		bnRemote := buildnumber.New(remote)

		_, _ = bnRemote.Set("remote", user, email, 123)
		entry1, _ := bnLocal.Set("local-1", user, email, 123)
		entry2, _ := bnLocal.Set("local-2", user, email, 123)

		err := bnLocal.Mirror("origin")
		assert.NoError(t, err)

		ns, _ := bnRemote.Namespaces()
		assert.Equal(t, []buildnumber.Namespace{{Name: "local-1", Entry: *entry1}, {Name: "local-2", Entry: *entry2}}, ns)
	})
}

func TestMarshal(t *testing.T) {
	t.Run("valid entry", func(t *testing.T) {
		t.Parallel()

		entry := buildnumber.Entry{Number: 1, Hash: "a651ecb2072d58803f9fabdebecac5a569ac7e6b"}
		bytes, err := buildnumber.Marshal(entry)

		assert.NoError(t, err)
		assert.Equal(t, []byte("1 a651ecb2072d58803f9fabdebecac5a569ac7e6b"), bytes)
	})
	t.Run("zero value", func(t *testing.T) {
		t.Parallel()

		entry := buildnumber.Entry{Hash: "a651ecb2072d58803f9fabdebecac5a569ac7e6b"}
		bytes, err := buildnumber.Marshal(entry)

		assert.ErrorIs(t, err, buildnumber.ErrZeroBuildNumber)
		assert.Equal(t, []byte(nil), bytes)
	})
	t.Run("empty hash", func(t *testing.T) {
		t.Parallel()

		entry := buildnumber.Entry{Number: 1}
		bytes, err := buildnumber.Marshal(entry)

		assert.ErrorIs(t, err, buildnumber.ErrInvalidHash)
		assert.Equal(t, []byte(nil), bytes)
	})
}

func TestUnmarshal(t *testing.T) {
	t.Run("valid bytes", func(t *testing.T) {
		bytes := []byte("1   a651ecb2072d58803f9fabdebecac5a569ac7e6b  ")
		entry, err := buildnumber.Unmarshal(bytes)

		assert.NoError(t, err)
		assert.Equal(t, &buildnumber.Entry{Number: 1, Hash: "a651ecb2072d58803f9fabdebecac5a569ac7e6b"}, entry)
	})
	t.Run("invalid format", func(t *testing.T) {
		bytes := []byte("1")
		entry, err := buildnumber.Unmarshal(bytes)

		assert.ErrorIs(t, err, buildnumber.ErrInvalidFormat)
		assert.Nil(t, entry)
	})
	t.Run("invalid build number", func(t *testing.T) {
		bytes := []byte("abc a651ecb2072d58803f9fabdebecac5a569ac7e6b")
		entry, err := buildnumber.Unmarshal(bytes)

		assert.ErrorIs(t, err, buildnumber.ErrInvalidBuildNumber)
		assert.Nil(t, entry)
	})
}
