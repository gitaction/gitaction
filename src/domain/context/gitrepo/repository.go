package gitrepo

type RepoRepository interface {
	IsRepoExist(name string) bool
} 
