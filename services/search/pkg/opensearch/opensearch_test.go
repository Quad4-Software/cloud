package opensearch_test

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/Quad4-Software/cloud/services/search/internal/opensearchtest"
	"github.com/Quad4-Software/cloud/services/search/pkg/config"
)

var defaultConfig *config.Config

func TestMain(m *testing.M) {
	cfg, done, err := opensearchtest.SetupTests(context.Background())
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, "Failed to setup tests:", err)
		os.Exit(1)
		return
	}
	defaultConfig = cfg
	code := m.Run()
	done()
	os.Exit(code)
}
