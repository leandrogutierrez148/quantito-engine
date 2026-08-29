package infraestructures

import (
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockWriter struct {
	mock.Mock
}

func (w *mockWriter) Write(p []byte) (int, error) {
	args := w.Called(p)
	return args.Int(0), args.Error(1)
}

func TestSplitWritter(t *testing.T) {
	mockInfoWriter := &mockWriter{}
	mockInfoWriter.On("Write", mock.Anything).Return(4, nil)
	mockErrorWriter := &mockWriter{}
	mockErrorWriter.On("Write", mock.Anything).Return(5, nil)

	writer := NewSplitWriter(mockInfoWriter, mockErrorWriter)

	n, err := writer.WriteLevel(zerolog.InfoLevel, []byte("INFO"))
	assert.NoError(t, err)
	assert.Equal(t, 4, n)
	mockInfoWriter.AssertCalled(t, "Write", []byte("INFO"))

	n, err = writer.WriteLevel(zerolog.ErrorLevel, []byte("ERROR"))
	assert.NoError(t, err)
	assert.Equal(t, 5, n)
	mockErrorWriter.AssertCalled(t, "Write", []byte("ERROR"))
}

func TestSplitWritterNoLevel(t *testing.T) {
	mockInfoWriter := &mockWriter{}
	mockErrorWriter := &mockWriter{}

	mockInfoWriter.On("Write", []byte("Sin nivel")).Return(4, nil).Once()

	writer := NewSplitWriter(mockInfoWriter, mockErrorWriter)

	n, err := writer.WriteLevel(zerolog.NoLevel, []byte("Sin nivel"))
	assert.NoError(t, err)
	assert.Equal(t, 4, n)

	mockInfoWriter.AssertExpectations(t)
}
