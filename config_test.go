package gosql

import (
	"github.com/d3v-friends/go-pure/fnPanic"
	"testing"
)

func TestConfig(test *testing.T) {
	var cfg = fnPanic.OnPointer(Read(Path("./config.yaml")))

	test.Run("validate", func(t *testing.T) {
		var err error
		if err = cfg.Validate(); err != nil {
			t.Fatal(err)
		}
	})

}
