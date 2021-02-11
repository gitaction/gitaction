package persistece

import "fmt"

type RepoRepositoryImpl struct {}

func NewRepoRepositoryImpl() *RepoRepositoryImpl {
	return &RepoRepositoryImpl{}
}

func (repo *RepoRepositoryImpl)IsRepoExist(name string) bool {
	fmt.Println(name)
	return false
}