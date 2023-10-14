package trier

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func passOrFail(args ...any) error {
	if len(args) != 0 {
		return errors.New("failed passOrFail")
	}
	return nil
}

func failIfString(args ...any) error {
	var err error

	switch args[0].(type) {
	case string:
		err = errors.New("failedIfString")
	}
	return err
}

func TestNewTrier(t *testing.T) {
	// Act
	tr := NewTrier()

	// Assert
	assert.NotNil(t, tr)
}

func TestTrierTry(t *testing.T) {
	// Arrange
	tr := NewTrier()

	// Act
	tr.Try(passOrFail)

	// Assert
	x := *tr.err
	assert.Equal(t, "", x.Error())
}

func TestTrierTryError(t *testing.T) {
	// Arrange
	tr := NewTrier()

	// Act
	tr.Try(passOrFail, true)

	// Assert
	x := *tr.err
	assert.Equal(t, "failed passOrFail", x.Error())
}

func TestTrierTryErrorChainedTries(t *testing.T) {
	// Arrange
	tr := NewTrier()

	// Act
	tr.Try(passOrFail).
		Try(failIfString, 0).
		Try(passOrFail, true).
		Try(failIfString, "hi")

	// Assert
	x := *tr.err
	assert.Equal(t, "failed passOrFail", x.Error())
}

func TestTrierTryErrorWrapped(t *testing.T) {
	// Arrange
	tr := NewTrier()

	// Act
	tr.Try(passOrFail).
		Try(failIfString, 0).
		Try(passOrFail, true).
		TryJoin(failIfString, "hi")

	// Assert
	x := *tr.err
	assert.Equal(t, "failedIfString\nfailed passOrFail", x.Error())
}

func TestTrierErr(t *testing.T) {
	// Arrange
	tr := NewTrier()

	// Act
	tr.Try(passOrFail, true)

	// Assert
	assert.Equal(t, "failed passOrFail", tr.Err().Error())
}

func TestTrierErrChainedTries(t *testing.T) {
	// Arrange
	tr := NewTrier()

	// Act
	tr.Try(passOrFail).
		Try(failIfString, 0).
		Try(passOrFail, true).
		Try(failIfString, "hi")

	// Assert
	assert.Equal(t, "failed passOrFail", tr.Err().Error())
}

func TestTrierErrJoined(t *testing.T) {
	// Arrange
	tr := NewTrier()

	// Act
	tr.Try(passOrFail).
		Try(failIfString, 0).
		Try(passOrFail, true).
		TryJoin(failIfString, "hi")

	// Assert
	assert.Equal(t, "failedIfString\nfailed passOrFail", tr.Err().Error())
}
