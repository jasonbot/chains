# Chained Iterators in Go

I haven't done much in Go lately, and I definitely haven't played with generics or iterators yet. Why not do both?

So now you can do this:

```go
myArray := []int{1, 2, 3}

evensOnly := func(val int) bool { return int % 2 == 0}
square := func(val int) int{ return val * val}
for _, x := Chain[int](myArray).Filter(evensOnly).Map(square).Each() {
    // ...
}
```
