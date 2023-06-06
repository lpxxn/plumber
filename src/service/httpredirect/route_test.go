package httpredirect

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRoutePatternMatch(t *testing.T) {
	r := NewRouter()
	_, err := r.Add("/test")
	assert.Nil(t, err)
	_, err = r.Add("/test/:id")
	assert.Nil(t, err)
	_, err = r.Add("api")
	assert.Nil(t, err)
	_, err = r.Add("/api/v1")
	_, err = r.Add("*")
	assert.Nil(t, err)

	route := r.MatchRoute("/test")
	assert.NotNil(t, route)
	assert.Equal(t, "/test", route.OriginPath)

	route = r.MatchRoute("/test/123")
	assert.NotNil(t, route)
	assert.Equal(t, "/test/:id", route.OriginPath)

	route = r.MatchRoute("/api/v1")
	assert.NotNil(t, route)
	assert.Equal(t, "/api/v1", route.OriginPath)

	route = r.MatchRoute("/api/v1/123")
	assert.NotNil(t, route)
	assert.Equal(t, "/*", route.OriginPath)
}
