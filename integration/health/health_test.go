package integration

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"

	"go-api-newspaper/pkg"
)

func TestPing(t *testing.T) {
	endpoint := pkg.GetEndpoint("/health")
	res, err := http.Get(endpoint)
	assert.Nil(t, err)
	assert.Equal(t, 200, res.StatusCode)
}
