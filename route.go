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

func (g *Group) Route(prefix string, route Route, middleware ...echo.MiddlewareFunc) {
	t := newGroup(g.d, g.Group.Group(prefix, middleware...))
	g.d.inject(route)
	route.Register(t)
}

func (g *Group) RouteFunc(prefix string, register RegisterFunc, middleware ...echo.MiddlewareFunc) {
	t := newGroup(g.d, g.Group.Group(prefix, middleware...))
	register(t)
}
