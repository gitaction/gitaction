package gitserver

import (
	"github.com/gitaction/src/domain/context/gitrepo"
	"io"
)

type GitSmart struct {
	repoName   string
	rpcName    string
	repository gitrepo.RepoRepository
}

type GitSmartTransferProtocol interface {
	GetRefsInfo(w io.Writer) error
	ReceivePack(r io.Reader, w io.Writer) error
}

func NewGitSmartUseCase(repoName string, rpcName string, repository gitrepo.RepoRepository) GitSmartTransferProtocol {
	return &GitSmart{
		repoName:   repoName,
		rpcName:    rpcName,
		repository: repository,
	}
}

func (usecase *GitSmart) GetRefsInfo(w io.Writer) error {
	rw := NewRepoIO(nil, w)
	repo := gitrepo.NewRepo(usecase.repoName, rw,  usecase.repository)
	return repo.GetRefsInfo(usecase.rpcName)
}

func (usecase *GitSmart) ReceivePack(r io.Reader, w io.Writer) error {
	rw := NewRepoIO(r, w)
	repo := gitrepo.NewRepo(usecase.repoName, rw,  usecase.repository)
	return repo.ReceivePack(usecase.rpcName)
}

type GitSmartReaderWriter struct {
	r io.Reader
	w io.Writer
}

func NewRepoIO(r io.Reader, w io.Writer) gitrepo.RepoIO {
	return &GitSmartReaderWriter{
		r: r,
		w: w,
	}
}

func (rw *GitSmartReaderWriter) GetReferenceUpdatesReader() io.Reader {
	return rw.r
}

func (rw *GitSmartReaderWriter) GetChunkWriter() io.Writer {
	return rw.w
}