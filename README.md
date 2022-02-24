# golang-enum-to-ts
Transform Golang `enum` type to Typescript enum

## Function

### Before (Golang)
```golang
package some

type Status int

const (
	Todo Status = iota
	Done
	Pending
	InProgress
)

type Sex string

const (
	Female Sex = "female"
	Male   Sex = "male"
)

func Abctext() {
	//dadsad
}
```

### After (Typescript)

```ts
namespace some {
  export enum Sex {
    Female = 'female',
    Male = 'male',
  }
  export enum Status {
    Pending = 2,
    InProgress = 3,
    Todo = 0,
    Done = 1,
  }
}
```