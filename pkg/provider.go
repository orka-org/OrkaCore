package pkg

import (
	"github.com/google/wire"
	"github.com/orka-org/orkacore/pkg/tokens"
)

var PkgProviderSet = wire.NewSet(tokens.NewTokenProvider)
