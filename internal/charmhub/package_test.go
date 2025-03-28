// Copyright 2020 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package charmhub

import (
	"io"
	"net/http"
	"net/url"
	"testing"

	jujutesting "github.com/juju/testing"
	jc "github.com/juju/testing/checkers"
	gc "gopkg.in/check.v1"

	corelogger "github.com/juju/juju/core/logger"
	"github.com/juju/juju/internal/charmhub/path"
	loggertesting "github.com/juju/juju/internal/logger/testing"
)

//go:generate go run go.uber.org/mock/mockgen -typed -package charmhub -destination client_mock_test.go github.com/juju/juju/internal/charmhub HTTPClient,RESTClient,FileSystem,ProgressBar

func Test(t *testing.T) {
	gc.TestingT(t)
}

type baseSuite struct {
	jujutesting.IsolationSuite

	logger corelogger.Logger
}

func (s *baseSuite) SetUpTest(c *gc.C) {
	s.logger = loggertesting.WrapCheckLog(c)
}

func MustParseURL(c *gc.C, path string) *url.URL {
	u, err := url.Parse(path)
	c.Assert(err, jc.ErrorIsNil)
	return u
}

func MustMakePath(c *gc.C, p string) path.Path {
	u := MustParseURL(c, p)
	return path.MakePath(u)
}

type nopCloser struct {
	io.Reader
}

func MakeNopCloser(r io.Reader) nopCloser {
	return nopCloser{
		Reader: r,
	}
}

func (nopCloser) Close() error { return nil }

func MakeContentTypeHeader(name string) http.Header {
	h := make(http.Header)
	h.Set("content-type", name)
	return h
}

func MustNewRequest(c *gc.C, path string) *http.Request {
	req, err := http.NewRequest("GET", path, nil)
	c.Assert(err, jc.ErrorIsNil)

	return req
}
