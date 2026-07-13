package builder

import (
	"github.com/shelepuginivan/makewww/pkg/config"
	"github.com/shelepuginivan/makewww/pkg/resource"
)

type GlobalContext struct {
	Config    *config.Config
	Resources []resource.Resource
}
