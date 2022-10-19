package main

import (
	"fmt"
	"github.com/go-git/go-git/v5"
	"github.com/spf13/cast"
	"log"
	"os"
	"regexp"
	"sort"
	"strings"
)

// TODO must abstract all git bump versioning functionality out and into its own file, and include ability to pass minor and major flags
type Error string

func (e Error) Error() string {
	return string(e)
}

const (
	ErrStagedFiles      = Error("Error: you already have staged files")
	ErrVersionListEmpty = Error("error: version list empty")
	ErrVersionFormat    = Error("error: version list contains incorrectly formatted versions")
	ErrNoTagsFound      = Error("error: no tags found")
)

var (
	newTag = ""
)

/*type Repository interface {
	DeleteObject(hash plumbing.Hash) error
	Prune(opt git.PruneOptions) error
	Config() (*config.Config, error)
	SetConfig(cfg *config.Config) error
	ConfigScoped(scope config.Scope) (*config.Config, error)
	Remote(name string) (*git.Remote, error)
	Remotes() ([]*git.Remote, error)
	CreateRemote(c *config.RemoteConfig) (*git.Remote, error)
	CreateRemoteAnonymous(c *config.RemoteConfig) (*git.Remote, error)
	DeleteRemote(name string) error
	Branch(name string) (*config.Branch, error)
	CreateBranch(c *config.Branch) error
	DeleteBranch(name string) error
	CreateTag(name string, hash plumbing.Hash, opts *git.CreateTagOptions) (*plumbing.Reference, error)
	createTagObject(name string, hash plumbing.Hash, opts *git.CreateTagOptions) (plumbing.Hash, error)
	buildTagSignature(tag *object.Tag, signKey *openpgp.Entity) (string, error)
	Tag(name string) (*plumbing.Reference, error)
	DeleteTag(name string) error
	resolveToCommitHash(h plumbing.Hash) (plumbing.Hash, error)
	clone(ctx context.Context, o *git.CloneOptions) error
	cloneRefSpec(o *git.CloneOptions) []config.RefSpec
	setIsBare(isBare bool) error
	updateRemoteConfigIfNeeded(o *git.CloneOptions, c *config.RemoteConfig, head *plumbing.Reference) error
	fetchAndUpdateReferences(ctx context.Context, o *git.FetchOptions, ref plumbing.ReferenceName) (*plumbing.Reference, error)
	updateReferences(spec []config.RefSpec, resolvedRef *plumbing.Reference) (updated bool, err error)
	calculateRemoteHeadReference(spec []config.RefSpec, resolvedHead *plumbing.Reference) []*plumbing.Reference
	Fetch(o *git.FetchOptions) error
	FetchContext(ctx context.Context, o *git.FetchOptions) error
	Push(o *git.PushOptions) error
	PushContext(ctx context.Context, o *git.PushOptions) error
	Log(o *git.LogOptions) (object.CommitIter, error)
	log(from plumbing.Hash, commitIterFunc func(*object.Commit) object.CommitIter) (object.CommitIter, error)
	logAll(commitIterFunc func(*object.Commit) object.CommitIter) (object.CommitIter, error)
	logWithFile(fileName string, commitIter object.CommitIter, checkParent bool) object.CommitIter
	logWithPathFilter(pathFilter func(string) bool, commitIter object.CommitIter, checkParent bool) object.CommitIter
	logWithLimit(commitIter object.CommitIter, limitOptions object.LogLimitOptions) object.CommitIter
	Tags() (storer.ReferenceIter, error)
	Branches() (storer.ReferenceIter, error)
	Notes() (storer.ReferenceIter, error)
	TreeObject(h plumbing.Hash) (*object.Tree, error)
	TreeObjects() (*object.TreeIter, error)
	CommitObject(h plumbing.Hash) (*object.Commit, error)
	CommitObjects() (object.CommitIter, error)
	BlobObject(h plumbing.Hash) (*object.Blob, error)
	BlobObjects() (*object.BlobIter, error)
	TagObject(h plumbing.Hash) (*object.Tag, error)
	TagObjects() (*object.TagIter, error)
	Object(t plumbing.ObjectType, h plumbing.Hash) (object.Object, error)
	Objects() (*object.ObjectIter, error)
	Head() (*plumbing.Reference, error)
	Reference(name plumbing.ReferenceName, resolved bool) (*plumbing.Reference, error)
	References() (storer.ReferenceIter, error)
	Worktree() (*git.Worktree, error)
	ResolveRevision(rev plumbing.Revision) (*plumbing.Hash, error)
	resolveHashPrefix(hashStr string) []plumbing.Hash
	RepackObjects(cfg *git.RepackConfig) (err error)
	createNewObjectPack(cfg *git.RepackConfig) (h plumbing.Hash, err error)
}

*/func ensureStagingEmpty(statuses git.Status) error {
	for _, s := range statuses {
		if s.Staging == git.Unmodified || s.Staging == git.Untracked {
			continue
		}
		return ErrStagedFiles
	}
	return nil
}

