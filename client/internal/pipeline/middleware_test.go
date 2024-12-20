package pipeline

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDecorateHandler(t *testing.T) {
	var order []string

	base := HandlerFunc(func(_ *MiddlewareContext) error {
		order = append(order, "base")
		return nil
	})

	newMiddleware := func(name string) MiddlewareFunc {
		return func(next Handler) Handler {
			return HandlerFunc(func(ctx *MiddlewareContext) error {
				order = append(order, name+"_before")
				if err := next.Handle(ctx); err != nil {
					return err
				}
				order = append(order, name+"_after")
				return nil
			})
		}
	}

	h1 := newMiddleware("h1")
	h2 := newMiddleware("h2")
	h3 := newMiddleware("h3")

	t.Run("handler", func(t *testing.T) {
		order = []string{}
		expected := []string{"base"}
		err := base.Handle(nil)
		require.NoError(t, err)
		require.Equal(t, expected, order)
	})

	t.Run("decorated", func(t *testing.T) {
		t.Run("alone", func(t *testing.T) {
			order = []string{}
			expected := []string{"base"}
			err := NewPipeline(base).Handle(nil)
			require.NoError(t, err)
			require.Equal(t, expected, order)
		})

		t.Run("multiple", func(t *testing.T) {
			order = []string{}
			expected := []string{
				"h1_before",
				"h2_before",
				"h3_before",
				"base",
				"h3_after",
				"h2_after",
				"h1_after",
			}
			err := NewPipeline(base, h1, h2, h3).Handle(nil)
			require.NoError(t, err)
			require.Equal(t, expected, order)
		})
	})
}
