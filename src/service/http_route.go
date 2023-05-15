package service

type Router struct {
	routes map[string]*Route
}

type Route struct {
	Method string `json:"method"`
	Path   string `json:"path"`
}

func (r *Router) add(method, path string) *Route {

	route := &Route{
		Method: method,
		Path:   path,
	}
	r.routes[method+path] = route
	return route
}
