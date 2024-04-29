package viewutils

import (
	"context"
	"io"

	"github.com/a-h/templ"
)

// render all items in a list with a component
func RenderListWithComponent[T interface{}](list []T, component func(T) templ.Component) templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		for _, item := range list {
			err := component(item).Render(ctx, w)
			if err != nil {
				return err
			}
		}
		return nil
	})
}
