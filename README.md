# ddf

 De-Duplicate Files

```bash
 ❯ cp README.md readme2.md
 ❯ ddf
 .gitignore
 Makefile
 README.md
> readme2.md
 go.mod
 go.sum
 main.go

Found 1 repeat
```bash

`ddf` detects piped output and will automatically produce "quiet" output:

```bash
❯ cp README.md readme2.md
❯ ddf | cat
readme2.md
```
