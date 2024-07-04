package gorilla

import (
	"net/http"
	"testing"

	"go-admin-package/tests/common"

	"github.com/gavv/httpexpect"
)

func TestGorilla(t *testing.T) {
	common.ExtraTest(httpexpect.WithConfig(httpexpect.Config{
		Client: &http.Client{
			Transport: httpexpect.NewBinder(internalHandler()),
			Jar:       httpexpect.NewJar(),
		},
		Reporter: httpexpect.NewAssertReporter(t),
	}))
}
