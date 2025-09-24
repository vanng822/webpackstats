package webpackstats

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestLoading(t *testing.T) {
	WebpackURLFuncMap("./data/webpack-stats.json")

	assert.EventuallyWithT(t, func(collect *assert.CollectT) {
		assert.True(collect, Get() != nil)
	}, 5*time.Second, 100*time.Millisecond, "Failed to load webpack stats")
}
