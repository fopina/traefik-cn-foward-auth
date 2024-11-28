package cmd

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRunCommand(t *testing.T) {
	cmd := newRootCmd("v0.0.0")
	b := bytes.NewBufferString("")

	cmd.SetArgs([]string{"--version"})
	cmd.SetOut(b)

	err := cmd.Execute()
	out, _ := io.ReadAll(b)

	assert.NoError(t, err)
	assert.Equal(t, "traefik-cn-foward-auth version v0.0.0\n", string(out))
}

func TestHandlerRaw(t *testing.T) {
	rr := httptest.NewRecorder()
	options := defaultRunOptions()
	options.Raw = true
	handler := http.HandlerFunc(options.handler)

	req, err := http.NewRequest("GET", "/", nil)
	require.NoError(t, err)
	handler.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusForbidden, rr.Code)
	assert.Equal(t, "Forbidden", rr.Body.String())

	rr = httptest.NewRecorder()
	req.Header.Set(options.HeaderName, "X")
	req.Header.Set(options.AllowHeaderName, "X,Y")
	handler.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, "OK", rr.Body.String())

	rr = httptest.NewRecorder()
	req.Header.Set(options.HeaderName, "X")
	req.Header.Set(options.AllowHeaderName, "Y")
	handler.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusForbidden, rr.Code)
	assert.Equal(t, "Forbidden", rr.Body.String())
}

func TestHandler(t *testing.T) {
	rr := httptest.NewRecorder()
	options := defaultRunOptions()
	handler := http.HandlerFunc(options.handler)

	req, err := http.NewRequest("GET", "/", nil)
	require.NoError(t, err)
	handler.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusForbidden, rr.Code)
	assert.Equal(t, "Forbidden", rr.Body.String())

	rr = httptest.NewRecorder()
	req.Header.Set(options.HeaderName, url.QueryEscape("Subject=\"CN=me\""))
	req.Header.Set(options.AllowHeaderName, "me,you")
	handler.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, "OK", rr.Body.String())

	rr = httptest.NewRecorder()
	req.Header.Set(options.HeaderName, url.QueryEscape("Subject=\"CN=me\""))
	req.Header.Set(options.AllowHeaderName, "you")
	handler.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusForbidden, rr.Code)
	assert.Equal(t, "Forbidden", rr.Body.String())
}
