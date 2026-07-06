package domain

type GitRepository interface {
	Init(path string) error
}
