package domain

type GitRepository interface {
	Init(path string) error
}

type FileRepository interface {
	List(vaultPath string) ([]string, error)
	Read(vaultPath, filename string) (string, error)
	Save(vaultPath, filename, content string) error
	Delete(vaultPath, filename string) error
}
