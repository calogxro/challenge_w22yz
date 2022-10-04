package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreate(t *testing.T) {
	testStore.drop()
	// create
	id, _ := testStore.create(testAnswer)
	assert.NotNil(t, id)
}

func TestUpdate(t *testing.T) {
	testStore.drop()
	// update
	id, _ := testStore.update(testAnswer)
	assert.Nil(t, id)
}

func TestDelete(t *testing.T) {
	testStore.drop()
	// delete
	id, _ := testStore.delete(testAnswer.Key)
	assert.Nil(t, id)
}

// sequence allowed:
// create → delete → create → update
func TestCreateDeleteCreateUpdate(t *testing.T) {
	testStore.drop()
	// create
	id, _ := testStore.create(testAnswer)
	assert.NotNil(t, id)
	// delete
	id, _ = testStore.delete(testAnswer.Key)
	assert.NotNil(t, id)
	// create
	id, _ = testStore.create(testAnswer)
	assert.NotNil(t, id)
	// update
	id, _ = testStore.update(testAnswer)
	assert.NotNil(t, id)
}

// sequence allowed:
// create → update → delete → create → update
func TestCreateUpdateDeleteCreateUpdate(t *testing.T) {
	testStore.drop()
	// create
	id, _ := testStore.create(testAnswer)
	assert.NotNil(t, id)
	// update
	id, _ = testStore.update(testAnswer)
	assert.NotNil(t, id)
	// delete
	id, _ = testStore.delete(testAnswer.Key)
	assert.NotNil(t, id)
	// create
	id, _ = testStore.create(testAnswer)
	assert.NotNil(t, id)
	// update
	id, _ = testStore.update(testAnswer)
	assert.NotNil(t, id)
}

// sequence not allowed:
// create → delete → update
func TestCreateDeleteUpdate(t *testing.T) {
	testStore.drop()
	// create
	id, _ := testStore.create(testAnswer)
	assert.NotNil(t, id)
	// delete
	id, _ = testStore.delete(testAnswer.Key)
	assert.NotNil(t, id)
	// update
	id, err := testStore.update(testAnswer)
	assert.Nil(t, id)
	assert.NotNil(t, err)
}

// sequence not allowed:
// create → create
func TestCreateCreate(t *testing.T) {
	testStore.drop()
	// create
	id, _ := testStore.create(testAnswer)
	assert.NotNil(t, id)
	// create
	id, err := testStore.create(testAnswer)
	assert.Nil(t, id)
	assert.NotNil(t, err)
}
