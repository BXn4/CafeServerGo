package utils

// This is a ternary operator
//	utils.If(true, "foo", "bar") // returns "foo"
//	utils.If(false, "foo", "bar") // returns "bar"
func If[T any](condition bool, a T, b T) T {
  if condition {
    return a
  }
  return b
}
