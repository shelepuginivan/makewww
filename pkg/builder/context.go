package builder

import (
	"github.com/shelepuginivan/makewww/pkg/config"
	"github.com/shelepuginivan/makewww/pkg/source"
)

type GlobalContext struct {
	Config    *config.Config
	Documents []source.Document
}
