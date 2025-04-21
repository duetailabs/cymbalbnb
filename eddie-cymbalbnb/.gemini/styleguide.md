# Repository style guide

Generally, please follow the Go Google style guide. Below are a few rules specific to the Golang
code contained in this repository. For other languages besides Golang, please follow relevant industry leading 
style guides and best practices.

1. For sequences of `if-else` if statements that evaluate the same expression against different values, use
   a `switch` statement if you have more than two conditions. This improves code clarity and maintainability.

**Example**

The following `if-else` if statement has three conditions:

```go
if myStr == "foo" {
  return "bar", nil
} else if myStr == "baz" {
  return "qux", nil
} else {
  return "", fmt.Errorf("myStr value not supported: %s", myStr)
}
```

This should instead be written as a `switch` statement:

```go
switch myStr {
  case "foo":
    return "bar", nil
  case "baz":
    return "qux", nil
  default:
    return "", fmt.Errorf("myStr value not supported: %s", myStr)
}
```

And similar for all `if-else` statements above two condition checks.

2. Ensure consistent visibility for fields within a struct; avoid mixing exported (uppercase) and
   unexported (lowercase) fields.

**Example**

For the following structs:

```go
type Foo struct {
  bar string
  baz string
}

type Qux struct {
  Quux  string
  Corge string
}
```

Only a private (i.e. lowercased) field can added to `Foo`, and a public (i.e. uppercased) field can be added to `Qux`.

3. Ensure that proper and consistent indentation is used for all code added since we do not have linters in place
   for Golang at the moment. While the incorrectly indented code will still compile, it degrades maintainability to have
   inconsistently indented code.

4. Use `slog` for all logging statements. Explicitly, do not use `log`.
