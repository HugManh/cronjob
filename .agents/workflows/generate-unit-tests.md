---
description: 
---

## Unit Test Generation Workflow (Golang)

**Trigger:** When the user invokes this workflow.

### Instructions:

**Analyze all Go source files in the current active context**

* Identify all `.go` files excluding:

  * `_test.go` files
  * `main.go`
  * generated files (`*_gen.go`, `mock_*.go`, etc.)

---

**Create corresponding test files**

* For every source file (e.g., `utils.go`), create a corresponding test file named:

  ```bash
  utils_test.go
  ```

* The test file must be in the same package or use the `_test` package pattern:

  ```go
  package utils
  ```

  or

  ```go
  package utils_test
  ```

---

**Use Go's standard testing framework**

* Use the built-in `testing` package:

  ```go
  import "testing"
  ```

* Follow the standard naming convention:

  ```go
  func TestFunctionName(t *testing.T)
  ```

---

**Test coverage requirements**

For every exported and critical function, implement at least:

1. **Positive test case**

   * Validate expected behavior with valid input
   * Verify correct return values and no unexpected errors

2. **Edge case test**

   * Boundary values
   * Empty input
   * Nil input (if applicable)
   * Invalid input
   * Error scenarios

---

**Example**

**utils.go**

```go
package utils

func Add(a, b int) int {
    return a + b
}
```

**utils_test.go**

```go
package utils

import "testing"

func TestAdd_Positive(t *testing.T) {
    result := Add(2, 3)

    if result != 5 {
        t.Errorf("expected 5, got %d", result)
    }
}

func TestAdd_EdgeCase_Zero(t *testing.T) {
    result := Add(0, 0)

    if result != 0 {
        t.Errorf("expected 0, got %d", result)
    }
}
```

---

**Recommended best practices**

* Use table-driven tests:

```go
func TestAdd_TableDriven(t *testing.T) {
    tests := []struct {
        name string
        a    int
        b    int
        want int
    }{
        {"positive", 2, 3, 5},
        {"zero", 0, 0, 0},
        {"negative", -1, -2, -3},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got := Add(tt.a, tt.b)
            if got != tt.want {
                t.Errorf("expected %d, got %d", tt.want, got)
            }
        })
    }
}
```