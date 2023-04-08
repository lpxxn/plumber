//go:build wireinject
// +build wireinject

package log

import (
	"github.com/google/wire"
)

func initializeLog() Loger {
	wire.Build(logSet)
	return nil
}
