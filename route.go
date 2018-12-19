package dim

import (
	"github.com/labstack/echo"
)

type Group struct {
	*echo.Group
	d *Dim
}

type RegisterFunc func(g *Group)

type Route interface {
	Register(g *Group)
}

func newGroup(d *Dim, g *echo.Group) *Group {
	return &Group{
		Group: g,
		d:     d,
	}
}

func (g *Group) Route(prefix string, route Route, middlewares ...Middleware) {
	t := newGroup(g.d, g.Group.Group(prefix, middlewaresToFuncs(middlewares)...))
	g.d.inject(route)
	route.Register(t)
}

func (g *Group) RouteFunc(prefix string, register RegisterFunc, middlewares ...Middleware) {
	t := newGroup(g.d, g.Group.Group(prefix, middlewaresToFuncs(middlewares)...))
	register(t)
}

func (g *Group) Use(middlewares ...Middleware) {
	for _, middleware := range middlewares {
		g.d.inject(middleware)
	}
	g.Group.Use(middlewaresToFuncs(middlewares)...)
}

func (g *Group) UseRaw(middlewares ...echo.MiddlewareFunc) {
	g.Group.Use(middlewares...)
}
