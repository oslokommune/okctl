// Package action is a work in progress
package action

// Action represents an action that will be performed
type Action func(action interface{}) (response interface{}, err error)

// Middleware makes it possible to perform activities on actions prior to sending
type Middleware func(Action) Action

// Chain is a helper function for composing multiple actions
func Chain(outer Middleware, others ...Middleware) Middleware {
	return func(next Action) Action {
		for i := len(others) - 1; i >= 0; i-- {
			next = others[i](next)
		}

		return outer(next)
	}
}
