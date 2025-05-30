// Copyright 2013 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package tools_test

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/juju/tc"

	"github.com/juju/juju/environs/tools"
)

func TestMarshalSuite(t *testing.T) {
	tc.Run(t, &marshalSuite{})
}

type marshalSuite struct {
	streamMetadata map[string][]*tools.ToolsMetadata
}

func (s *marshalSuite) SetUpTest(c *tc.C) {
	s.streamMetadata = map[string][]*tools.ToolsMetadata{
		"released": releasedToolMetadataForTesting,
		"proposed": proposedToolMetadataForTesting,
	}
}

func (s *marshalSuite) TestLargeNumber(c *tc.C) {
	metadata := map[string][]*tools.ToolsMetadata{
		"released": {
			{
				Release:  "ubuntu",
				Version:  "1.2.3.4",
				Arch:     "arm",
				Size:     9223372036854775807,
				Path:     "/somewhere/over/the/rainbow.tar.gz",
				FileType: "tar.gz",
			}},
	}
	_, _, products, err := tools.MarshalToolsMetadataJSON(metadata, time.Now())
	c.Assert(err, tc.ErrorIsNil)
	c.Assert(string(products["released"]), tc.Contains, `"size": 9223372036854775807`)
}

var expectedIndex = `{
    "index": {
        "com.ubuntu.juju:proposed:agents": {
            "updated": "Thu, 01 Jan 1970 00:00:00 +0000",
            "format": "products:1.0",
            "datatype": "content-download",
            "path": "streams/v1/com.ubuntu.juju-proposed-agents.json",
            "products": [
                "com.ubuntu.juju:ubuntu:arm64",
                "com.ubuntu.juju:ubuntu:ppc64el"            
            ]
        },
        "com.ubuntu.juju:released:agents": {
            "updated": "Thu, 01 Jan 1970 00:00:00 +0000",
            "format": "products:1.0",
            "datatype": "content-download",
            "path": "streams/v1/com.ubuntu.juju-released-agents.json",
            "products": [
                "com.ubuntu.juju:ubuntu:amd64",
                "com.ubuntu.juju:ubuntu:arm"
            ]
        }
    },
    "updated": "Thu, 01 Jan 1970 00:00:00 +0000",
    "format": "index:1.0"
}`

var expectedLegacyIndex = `{
    "index": {
        "com.ubuntu.juju:released:agents": {
            "updated": "Thu, 01 Jan 1970 00:00:00 +0000",
            "format": "products:1.0",
            "datatype": "content-download",
            "path": "streams/v1/com.ubuntu.juju-released-agents.json",
            "products": [
                "com.ubuntu.juju:ubuntu:amd64",
                "com.ubuntu.juju:ubuntu:arm"
            ]
        }
    },
    "updated": "Thu, 01 Jan 1970 00:00:00 +0000",
    "format": "index:1.0"
}`

var expectedReleasedProducts = `{
    "products": {
        "com.ubuntu.juju:ubuntu:amd64": {
            "version": "4.3.2.1",
            "arch": "amd64",
            "versions": {
                "19700101": {
                    "items": {
                        "4.3.2.1-ubuntu-amd64": {
                            "release": "ubuntu",
                            "version": "4.3.2.1",
                            "arch": "amd64",
                            "size": 0,
                            "path": "whatever.tar.gz",
                            "ftype": "tar.gz",
                            "sha256": "afb14e65c794464e378def12cbad6a96f9186d69"
                        }
                    }
                }
            }
        },
        "com.ubuntu.juju:ubuntu:arm": {
            "version": "1.2.3.4",
            "arch": "arm",
            "versions": {
                "19700101": {
                    "items": {
                        "1.2.3.4-ubuntu-arm": {
                            "release": "ubuntu",
                            "version": "1.2.3.4",
                            "arch": "arm",
                            "size": 9223372036854775807,
                            "path": "/somewhere/over/the/rainbow.tar.gz",
                            "ftype": "tar.gz",
                            "sha256": ""
                        }
                    }
                }
            }
        }
    },
    "updated": "Thu, 01 Jan 1970 00:00:00 +0000",
    "format": "products:1.0",
    "content_id": "com.ubuntu.juju:released:agents"
}`

var expectedProposedProducts = `{
    "products": {
        "com.ubuntu.juju:ubuntu:arm64": {
            "version": "1.2-beta1",
            "arch": "arm64",
            "versions": {
                "19700101": {
                    "items": {
                        "1.2-beta1-ubuntu-arm64": {
                            "release": "ubuntu",
                            "version": "1.2-beta1",
                            "arch": "arm64",
                            "size": 42,
                            "path": "gotham.tar.gz",
                            "ftype": "tar.gz",
                            "sha256": ""
                        }
                    }
                }
            }
        },
        "com.ubuntu.juju:ubuntu:ppc64el": {
            "version": "1.2-alpha1",
            "arch": "ppc64el",
            "versions": {
                "19700101": {
                    "items": {
                        "1.2-alpha1-ubuntu-ppc64el": {
                            "release": "ubuntu",
                            "version": "1.2-alpha1",
                            "arch": "ppc64el",
                            "size": 9223372036854775807,
                            "path": "/funkytown.tar.gz",
                            "ftype": "tar.gz",
                            "sha256": ""
                        }
                    }
                }
            }
        }
    },
    "updated": "Thu, 01 Jan 1970 00:00:00 +0000",
    "format": "products:1.0",
    "content_id": "com.ubuntu.juju:proposed:agents"
}`

