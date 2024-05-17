package binding

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/julienschmidt/httprouter"
	"github.com/stretchr/testify/assert"
)

func TestURIParser_Parse(t *testing.T) {
	req, err := http.NewRequest(http.MethodPost, "/name/wardonne/address/shanghai/age/10/valid/true", nil)
	assert.Nil(t, err)

	rr := httptest.NewRecorder()
	hr := httprouter.New()
	hr.Handle(http.MethodPost, "/name/:name/address/:address/age/:age/valid/:valid", func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		ctx := r.Context()
		ctx = context.WithValue(ctx, httprouter.ParamsKey, p)
		r = r.WithContext(ctx)
		container := &struct {
			Name    string `param:"name"`
			Address string `param:"address"`
			Age     int    `param:"age"`
			Valid   bool   `param:"valid"`
		}{}

		assert.Nil(t, URI(r, container))
		assert.Equal(t, "wardonne", container.Name)
		assert.Equal(t, "shanghai", container.Address)
		assert.Equal(t, 10, container.Age)
		assert.Equal(t, true, container.Valid)
	})
	hr.ServeHTTP(rr, req)
}
