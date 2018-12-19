package dim

import "github.com/labstack/echo"

type Middleware interface {
	Act(c echo.Context) error
}

func middlewareFunc(middleware Middleware) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if err := middleware.Act(c); err != nil {
				return err
			}
			return next(c)
		}
	}
}

func middlewaresToFuncs(middlewares []Middleware) []echo.MiddlewareFunc {
	out := make([]echo.MiddlewareFunc, 0, len(middlewares))
	for _, middleware := range middlewares {
		out = append(out, middlewareFunc(middleware))
	}
	return out
}
