package gitserver

import (
	"github.com/sunwei/gitaction/src/domain/context/gitrepo"
)

type GitSmart struct {
	repoName   string
	rpcName    string
	repository gitrepo.RepoRepository
}

type GitSmartTransferProtocol interface {
	GetRefsInfo() ([]byte, error)
}

func NewGitSmartUseCase(repoName string, rpcName string, repository gitrepo.RepoRepository) GitSmartTransferProtocol {
	return &GitSmart{
		repoName:   repoName,
		rpcName:    rpcName,
		repository: repository,
	}
}

func (usecase *GitSmart) GetRefsInfo() ([]byte, error) {
	repo := gitrepo.NewRepo(usecase.repoName, usecase.repository)
	return repo.GetRefsInfo(usecase.rpcName)
}
