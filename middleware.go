package dim

import "github.com/labstack/echo"

type Middleware interface {
}

type NormalMiddleware interface {
	Act(c echo.Context) error
}

type RawMiddleware interface {
	Act(next echo.HandlerFunc) echo.HandlerFunc
}

func middlewareFunc(middleware Middleware) echo.MiddlewareFunc {
	if mw, ok := middleware.(NormalMiddleware); ok {
		return func(next echo.HandlerFunc) echo.HandlerFunc {
			return func(c echo.Context) error {
				if err := mw.Act(c); err != nil {
					return err
				}
				return next(c)
			}
		}
	} else if mw, ok := middleware.(RawMiddleware); ok {
		return mw.Act
	}
	panic("Invalid middleware")
}

func middlewaresToFuncs(middlewares []Middleware) []echo.MiddlewareFunc {
	out := make([]echo.MiddlewareFunc, 0, len(middlewares))
	for _, middleware := range middlewares {
		out = append(out, middlewareFunc(middleware))
	}
	return out
}
