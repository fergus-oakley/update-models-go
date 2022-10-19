package main

import (
	"fmt"
	"github.com/go-git/go-git/v5"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestEnsureStagingEmpty(t *testing.T) {
	testCases := []struct {
		codeFile1 git.StatusCode
		codeFile2 git.StatusCode
		want      error
	}{
		{git.Unmodified, git.Unmodified, nil},
		{git.Untracked, git.Untracked, nil},
		{git.Modified, git.Modified, ErrStagedFiles},
		{git.Added, git.Added, ErrStagedFiles},
		{git.Deleted, git.Deleted, ErrStagedFiles},
		{git.Renamed, git.Renamed, ErrStagedFiles},
		{git.Copied, git.Copied, ErrStagedFiles},
		{git.UpdatedButUnmerged, git.UpdatedButUnmerged, ErrStagedFiles},

		{git.Unmodified, git.Modified, ErrStagedFiles},
		{git.Modified, git.Unmodified, ErrStagedFiles},
	}
	for _, tc := range testCases {
		t.Run(fmt.Sprintf("git.Status 1: %s, git.Status 2: %s, expect: %v", string(tc.codeFile1), string(tc.codeFile2), tc.want), func(t *testing.T) {
			statuses := map[string]*git.FileStatus{
				"filepath1": {
					Staging: tc.codeFile1,
				},
				"filepath2": {
					Staging: tc.codeFile2,
				},
			}

			err := ensureStagingEmpty(statuses)

			assert.ErrorIs(t, err, tc.want)
		})
	}
}

func TestRetrieveLatestVersion(t *testing.T) {
	t.Run("no tags found", func(t *testing.T) {
		_, got := retrieveLatestVersion([]string{})
		want := ErrVersionListEmpty

		assert.ErrorIs(t, got, want)
	})
	t.Run("tags formatted incorrectly", func(t *testing.T) {
		_, got := retrieveLatestVersion([]string{
			"v1.1.12",
			"POORLY FORMATTED VERSION",
		})
		want := ErrVersionFormat

		assert.ErrorIs(t, got, want)
	})
}

/*func TestRetrieveLatestRemoteGitTag(t *testing.T) {
	t.Run("no tags found", func(t *testing.T) {
		repo := &mockRepository{refs: nil}
		_, got := retrieveLatestRemoteGitTag(repo)
	})
}
*/
//type mockRepository struct {
//	refs []*plumbing.Reference
//}
//
//func (r *mockRepository) DeleteObject(hash plumbing.Hash) error {
//	return nil
//}
//func (r *mockRepository) Prune(opt git.PruneOptions) error {
//	return nil
//}
//func (r *mockRepository) Config() (*config.Config, error) {
//	return nil, nil
//}
//func (r *mockRepository) SetConfig(cfg *config.Config) error {
//	return nil
//}
//func (r *mockRepository) ConfigScoped(scope config.Scope) (*config.Config, error) {
//	return nil, nil
//}
//func (r *mockRepository) Remote(name string) (*git.Remote, error) {
//	return nil, nil
//}
//func (r *mockRepository) Remotes() ([]*git.Remote, error) {
//	return nil, nil
//}
//func (r *mockRepository) CreateRemote(c *config.RemoteConfig) (*git.Remote, error) {
//	return nil, nil
//}
//func (r *mockRepository) CreateRemoteAnonymous(c *config.RemoteConfig) (*git.Remote, error) {
//	return nil, nil
//}
//func (r *mockRepository) DeleteRemote(name string) error {
//	return nil
//}
//func (r *mockRepository) Branch(name string) (*config.Branch, error) {
//	return nil, nil
//}
//func (r *mockRepository) CreateBranch(c *config.Branch) error {
//	return nil
//}
//func (r *mockRepository) DeleteBranch(name string) error {
//	return nil
//}
//func (r *mockRepository) CreateTag(name string, hash plumbing.Hash, opts *git.CreateTagOptions) (*plumbing.Reference, error) {
//	return nil, nil
//}
//func (r *mockRepository) createTagObject(name string, hash plumbing.Hash, opts *git.CreateTagOptions) (plumbing.Hash, error) {
//	return plumbing.Hash{}, nil
//}
//func (r *mockRepository) buildTagSignature(tag *object.Tag, signKey *openpgp.Entity) (string, error) {
//	return "", nil
//}
//func (r *mockRepository) Tag(name string) (*plumbing.Reference, error) {
//	return nil, nil
//}
//func (r *mockRepository) DeleteTag(name string) error {
//	return nil
//}
//func (r *mockRepository) resolveToCommitHash(h plumbing.Hash) (plumbing.Hash, error) {
//	return plumbing.Hash{}, nil
//}
//func (r *mockRepository) clone(ctx context.Context, o *git.CloneOptions) error {
//	return nil
//}
//func (r *mockRepository) cloneRefSpec(o *git.CloneOptions) []config.RefSpec {
//	return nil
//}
//func (r *mockRepository) setIsBare(isBare bool) error {
//	return nil
//}
//func (r *mockRepository) updateRemoteConfigIfNeeded(o *git.CloneOptions, c *config.RemoteConfig, head *plumbing.Reference) error {
//	return nil
//}
//func (r *mockRepository) fetchAndUpdateReferences(ctx context.Context, o *git.FetchOptions, ref plumbing.ReferenceName) (*plumbing.Reference, error) {
//	return nil, nil
//}
//func (r *mockRepository) updateReferences(spec []config.RefSpec, resolvedRef *plumbing.Reference) (updated bool, err error) {
//	return false, nil
//}
//func (r *mockRepository) calculateRemoteHeadReference(spec []config.RefSpec, resolvedHead *plumbing.Reference) []*plumbing.Reference {
//	return nil
//}
//func (r *mockRepository) Fetch(o *git.FetchOptions) error {
//	return nil
//}
//func (r *mockRepository) FetchContext(ctx context.Context, o *git.FetchOptions) error {
//	return nil
//}
//func (r *mockRepository) Push(o *git.PushOptions) error {
//	return nil
//}
//func (r *mockRepository) PushContext(ctx context.Context, o *git.PushOptions) error {
//	return nil
//}
//func (r *mockRepository) Log(o *git.LogOptions) (object.CommitIter, error) {
//	return nil, nil
//}
//func (r *mockRepository) log(from plumbing.Hash, commitIterFunc func(*object.Commit) object.CommitIter) (object.CommitIter, error) {
//	return nil, nil
//}
//func (r *mockRepository) logAll(commitIterFunc func(*object.Commit) object.CommitIter) (object.CommitIter, error) {
//	return nil, nil
//}
//func (r *mockRepository) logWithFile(fileName string, commitIter object.CommitIter, checkParent bool) object.CommitIter {
//	return nil
//}
//func (r *mockRepository) logWithPathFilter(pathFilter func(string) bool, commitIter object.CommitIter, checkParent bool) object.CommitIter {
//	return nil
//}
//func (r *mockRepository) logWithLimit(commitIter object.CommitIter, limitOptions object.LogLimitOptions) object.CommitIter {
//	return nil
//}
//func (r *mockRepository) Tags() (storer.ReferenceIter, error) {
//	return nil, nil
//}
//func (r *mockRepository) Branches() (storer.ReferenceIter, error) {
//	return nil, nil
//}
//func (r *mockRepository) Notes() (storer.ReferenceIter, error) {
//	return nil, nil
//}
//func (r *mockRepository) TreeObject(h plumbing.Hash) (*object.Tree, error) {
//	return nil, nil
//}
//func (r *mockRepository) TreeObjects() (*object.TreeIter, error) {
//	return nil, nil
//}
//func (r *mockRepository) CommitObject(h plumbing.Hash) (*object.Commit, error) {
//	return nil, nil
//}
//func (r *mockRepository) CommitObjects() (object.CommitIter, error) {
//	return nil, nil
//}
//func (r *mockRepository) BlobObject(h plumbing.Hash) (*object.Blob, error) {
//	return nil, nil
//}
//func (r *mockRepository) BlobObjects() (*object.BlobIter, error) {
//	return nil, nil
//}
//func (r *mockRepository) TagObject(h plumbing.Hash) (*object.Tag, error) {
//	return nil, nil
//}
//func (r *mockRepository) TagObjects() (*object.TagIter, error) {
//	return nil, nil
//}
//func (r *mockRepository) Object(t plumbing.ObjectType, h plumbing.Hash) (object.Object, error) {
//	return nil, nil
//}
//func (r *mockRepository) Objects() (*object.ObjectIter, error) {
//	return nil, nil
//}
//func (r *mockRepository) Head() (*plumbing.Reference, error) {
//	return nil, nil
//}
//func (r *mockRepository) Reference(name plumbing.ReferenceName, resolved bool) (*plumbing.Reference, error) {
//	return nil, nil
//}
//func (r *mockRepository) References() (storer.ReferenceIter, error) {
//	return nil, nil
//}
//func (r *mockRepository) Worktree() (*git.Worktree, error) {
//	return nil, nil
//}
//func (r *mockRepository) ResolveRevision(rev plumbing.Revision) (*plumbing.Hash, error) {
//	return nil, nil
//}
//func (r *mockRepository) resolveHashPrefix(hashStr string) []plumbing.Hash {
//	return nil
//}
//func (r *mockRepository) RepackObjects(cfg *git.RepackConfig) (err error) {
//	return nil
//}
//func (r *mockRepository) createNewObjectPack(cfg *git.RepackConfig) (h plumbing.Hash, err error) {
//	return plumbing.Hash{}, nil
//}
