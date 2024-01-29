# MarkOut

Simplistic cli built with cobra for converting markdown files to html.

This is not a static site generator. This is not pandoc.

Accepts file list or stdin, writes to files or stdout.

```sh
markout example.md -o > output.html
```

## Usage

```text
> markout --help
Convert Markdown files to HTML

Usage:
  markout [flags]

Flags:
  -e, --extension string   Output file extension (default ".html")
  -h, --help               help for markout
  -d, --outdir string      Output directory (default "./outputs")
  -w, --overwrite          Overwrite existing output files
  -s, --stdin              Read input from stdin
  -o, --stdout             Print output to stdout
```

If stdin is used, the output file will be named `MarkOut.html`.

Further examples:

```sh
> markout example.md
Successfully converted example.md to outputs/example.html
```

```sh
> markout example.md -d mypath
Successfully converted example.md to mypath/example.html
```

```sh
markout hw.md -o
```

```sh
markout -w hw.md hello.md example.md ./tests/test.md
```

```sh
cat hello.md | markout -s -o
```

```sh
cat example.md | markout -s -o > mypath/example.html
```

```
./markout -s -w <<EOF
# Hello World!

  - How are you?
  - I'm fine.
EOF
```

## Development


Requires [Task](https://taskfile.dev/), easily installable with `go install github.com/go-task/task/v3/cmd/task@latest`

```sh
task run     # Run the application to print example input to console
task example # Example to generate output file
task all     # Run some build stuff
task build   # Build the application
task test    # Run tests
task lint    # Run linting
```
