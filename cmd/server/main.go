package main

import (
	"context"
	"os/signal"
	"syscall"

	"github.com/Onnywrite/lms-golang-24/internal/app"

	"dev.gaijin.team/go/golib/must"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	defer cancel()

	must.NoErr(app.New().Run(ctx))
}
