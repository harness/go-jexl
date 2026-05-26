This repository is a Go port of [Apache JEXL](https://commons.apache.org/proper/commons-jexl/), the Java Expression Language from Apache Commons. It was created by taking the [expr-lang/expr](https://github.com/expr-lang/expr) library as a starting point and adapting its interpreter to evaluate JEXL syntax and semantics rather than expr's native language.

## Testing

All tests must pass. The integration and conformance tests in `tests/**` are especially important.

**Do not modify the synthetic, integration, or conformance tests in `tests/**`.** These define language conformance — if they fail, fix the implementation, not the test.

## Coverage

Run the coverage report before submitting changes:

```
go test -cover ./...
```

When adding or editing new code, please make a reasonable effort to ensure it is covered by tests. Please prefer meaningful tests over tests written purely to hit coverage numbers.
