//go:build wireinject
// +build wireinject

package wire

import (
	"go_code/simplek8s/core/application/dao"
	"go_code/simplek8s/core/application/handler"
	"go_code/simplek8s/core/application/service"
	"go_code/simplek8s/internal/database"
	"go_code/simplek8s/server"
	"net/http"

	"github.com/google/wire"
)

func InitializeRouter() (http.Handler, error) {
	wire.Build(
		database.NewDB,
		dao.NewClusterDao,
		service.NewClusterService,
		handler.NewClusterHandler,
		server.NewRouter,
	)
	return nil, nil
}
