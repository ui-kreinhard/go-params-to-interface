What's this?
===

It's a small tooling based on golang ast parsing and jennifer.

Sometimes it's unavoidably to have methods or functions with a lot of parameters. A good idea is to use a builder like pattern.

E.g.
```

type Test struct {
	
}

func (t *Test) add(a int, b int, c int, d int) {
	
}
```

```
	WithTest(&Test{}).
		WithA(1).
		WithB(1).
		WithC(2).
		WithD(3).
	Add()
```

This small application will autogenerate all interfaces for this kind :)

How to use
===
E.g. 

```go run main.go -file samples/test.go -function add -struct *Test```