/*	vList := []string{
		"v0.0.1",
		"v0.1.0",
		"v0.1.1",
		"v0.2.0",
		"v1.0.0",
		"v1.0.1",
		"v1.1.0",
		"v1.2.0",
		"v1.3.0",
		"v1.3.1",
		"v1.3.2",
		"v1.4.0",
		"v1.4.1",
		"v1.4.2",
		"v1.4.3",
		"v1.4.4",
		"v1.4.5",
		"v1.4.6",
		"v1.5.0",
		"v1.5.1",
		"v1.5.10",
		"v1.5.11",
		"v1.5.12",
		"v1.5.13",
		"v1.5.14",
		"v1.5.15",
		"v1.5.16",
		"v1.5.17",
		"v1.5.18",
		"v1.5.19",
		"v1.5.2",
		"v1.5.20",
		"v1.5.21",
		"v1.5.22",
		"v1.5.23",
		"v1.5.24",
		"v1.5.25",
		"v1.5.26",
		"v1.5.27",
		"v1.5.28",
		"v1.5.29",
		"v1.5.3",
		"v1.5.4",
		"v1.5.5",
		"v1.5.6",
		"v1.5.7",
		"v1.5.8",
		"v1.5.9",
		"v1.6.0",
		"v1.6.1",
		"v1.6.2",
		"v1.6.3",
		"v1.6.4",
		"v1.6.5",
		"v1.6.6",
		"v1.7.0",
		"v1.7.1",
		"v1.7.10",
		"v1.7.11",
		"v1.7.12",
		"v1.7.13",
		"v1.7.14",
		"v1.7.15",
		"v1.7.16",
		"v1.7.17",
		"v1.7.18",
		"v1.7.19",
		"v1.7.2",
		"v1.7.20",
		"v1.7.21",
		"v1.7.22",
		"v1.7.23",
		"v1.7.24",
		"v1.7.25",
		"v1.7.26",
		"v1.7.27",
		"v1.7.28",
		"v1.7.3",
		"v1.7.4",
		"v1.7.5",
		"v1.7.6",
		"v1.7.7",
		"v1.7.8",
		"v1.7.9",
		"v1.8.0",
		"v1.8.1",
		"v1.8.10",
		"v1.8.11",
		"v1.8.12",
		"v1.8.13",
		"v1.8.14",
		"v1.8.15",
		"v1.8.16",
		"v1.8.17",
		"v1.8.18",
		"v1.8.19",
		"v1.8.2",
		"v1.8.20",
		"v1.8.21",
		"v1.8.22",
		"v1.8.23",
		"v1.8.24",
		"v1.8.25",
		"v1.8.26",
		"v1.8.27",
		"v1.8.28",
		"v1.8.29",
		"v1.8.3",
		"v1.8.30",
		"v1.8.31",
		"v1.8.32",
		"v1.8.33",
		"v1.8.34",
		"v1.8.35",
		"v1.8.36",
		"v1.8.37",
		"v1.8.38",
		"v1.8.39",
		"v1.8.4",
		"v1.8.40",
		"v1.8.41",
		"v1.8.42",
		"v1.8.43",
		"v1.8.44",
		"v1.8.45",
		"v1.8.46",
		"v1.8.47",
		"v1.8.48",
		"v1.8.49",
		"v1.8.5",
		"v1.8.50",
		"v1.8.51",
		"v1.8.52",
		"v1.8.53",
		"v1.8.54",
		"v1.8.55",
		"v1.8.6",
		"v1.8.7",
		"v1.8.8",
	}
*/

func main() {
	repo, err := git.PlainOpen("../..")
	if err != nil {
		panic(err)
	}

	if err := tag(repo, "v1.8.56"); err != nil {
		log.Fatal(err)
	}

	return

	worktree, err := repo.Worktree()
	if err != nil {
		panic(err)
	}

	status, err := worktree.Status()
	if err != nil {
		panic(err)
	}

	latestTag, err := retrieveLatestRemoteGitTag(repo)
	if err != nil {
		if err != ErrNoTagsFound {
			log.Fatal(err)
		}
		fmt.Println("no tag found")
		newTag = "v0.0.1"
		//tag(repo, newTag)
		os.Exit(0)
	}
	fmt.Println("last tag: ", latestTag)

	v := versionToIntSlice(latestTag)
	newTag = fmt.Sprintf("v%d.%d.%d", v[0], v[1], v[2]+1)

	/*	if err := tag(repo, newTag); err != nil {
			log.Fatal(err)
		}
	*/
	return

	if err := ensureStagingEmpty(status); err != nil {
		log.Fatal(err)
	}

	if err := worktree.AddGlob("*pkg/model*"); err != nil {
		if err == git.ErrGlobNoMatches {
			log.Fatal("Error: No changes to pkg/model to update.")
		}
		panic(err)
	}

	if _, err := worktree.Commit("chore: updated models", &git.CommitOptions{}); err != nil {
		panic(err)
	}

	if err := repo.Push(&git.PushOptions{}); err != nil {
		panic(err)
	}

	// NOTE a tag is an arbitrary identifier for a commit (that is more readable than
	// the hash that git generates automatically.). Go uses git tags - it piggybacks
	// on some git features like this. There are minor and major parts of the version
	// number generally to indicate small or breaking changes respectively. Here,
	// after the git push, we then increment the version locally using bump, and as
	// part of the bump we push the new version tag to git. Matt is retroactively
	// creating a new tag version post commit push, he then pushes THE TAG up to the
	// cloud as well and retroactively tags the commit with it.

}

