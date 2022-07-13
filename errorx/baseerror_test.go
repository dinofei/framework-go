package errorx_test

import (
	"testing"

	"github.com/dinofei/framework-go/errorx"
	"github.com/stretchr/testify/assert"
)

func TestNewBizError(t *testing.T) {
	be := errorx.NewBizError(101, "test error")
	assert.Equal(t, 101, be.Code)
	assert.Equal(t, "test error", be.Message)
	assert.NotEmpty(t, be.Error())
}

func TestIsBizError(t *testing.T) {
	be := errorx.NewBizError(101, "test error")
	assert.True(t, errorx.IsBizError(be))
}
