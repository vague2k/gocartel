package gocartel

import (
	"testing"

	"github.com/sanity-io/litter"
	"github.com/stretchr/testify/assert"
)

func TestAccount(t *testing.T) {
	c := TestClient()
	acc, err := c.Account()
	assert.NoError(t, err)
	assert.Equal(t, "Blackheaven Records", acc.StoreName)
	t.Log(litter.Sdump(acc))
}

func TestAccountByID(t *testing.T) {
	c := TestClient()
	internalID := InternalStoreID()
	acc, err := c.AccountByID(internalID)
	assert.NoError(t, err)
	assert.Equal(t, "Blackheaven Records", acc.StoreName)
	t.Log(litter.Sdump(acc))
}

func TestAccountByIDErrors(t *testing.T) {
	c := TestClient()
	acc, err := c.AccountByID("0")
	assert.Error(t, err)
	assert.Equal(t, err.Error(), "no account data found")
	assert.Nil(t, acc)
}
