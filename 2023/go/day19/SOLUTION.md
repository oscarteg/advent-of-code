# Day 19: Aplenty

Each workflow is a list of rules: a category (`x`, `m`, `a`, `s`), a `<` or
`>` comparison, a threshold and a target. Only the last rule of a workflow is
unconditional, which parsing enforces.

Part 1 walks every part from `in`, always taking the first rule whose condition
holds, until it reaches `A` or `R`, and sums the ratings of the accepted parts:
`O(parts * rules)`.

Part 2 pushes a whole search space (`1..4000` per category) through the same
rules instead of a single part. Each condition splits the incoming space into
the ratings it routes away and the ratings the next rule sees, so accepted
fragments are disjoint and the answer is the sum of the products of their four
span widths. No rating is ever enumerated.

From `2023/go`:

```sh
go test -race ./day19
go vet ./day19
go run ./day19 < day19/testdata/example.txt
```

Expected output:

```text
Part 1: 19114
Part 2: 167409079868000
```
