package domain

type GitRepository interface {
	Init(path string) error
}

type FileRepository interface {
	List(vaultPath string) ([]string, error)
	Read(vaultPath, filename string) (string, error)
}
