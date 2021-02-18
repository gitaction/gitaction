package gitrepo

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
)

type RepoIO interface {
	GetReferenceUpdatesReader() io.Reader
	GetChunkWriter() io.Writer
}

type Repo struct {
	name       string
	path       string
	io         RepoIO
	repository RepoRepository
}

func NewRepo(repoName string, repoIO RepoIO, repository RepoRepository) *Repo {
	return &Repo{name: repoName, io: repoIO, repository: repository}
}

func (r *Repo) SameIdentityAs(e Repo) bool {
	return r.name == e.name
}

func (r *Repo) GetRefsInfo(rpc string) error {
	if err := r.setRepoDir(); err != nil {
		return err
	}
	if err := r.provisionRepo(); err != nil {
		return err
	}

	return r.runRpc(rpc, "--stateless-rpc", "--advertise-refs")
}

func (r *Repo) ReceivePack(rpc string) error {
	if err := r.setRepoDir(); err != nil {
		return err
	}
	if err := r.provisionRepo(); err != nil {
		return err
	}

	return r.runRpc(rpc, "--stateless-rpc")
}

func (r *Repo) runRpc(rpc string, arg ...string) error {
	cmd := exec.Command("git", strings.TrimPrefix(rpc, "git-"), strings.Join(arg, ","), r.path)
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	cmd.Env = append(os.Environ(),
		fmt.Sprintf("RECEIVE_APP=%s", "test123"),
	)
	reader, _ := cmd.StdoutPipe()
	if err := cmd.Start(); err != nil {
		return err
	}
	if rur := r.io.GetReferenceUpdatesReader(); rur != nil {
		writer, err := cmd.StdinPipe()
		if err != nil {
			return err
		}
		if _, err := io.Copy(writer, rur); err != nil {
			return err
		}
	}
	if cw := r.io.GetChunkWriter(); cw != nil {
		if _, err := io.Copy(cw, reader); err != nil {
			return err
		}
	}
	if err := cmd.Wait(); err != nil {
		return err
	}

	return nil
}

func (r *Repo) provisionRepo() error {
	if !r.repository.IsRepoExist(r.name) {
		if err := r.provisionEmptyRepo(); err != nil {
			return err
		}
	} else {
		// TODO: repo exist
		return nil
	}
	return nil
}

func (r *Repo) setRepoDir() error {
	path, err := r.getTmpDir()
	fmt.Printf("???? - %s\n", path)
	if err != nil {
		return err
	}
	r.path = path
	return nil
}

func (r *Repo) provisionEmptyRepo() error {
	if err := r.initRepo(); err != nil {
		return err
	}
	if err := r.disableGCAutoDetach(); err != nil {
		return err
	}
	if err := r.writePreReceiveHook(); err != nil {
		return err
	}
	return nil
}

func (r *Repo) getRepoPath() (path string, err error) {
	if !r.repository.IsRepoExist(r.name) {
		return r.getTmpDir()
	}

	return r.getTmpDir()
}

// TODO: remove tmp dir
func (r *Repo) getTmpDir() (path string, err error) {
	return ioutil.TempDir("", "gan-repo-"+r.name)
}

func (r *Repo) initRepo() error {
	cmd := exec.Command("git", "init")
	cmd.Dir = r.path
	return cmd.Run()
}

func (r *Repo) disableGCAutoDetach() error {
	cmd := exec.Command("git", "config", "--bool", "gc.autoDetach", "false")
	cmd.Dir = r.path
	return cmd.Run()
}

func (r *Repo) writePreReceiveHook() error {
	return ioutil.WriteFile(filepath.Join(r.path, ".git/hooks", "pre-receive"), preReceiveHook, 0755)
}
