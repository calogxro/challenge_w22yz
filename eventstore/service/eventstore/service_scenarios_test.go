package eventstore

import (
	"testing"

	esmemdb "github.com/calogxro/qaservice/eventstore/repository/memory"

	"github.com/stretchr/testify/assert"
)

//var testAnswer = domain.Answer{Key: "name", Value: "John"}

// sequence allowed:
// create → delete → create → update
func TestCreateDeleteCreateUpdate(t *testing.T) {
	service := New(esmemdb.New())

	// Create
	_, err := service.CreateAnswer(testAnswer)
	assert.Nil(t, err)

	// Delete
	_, err = service.DeleteAnswer(testAnswer.Key)
	assert.Nil(t, err)

	// Create
	_, err = service.CreateAnswer(testAnswer)
	assert.Nil(t, err)

	// Update
	_, err = service.UpdateAnswer(testAnswer)
	assert.Nil(t, err)
}

// sequence allowed:
// create → update → delete → create → update
func TestCreateUpdateDeleteCreateUpdate(t *testing.T) {
	service := New(esmemdb.New())

	// Create
	_, err := service.CreateAnswer(testAnswer)
	assert.Nil(t, err)

	// Update
	_, err = service.UpdateAnswer(testAnswer)
	assert.Nil(t, err)

	// Delete
	_, err = service.DeleteAnswer(testAnswer.Key)
	assert.Nil(t, err)

	// Create
	_, err = service.CreateAnswer(testAnswer)
	assert.Nil(t, err)

	// Update
	_, err = service.UpdateAnswer(testAnswer)
	assert.Nil(t, err)
}

// sequence not allowed:
// create → delete → update
func TestCreateDeleteUpdate(t *testing.T) {
	service := New(esmemdb.New())

	// Create
	_, err := service.CreateAnswer(testAnswer)
	assert.Nil(t, err)

	// Delete
	_, err = service.DeleteAnswer(testAnswer.Key)
	assert.Nil(t, err)

	// Update
	_, err = service.UpdateAnswer(testAnswer)
	assert.NotNil(t, err)
}

// sequence not allowed:
// create → create
func TestCreateCreate(t *testing.T) {
	service := New(esmemdb.New())

	// Create
	_, err := service.CreateAnswer(testAnswer)
	assert.Nil(t, err)

	// Create
	_, err = service.CreateAnswer(testAnswer)
	assert.NotNil(t, err)
}

// sequence not allowed:
// create → delete → delete
func TestCreateDeleteDelete(t *testing.T) {
	service := New(esmemdb.New())

	// Create
	_, err := service.CreateAnswer(testAnswer)
	assert.Nil(t, err)

	// Delete
	_, err = service.DeleteAnswer(testAnswer.Key)
	assert.Nil(t, err)

	// Delete
	_, err = service.DeleteAnswer(testAnswer.Key)
	assert.NotNil(t, err)
}
