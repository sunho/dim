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

func (g *Group) Add(method, path string, handler echo.HandlerFunc, middlewares ...Middleware) *echo.Route {
	for _, middleware := range middlewares {
		g.d.inject(middleware)
	}
	return g.Group.Add(method, path, handler, middlewaresToFuncs(middlewares)...)
}

func (g *Group) Any(path string, handler echo.HandlerFunc, middlewares ...Middleware) []*echo.Route {
	for _, middleware := range middlewares {
		g.d.inject(middleware)
	}
	return g.Group.Any(path, handler, middlewaresToFuncs(middlewares)...)
}

func (g *Group) Match(methods []string, path string, handler echo.HandlerFunc, middlewares ...Middleware) []*echo.Route {
	for _, middleware := range middlewares {
		g.d.inject(middleware)
	}
	return g.Group.Match(methods, path, handler, middlewaresToFuncs(middlewares)...)
}

func (g *Group) GET(path string, handler echo.HandlerFunc, middlewares ...Middleware) *echo.Route {
	for _, middleware := range middlewares {
		g.d.inject(middleware)
	}
	return g.Group.GET(path, handler, middlewaresToFuncs(middlewares)...)
}

func (g *Group) POST(path string, handler echo.HandlerFunc, middlewares ...Middleware) *echo.Route {
	for _, middleware := range middlewares {
		g.d.inject(middleware)
	}
	return g.Group.POST(path, handler, middlewaresToFuncs(middlewares)...)
}

func (g *Group) PUT(path string, handler echo.HandlerFunc, middlewares ...Middleware) *echo.Route {
	for _, middleware := range middlewares {
		g.d.inject(middleware)
	}
	return g.Group.PUT(path, handler, middlewaresToFuncs(middlewares)...)
}

func (g *Group) PATCH(path string, handler echo.HandlerFunc, middlewares ...Middleware) *echo.Route {
	for _, middleware := range middlewares {
		g.d.inject(middleware)
	}
	return g.Group.PATCH(path, handler, middlewaresToFuncs(middlewares)...)
}

func (g *Group) CONNECT(path string, handler echo.HandlerFunc, middlewares ...Middleware) *echo.Route {
	for _, middleware := range middlewares {
		g.d.inject(middleware)
	}
	return g.Group.CONNECT(path, handler, middlewaresToFuncs(middlewares)...)
}

func (g *Group) DELETE(path string, handler echo.HandlerFunc, middlewares ...Middleware) *echo.Route {
	for _, middleware := range middlewares {
		g.d.inject(middleware)
	}
	return g.Group.DELETE(path, handler, middlewaresToFuncs(middlewares)...)
}

func (g *Group) TRACE(path string, handler echo.HandlerFunc, middlewares ...Middleware) *echo.Route {
	for _, middleware := range middlewares {
		g.d.inject(middleware)
	}
	return g.Group.TRACE(path, handler, middlewaresToFuncs(middlewares)...)
}

func (g *Group) OPTIONS(path string, handler echo.HandlerFunc, middlewares ...Middleware) *echo.Route {
	for _, middleware := range middlewares {
		g.d.inject(middleware)
	}
	return g.Group.OPTIONS(path, handler, middlewaresToFuncs(middlewares)...)
}

func (g *Group) HEAD(path string, handler echo.HandlerFunc, middlewares ...Middleware) *echo.Route {
	for _, middleware := range middlewares {
		g.d.inject(middleware)
	}
	return g.Group.HEAD(path, handler, middlewaresToFuncs(middlewares)...)
}
