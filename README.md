Preview Github README.md
========================

The preview that mimics repo's homepage in [GitHub](https://github.com).

## Requirements ##

* Go 1.1
* Internet so `&http.Client{}` can make a request to https://api.github.com/markdown/raw

## How to use

```
$ go install github.com/gedex/preview-github-readme
$ preview-github-readme README.md           # Will output generated HTML into stdout
$ preview-github-readme /my/repo/readme.md  # Works on any path
```

Pipes it to browser:

```
$ preview-github-readme README.md | browser
```

or:

```
$ preview-github-readme --serve 8080 README.md # Previewer is available on http://localhost:8080
````
