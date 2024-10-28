package record

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_PathParser(t *testing.T) {

	tokens := ParsePath("aaa.bbb.ccc[0].ddd")

	if !assert.Equal(t, 5, len(tokens)) {
		return
	}

	assert.Equal(t, "aaa", tokens[0].Value)
	assert.Equal(t, "bbb", tokens[1].Value)
	assert.Equal(t, "ccc", tokens[2].Value)
	assert.Equal(t, "0", tokens[3].Value)
	assert.Equal(t, "ddd", tokens[4].Value)
}
