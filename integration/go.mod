module github.com/Zweih/go-rpmdb/integration

go 1.24.3

require github.com/stretchr/testify v1.10.0

require (
	github.com/Zweih/go-rpmdb/pkg v0.0.0
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace github.com/Zweih/go-rpmdb/pkg => ../pkg
