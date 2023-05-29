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

func (r *Router) Add(pathRaw string) (*Route, error) {
	if pathRaw == "" {
		pathRaw = "/"
	}
	// Path always start with a '/'
	if pathRaw[0] != '/' {
		pathRaw = "/" + pathRaw
	}
	route := &Route{
		OriginPath:  pathRaw,
		RouteParser: ParseRoute(pathRaw),
	}
	if _, ok := r.routes[pathRaw]; ok {
		log.Errorf("route already exists: %s", pathRaw)
		return nil, errors.New("route already exists")
	}
	r.routes[pathRaw] = route
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
