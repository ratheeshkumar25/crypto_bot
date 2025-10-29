package main

import (
	"github.com/ratheeshkumar25/forex_bot/backend/pkg/di"
	"go.uber.org/fx"
)

func main() {
	fx.New(di.Module).Run()
}
