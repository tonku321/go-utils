1. make changes
2. commit
3. call this with new version: `git tag v0.1.5 && git push origin main --tags`
4. run this again in a project: `go get github.com/tonku321/go-utils@latest` or change manually in go.mod then `go mod tidy`