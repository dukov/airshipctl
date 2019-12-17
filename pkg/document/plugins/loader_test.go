package plugins_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"opendev.org/airship/airshipctl/pkg/document"
	"opendev.org/airship/airshipctl/testutil"
)

func TestLoader(t *testing.T) {
	t.Run("Try load non-existent plugin", func(t *testing.T) {
		_, err := document.NewBundle(testutil.SetupTestFs(t, "testdata/unknownplugin"), "/", "/")
		assert.Error(t, err)
	})
}
