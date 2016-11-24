package api

import (
	"github.com/fgrehm/sls-web/src/core/compiler"
	"github.com/fgrehm/sls-web/src/core/models"
	"github.com/fgrehm/sls-web/src/core/renderer"
	"github.com/fgrehm/sls-web/src/core/solver"

	"github.com/facebookgo/grace/gracehttp"
	"github.com/labstack/echo"
	"github.com/labstack/echo/engine/standard"
	"github.com/labstack/echo/middleware"
)

type RunOpts struct {
	Port         string
	CompilerSvc  compiler.Service
	RendererSvc  renderer.Service
	SanModelsSvc models.Service
	SolverSvc    solver.Service
}

func Run(o RunOpts) error {
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	// FIXME: No need for this in production
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{}))

	// TODO
	// println("WARNING: Client assets were not found, make sure you are running the webpack dev server")
	e.File("/", "client/index.html")
	e.Static("/static", "client/static")

	registerCompilerEndpoints(e, o.CompilerSvc)
	registerRendererEndpoints(e, o.RendererSvc)
	registerSanModelsEndpoints(e, o.SanModelsSvc)
	registerSolverEndpoints(e, o.SolverSvc)

	std := standard.New(":" + o.Port)
	std.SetHandler(e)
	return gracehttp.Serve(std.Server)
}
