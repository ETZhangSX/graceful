package graceful

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestListenAndServe(t *testing.T) {
	err := ListenAndServe("bbb", nil)

	assert.Equal(t, "[graceful] listen tcp: address bbb: missing port in address", err.Error())
}
