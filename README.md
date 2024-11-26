# Chained Iterators in Go

I haven't done much in Go lately, and I definitely haven't played with generics or iterators yet. Why not do both?

So now you can do this:

```go
myArray := []int{1, 2, 3}

returnArray := Chain(
    []int{8, 10, 145, 3},
).Map(
    func(i int) int {
        return i * 3
    },
).Filter(
    func(i int) bool {
        return i%2 == 0
    },
).A()
// returnArray == []int{24, 30}
```

## Warts

The Go templating system is a little limited, so you can't do something like this:

```go
secondArray := Chain(
    []int{8, 10, 145, 3},
).Map(
    func(i int) str {
        return fmt.Sprintf("%v", i)
    },
)
```

The generic system does not allow for templated methods, so chaining methoids and expecting to go from `Chain[T]` to `Chain[V]` isn't possible.

You need to give the templating system a hint with a junction, telling it there's 2 types involved:

```go
mapFunc := func(i int) string { return fmt.Sprintf("%v", i) }
array := []int{1, 2, 3, 4}
returnArray := Junction2[int, string](
    Chain(
        array,
    ).Filter(
        func(i int) bool {
            return i%2 == 0
        },
    )).Map(mapFunc).A()
// returnArray == []string{"2", "4"}
```
