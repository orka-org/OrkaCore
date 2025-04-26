//go:build wireinject
// +build wireinject

// The build tag makes sure the stub is not built in the final build.

package main

import (
	"github.com/orka-org/orkacore/internal/biz"
	"github.com/orka-org/orkacore/internal/conf"
	"github.com/orka-org/orkacore/internal/data"
	"github.com/orka-org/orkacore/internal/server"
	"github.com/orka-org/orkacore/internal/service"
	"github.com/orka-org/orkacore/pkg"

	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
)

// wireApp init kratos application.
func wireApp(*conf.Bootstrap, *conf.Server, *conf.Data, log.Logger) (*kratos.App, func(), error) {
	panic(
		wire.Build(
			pkg.PkgProviderSet,
			data.DataProviderSet,
			biz.BizProviderSet,
			service.ServiceProviderSet,
			server.ServerProviderSet,
			newApp,
		),
	)
}
