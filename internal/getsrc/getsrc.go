package getsrc

import (
	"log"
	"sort"
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/plumbing/storer"
)

type Repo struct {
	Repo *git.Repository
	Name string
}

func NewRepo(name string, path string) (*Repo, error) {
	repo, err := git.PlainOpen(path)
	if err != nil {
		return nil, err
	}
	return &Repo{Repo: repo, Name: name}, nil

}

func (rr *Repo) ToList(it storer.ReferenceIter) []*plumbing.Reference {
	ret := []*plumbing.Reference{}
	ref, err := it.Next()
	for err == nil {
		ret = append(ret, ref)
		ref, err = it.Next()
	}
	if err != nil {
		log.Println(err)
	}
	return ret
}

func (rr *Repo) ToListObject(it *object.ObjectIter) []object.Object {
	ret := []object.Object{}
	ref, err := it.Next()
	for err == nil {
		ret = append(ret, ref)
		ref, err = it.Next()
	}
	if err != nil {
		log.Println(err)
	}
	return ret
}

func (rr *Repo) ToListFiles(it *object.FileIter) []*object.File {
	ret := []*object.File{}
	ref, err := it.Next()
	for err == nil {
		ret = append(ret, ref)
		ref, err = it.Next()
		// ref.Entries
	}
	if err != nil {
		log.Println(err)
	}
	return ret
}

func (rr *Repo) ToListTree(it *object.TreeIter) []*object.Tree {
	ret := []*object.Tree{}
	ref, err := it.Next()
	for err == nil {
		ret = append(ret, ref)
		ref, err = it.Next()
		// ref.Entries
	}
	if err != nil {
		log.Println(err)
	}
	return ret
}

func (rr *Repo) CommitCount() int {
	cnt := 0

	if c, err := rr.Repo.Log(&git.LogOptions{}); err == nil {
		c.ForEach(func(cc *object.Commit) error {
			cnt++
			return nil
		})
	}

	// if c, err := rr.Repo.CommitObjects(); err == nil {
	// 	c.ForEach(func(cc *object.Commit) error {
	// 		cnt++
	// 		return nil
	// 	})
	// }
	return cnt
}

type ByName []object.TreeEntry

func (a ByName) Len() int      { return len(a) }
func (a ByName) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a ByName) Less(i, j int) bool {
	if a[i].Mode.IsFile() == a[j].Mode.IsFile() {
		return strings.Compare(a[i].Name, a[j].Name) < 0
	} else {
		return !a[i].Mode.IsFile()
	}
}

func (rr *Repo) Files() ([]object.TreeEntry, error) {
	var err error
	if h, err := rr.Repo.Head(); err == nil {
		if c, err := rr.Repo.CommitObject(h.Hash()); err == nil {
			if tree, err := c.Tree(); err == nil {
				ret := tree.Entries
				sort.Sort(ByName(ret))
				return ret, nil
			}
		}
	}
	return nil, err
}
