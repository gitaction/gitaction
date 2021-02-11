package gitrepo

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
)

type Repo struct {
	name       string
	path       string
	repository RepoRepository
}

func NewRepo(repoName string, repository RepoRepository) *Repo {
	return &Repo{name: repoName, repository: repository}
}

func (r *Repo) SameIdentityAs(e Repo) bool {
	return r.name == e.name
}

func (r *Repo) GetRefsInfo(rpc string) ([]byte, error) {
	if err := r.setRepoDir(); err != nil {
		return nil, err
	}
	if err := r.provisionRepo(); err != nil {
		return nil, err
	}

	return r.runRpc(rpc)
}

func (r *Repo) runRpc(rpc string) ([]byte, error) {
	cmd := exec.Command("git", strings.TrimPrefix(rpc, "git-"), "--stateless-rpc", "--advertise-refs", r.path)
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	cmd.Env = append(os.Environ(),
		fmt.Sprintf("RECEIVE_APP=%s", "test123"),
	)
	output, err := cmd.Output() 
	if err != nil {
		return nil, err
	}
	//todo: get stdin ready
	//if err := cmd.Wait(); err != nil {
	//	return nil, err
	//}
	
	return output, nil
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
