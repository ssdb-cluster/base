# base

```
# excludes `Test***`
go test -bench Bench -run B -benchmem
                          ~

go test -bench Bench -run Bench -benchmem
               ~~~~~      ~~~~~
```
