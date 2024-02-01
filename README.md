# MarkOut

Simple cli for md => html

> This is not a static site generator. This is not pandoc. It's Markdown in, HTML out.

Accepts file list or stdin, writes files to directory or logs result to stdout.

```sh
markout input.md -o > output.html
```

## Usage

```text
Convert Markdown files to HTML

Usage:
  markout [flags]

Flags:
  -h, --help               help for markout
  -i, --stdin              Read input from stdin
  -o, --stdout             Print output to stdout
  -w, --overwrite          Overwrite existing output files
  -r, --recurse            Run recursively on subdirectory contents
  -e, --extension          Output file extension (default ".html")
  -d, --outdir             Output directory (default "./markoutput")
  -f, --full               Write complete HTML page (including head, with md content in body)
  -s, --style              Path to a css file. Contents are injected into a <style> block (default "none")
  -l, --link               Text value to insert into the href attribute of <link rel="stylesheet" />. (default "none")
  -t, --theme              A predefined css to embed. Options include "none", "pandoc", and "retro". (default "none")
```

If stdin is used, the output file will be named `MarkOut.html`.

If stdout is used, the response will be written as raw html (no use of `--full`), otherwise `-f` is assumed by default for file outputs.

If no input file is supplied, all files in the current working directory with an `.md` or `.markdown` extension will be converted.

```sh
cp ./markout ./examples/markout && cd ./examples && ./markout
```

Further examples (run from examples directory):

Process all markdown files in working directory, generating full page output with default embedded theme.

```sh
markout -f
```

Overwrite outputs and remove stylesheets

```sh
markout -f -w -t none
```

Use link (cdn theme)

```sh
markout -f -l bootstrap
```

Use local stylesheet file

```sh
markout -f -w -s css.css
```

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
cat hello.md | markout -i -o
```

Cat file contents into stdout and pipe to file

```sh
cat example.md | markout -i -o > mypath/example.html
```

Skip the page wrapping

```sh
$ echo "# Hello World" | markout -i -o
```

```html
<h1>Hello World</h1>

```

Add the page wrapping

```sh
$ echo "# Hello World" | markout -f -i -o
```

```html
<!DOCTYPE html>
<html>
  <head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>MarkOut</title>
  </head>
  <body>
    <h1>Hello World</h1>
  </body>
</html>
```

Read in content from heredoc

```
markout -i -w <<EOF
# Hello World!

  - How are you?
  - I'm fine.
EOF
```

## Page Style

CSS Themes can be injected into the *page* to make your markdown pretty, assuming that `--full` mode has been used (i.e. without the `-f` flag, the resulting html will be style-less).

There are three types of CSS support for style management included in markout, namely:

  - **themes**: These baked in themes are minified and injected directly into the page in a `<style>` block. The are assigned by name (e.g. `retro`). If no other style options are specified, the *pandoc* theme will be applied by default. Styles can be removed entirely by using `-t none` to disable the default theme.
  - **styles**: The path to a stylesheet file that will have its contents injected directly into the page within a `<style>` block.
  - **links**: Links come from the interwebs and load via CDNs (so therefore require an active internet connection and  to function on the output pages). This has a ton of benefits, but also caveats (e.g. CORS issues) when loading your pages locally. You can specify a link to include as an href attribute, or select one of the predefined css frameworks by name (e.g. `-l bulma`).

Although each parameter is only designed to support one configuration at a time, these options are not mutually exclusive and can be combined to inject multiple stylesheets into a single page.

If you need more control over your markdown styling, it is likely that you are using the wrong tool for your needs.

### Themes

The following themes are included OOTB:

  - [pandoc](https://gist.github.com/killercup/5917178)
  - [retro](https://markdowncss.github.io/retro/)

Sourced from the community, they may have their own license considerations. Yolo.

### Named Links

The following links (online only themes) are included by name:

  - [milligram](https://milligram.io/#getting-started)
  - [wing](https://kbrsh.github.io/wing/)
  - [pico](https://picocss.com/docs/)
  - [bootstrap](https://www.bootstrapcdn.com/)
  - [bulma](https://bulma.io/documentation/overview/start/)
  - [tachyons](https://tachyons.io/)

Note that all named links include the attributes `crossorigin="anonymous" referrerpolicy="no-referrer"` and none of them will include integrity check attributes.

Specifying any other string value will inject that value as the href content directly, and modifying link attributes is unsupported.

### Scripts

Explicit support for script etc aren't expected to be supported, because it's not a tool meant for that. With that said, all content is pretty naively injected from file contents, so you can just break out of a user supplied local style file.

E.g. `</style><script>alert('Naughty!')</script><style>` and then `markout README.md -o -f -w -l milligram -s injection.js > mischief.html`.

## Development

Requires [Task](https://taskfile.dev/), easily installable with `go install github.com/go-task/task/v3/cmd/task@latest`

Then, just

```sh
task
```

Or for individual operations:

```sh
task run     # Run the application with args
task example # Example to generate stdout
task tidy    # tidy build+generate modules 
task lint    # Code quality
task test    # Run tests
task build   # Build the application
```
