# MarkOut

Simplistic cli built with cobra for converting markdown files to html.

This is not a static site generator. This is not pandoc. This is Markdown in, Markup (in HTML) out.

Accepts file list or stdin, writes files to directory or logs result to stdout.

```sh
markout input.md -o > output.html
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
  -d, --outdir string      Output directory (default "./markoutput")
  -w, --overwrite          Overwrite existing output files
  -x, --raw                Raw transform of content to HTML (no page head and body)
  -r, --recurse            Run recursively on subdirectory contents
  -s, --stdin              Read input from stdin
  -o, --stdout             Print output to stdout
```

If stdin is used, the output file will be named `MarkOut.html`.

If no input file is supplied, all file in the working directory with an `.md` or `.markdown` extension will be converted.

```sh
cp ./markout ./examples/markout && cd ./examples && ./markout
```

Further examples (run from examples directory):

File input

```sh
$ markout example.md
Successfully converted example.md to markout/example.html
```

Directory output

```sh
$ markout example.md -d mypath
Successfully converted example.md to mypath/example.html
```

Process nested items from working directory (currently ignores subdir name in output)

```sh
markout -r -w
```

Log parsed file content to stdout

```sh
markout hw.md -o
```

Process multiple files (and overwrite existing outputs)

```sh
markout -w hw.txt hello.md example.md ../tests/test.md
```

Cat file contents into stdout

```sh
cat hello.md | markout -s -o
```

Cat file contents into stdout and pipe to file

```sh
cat example.md | markout -s -o > mypath/example.html
```

Skip the page wrapping

```sh
$ echo "# Hello World" | markout -x -s -o
<h1>Hello World</h1>

```

Read in content from heredoc

```
markout -s -w <<EOF
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
