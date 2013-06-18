Preview Github README.md
========================

The preview that mimics repo's homepage in [GitHub](https://github.com).

## Requirements ##

* Go 1.1
* Internet so `&http.Client{}` can make a request to https://api.github.com/markdown/raw

## How to use

```
$ go build -o previewer main.go
$ alias previewer /path/to/previewer # Better to put this into your shell's profile
$ previewer README.md                # Will output generated HTML into stdout
$ previewer /my/repo/readme.md       # Works on any path
```

Pipes it to browser:

```
$ ./previewer README.md | browser
```

or:

```
$ ./previewer --serve 8080 README.md # Previewer is available on http://localhost:8080
````
