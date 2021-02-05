package gitserver

import (
	"github.com/sunwei/gitaction/src/domain/context/gitrepo"
)

type GitSmart struct {
	name string
}

type GitSmartTransferProtocol interface {
 GetRefsInfo(name string) (info int64, err error)
} 

func NewGitSmartUseCase(repoName string) GitSmartTransferProtocol {
	return &GitSmart{repoName }
}

func (usecase *GitSmart) GetRefsInfo(name string) (info int64, err error) {
	repo := gitrepo.NewRepo(name)
	
}