var releasedToolMetadataForTesting = []*tools.ToolsMetadata{
	{
		Release:  "ubuntu",
		Version:  "1.2.3.4",
		Arch:     "arm",
		Size:     9223372036854775807,
		Path:     "/somewhere/over/the/rainbow.tar.gz",
		FileType: "tar.gz",
	},
	{
		Release:  "ubuntu",
		Version:  "4.3.2.1",
		Arch:     "amd64",
		Path:     "whatever.tar.gz",
		FileType: "tar.gz",
		SHA256:   "afb14e65c794464e378def12cbad6a96f9186d69",
	},
	{
		Release:  "xuanhuaceratops",
		Version:  "4.3.2.1",
		Arch:     "amd64",
		Size:     42,
		Path:     "dinodance.tar.gz",
		FileType: "tar.gz",
	},
	{
		Release:  "xuanhanosaurus",
		Version:  "5.4.3.2",
		Arch:     "amd64",
		Size:     42,
		Path:     "dinodisco.tar.gz",
		FileType: "tar.gz",
	},
}

var proposedToolMetadataForTesting = []*tools.ToolsMetadata{
	{
		Release:  "ubuntu",
		Version:  "1.2-alpha1",
		Arch:     "ppc64el",
		Size:     9223372036854775807,
		Path:     "/funkytown.tar.gz",
		FileType: "tar.gz",
	},
	{
		Release:  "ubuntu",
		Version:  "1.2-beta1",
		Arch:     "arm64",
		Size:     42,
		Path:     "gotham.tar.gz",
		FileType: "tar.gz",
	},
	{
		Release:  "xuanhuaceratops",
		Version:  "4.3.2.1",
		Arch:     "amd64",
		Size:     42,
		Path:     "dinodance.tar.gz",
		FileType: "tar.gz",
	},
	{
		Release:  "xuanhanosaurus",
		Version:  "5.4.3.2",
		Arch:     "amd64",
		Size:     42,
		Path:     "dinodisco.tar.gz",
		FileType: "tar.gz",
	},
}

func (s *marshalSuite) TestMarshalIndex(c *tc.C) {
	index, legacyIndex, err := tools.MarshalToolsMetadataIndexJSON(c.Context(), s.streamMetadata, time.Unix(0, 0).UTC())
	c.Assert(err, tc.ErrorIsNil)
	assertIndex(c, index, expectedIndex)
	assertIndex(c, legacyIndex, expectedLegacyIndex)
}

func assertIndex(c *tc.C, obtainedIndex []byte, expectedIndex string) {
	// Unmarshall into objects so an order independent comparison can be done.
	var obtained interface{}
	err := json.Unmarshal(obtainedIndex, &obtained)
	c.Assert(err, tc.ErrorIsNil)
	var expected interface{}
	err = json.Unmarshal([]byte(expectedIndex), &expected)
	c.Assert(err, tc.ErrorIsNil)
	c.Assert(obtained, tc.DeepEquals, expected)
}

func (s *marshalSuite) TestMarshalProducts(c *tc.C) {
	products, err := tools.MarshalToolsMetadataProductsJSON(s.streamMetadata, time.Unix(0, 0).UTC())
	c.Assert(err, tc.ErrorIsNil)
	assertProducts(c, products)
}

func assertProducts(c *tc.C, obtainedProducts map[string][]byte) {
	c.Assert(obtainedProducts, tc.HasLen, 2)
	c.Assert(string(obtainedProducts["released"]), tc.Equals, expectedReleasedProducts)
	c.Assert(string(obtainedProducts["proposed"]), tc.Equals, expectedProposedProducts)
}

func (s *marshalSuite) TestMarshal(c *tc.C) {
	index, legacyIndex, products, err := tools.MarshalToolsMetadataJSON(s.streamMetadata, time.Unix(0, 0).UTC())
	c.Assert(err, tc.ErrorIsNil)
	assertIndex(c, index, expectedIndex)
	assertIndex(c, legacyIndex, expectedLegacyIndex)
	assertProducts(c, products)
}

func (s *marshalSuite) TestMarshalNoReleaseStream(c *tc.C) {
	metadata := s.streamMetadata
	delete(metadata, "released")
	index, legacyIndex, products, err := tools.MarshalToolsMetadataJSON(s.streamMetadata, time.Unix(0, 0).UTC())
	c.Assert(err, tc.ErrorIsNil)
	c.Assert(legacyIndex, tc.IsNil)
	c.Assert(index, tc.NotNil)
	c.Assert(products, tc.NotNil)
}
