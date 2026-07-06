package adapters

import (
	"github.com/go-git/go-git/v6"
	"github.com/ogurechnikov/logbook/internal/domain"
)

var _ domain.GitRepository = (*GoGitRepository)(nil)

type GoGitRepository struct{}

func NewGoGitRepository() *GoGitRepository {
	return &GoGitRepository{}
}

func (g *GoGitRepository) Init(path string) error {
	_, err := git.PlainInit(path, false)
	return err
}
