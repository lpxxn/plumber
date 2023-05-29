package httpredirect

import (
	"errors"
	"net"

	"github.com/lpxxn/plumber/src/log"
)

type Router struct {
	routes map[string]*Route

	ForwardHost func() *net.Conn
}

type Route struct {
	OriginPath  string `json:"path"`
	RouteParser routeParser
	ForwardHost func() *net.Conn
}

func (r *Router) Add(path string) (*Route, error) {
	route := &Route{
		OriginPath:  path,
		RouteParser: ParseRoute(path),
	}
	if _, ok := r.routes[path]; ok {
		log.Errorf("route already exists: %s", path)
		return nil, errors.New("route already exists")
	}
	r.routes[path] = route
	return route, nil
}

func (r *Router) MatchRoute(path string) *Route {
	for _, route := range r.routes {
		if RoutePatternMatch(route.OriginPath, path) {
			return route
		}
	}
	return nil
}