// retrieveLatestRemoteGitTag retrieves the latest tag available on the remote repository.
func retrieveLatestRemoteGitTag(repo *git.Repository) (string, error) {
	remote, err := repo.Remote("origin")
	if err != nil {
		return "", err
	}
	refs, err := remote.List(&git.ListOptions{})
	if err != nil {
		return "", err
	}

	if len(refs) == 0 {
		return "", ErrNoTagsFound
	}

	var versionList []string
	for _, r := range refs {
		version := r.Name().Short()
		versionFormat := regexp.MustCompile("v[0-9]*\\.[0-9]*\\.[0-9]*")
		if versionFormat.MatchString(version) {
			versionList = append(versionList, version)
		}
	}

	latestTag, err := retrieveLatestVersion(versionList)
	if err != nil {
		return "", err
	}

	return latestTag, nil
	/*	iter, err := repo.Tags()
		if err != nil {
			panic(err)
		}

		var versionList []string
		if err := iter.ForEach(func(ref *plumbing.Reference) error {
			version := ref.Name().Short()
			versionList = append(versionList, version)
			return nil
		}); err != nil {
			// Handle outer iterator error
			panic(err)
		}
	*/
}

// retrieveLatestVersion applies a natural sort to a list of versions (of format "v1.1.1") and returns the latest version.
func retrieveLatestVersion(versionList []string) (string, error) {
	var formattedVersionList [][]int
	var majorElements []int
	var minorElements []int
	var bugFixElements []int

	if len(versionList) == 0 {
		return "", ErrVersionListEmpty
	}

	expr := regexp.MustCompile("v[0-9]*\\.[0-9]*\\.[0-9]*")
	for _, v := range versionList {
		if !expr.MatchString(v) {
			return "", ErrVersionFormat
		}
	}

	for _, v := range versionList {
		versionIntSlice := versionToIntSlice(v)
		formattedVersionList = append(formattedVersionList, versionIntSlice)
	}

	for _, v := range formattedVersionList {
		majorElements = append(majorElements, v[0])
	}
	sort.Ints(majorElements)
	mostRecentMajorElement := majorElements[len(majorElements)-1]

	for _, v := range formattedVersionList {
		if v[0] == mostRecentMajorElement {
			minorElements = append(minorElements, v[1])
		}
	}
	sort.Ints(minorElements)
	mostRecentMinorElement := minorElements[len(minorElements)-1]

	for _, v := range formattedVersionList {
		if v[0] == mostRecentMajorElement && v[1] == mostRecentMinorElement {
			bugFixElements = append(bugFixElements, v[2])
		}
	}
	sort.Ints(bugFixElements)
	mostRecentBugFixElement := bugFixElements[len(bugFixElements)-1]

	mostRecentVersion := fmt.Sprintf("v%d.%d.%d", mostRecentMajorElement, mostRecentMinorElement, mostRecentBugFixElement)

	return mostRecentVersion, nil
}

// versionToIntSlice takes a version string of format "v1.4.12" and converts it to []int{"1", "4", "12"}
func versionToIntSlice(v string) []int {
	v = strings.TrimPrefix(v, "v")

	var versionElementsInt []int

	es := strings.Split(v, ".")
	versionElementsInt = cast.ToIntSlice(es)
	return versionElementsInt
}

// tag creates a new tag and pushes all local tags to the remote repository.
func tag(repo *git.Repository, newTag string) error {
	fmt.Println("new tag: ", newTag)

	h, err := repo.Head()
	if err != nil {
		log.Fatal("get HEAD error: %s", err)
		return err
	}

	if _, err := repo.CreateTag(newTag, h.Hash(), nil); err != nil {
		return err
	}

	/*	if err := repo.Push(&git.PushOptions{RemoteName: "origin", RefSpecs: []config.RefSpec{config.RefSpec("refs/tags/*:refs/tags/*")}, Progress: os.Stdout}); err != nil {
			fmt.Println("claims origin up to date")
			return err
		}
	*/return nil
}

/*func setTag(r *git.Repository, tag string, tagger *object.Signature) (bool, error) {
	if tagExists(tag, r) {
		log.Infof("tag %s already exists", tag)
		return false, nil
	}
	log.Infof("Set tag %s", tag)
	h, err := r.Head()
	if err != nil {
		log.Errorf("get HEAD error: %s", err)
		return false, err
	}
	_, err = r.CreateTag(tag, h.Hash(), &git.CreateTagOptions{
		Tagger:  tagger,
		Message: tag,
	})
	if err != nil {
		log.Errorf("create tag error: %s", err)
		return false, err
	}
	return true, nil
}
*/