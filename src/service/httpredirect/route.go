package httpredirect

import (
	"errors"
	"net"
	"strings"

	"github.com/lpxxn/plumber/src/log"
)

type Router struct {
	routes    []*Route
	routesMap map[string]*Route

	ForwardConn func() (net.Conn, error)
}

func NewRouter() *Router {
	return &Router{
		routes:    []*Route{},
		routesMap: make(map[string]*Route),
	}
}

type Route struct {
	OriginPath  string `json:"path"`
	RouteParser routeParser
	ForwardConn func() (net.Conn, error)
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
	if _, ok := r.routesMap[pathRaw]; ok {
		log.Errorf("route already exists: %s", pathRaw)
		return nil, errors.New("route already exists")
	}
	r.routesMap[pathRaw] = route
	r.routes = append(r.routes, route)
	return route, nil
}

func (r *Router) MatchRoute(path string) *Route {
	for _, route := range r.routes {
		if route.Match(path) {
			return route
		}
	}
	return nil
}

func (r *Route) Match(path string) bool {
	if path == "" {
		path = "/"
	}
	patternPretty := r.OriginPath
	// Strict routing, remove trailing slashes
	if len(patternPretty) > 1 {
		patternPretty = strings.TrimRight(patternPretty, "/")
	}
	if r.OriginPath == "/" && path == "/" {
		return true
		// '*' wildcard matches any path
	} else if r.OriginPath == "/*" {
		return true
	}
	parser := r.RouteParser
	var ctxParams [maxParams]string
	// Does this route have parameters
	if len(parser.params) > 0 {
		if match := parser.getMatch(path, path, &ctxParams, false); match {
			return true
		}
	}
	// Check for a simple match
	patternPretty = RemoveEscapeChar(patternPretty)
	if len(patternPretty) == len(path) && patternPretty == path {
		return true
	}
	// No match
	return false
}
