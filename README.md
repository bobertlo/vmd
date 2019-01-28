[![Build
Status](https://travis-ci.org/bobertlo/vmd.svg?branch=master)](https://travis-ci.org/bobertlo/vmd)
[![Go Report
Card](https://goreportcard.com/badge/github.com/bobertlo/vmd)](https://goreportcard.com/report/github.com/bobertlo/vmd)

# Versioned Markdown (VMD)

Versioned Markdown is a [markdown](https://en.wikipedia.org/wiki/Markdown)
specification tailored to group collaboration on documents. A key component of
this project is a formatting tool, similar to `go fmt`, which can parse markdown
into a tree, and re-render it in a standard and stable output format.

## vmdfmt Auto Formatter

The included `vmdfmt` formatter tool works very similiarly to the `gofmt` tool
included with go. It will format files given as arguments, with the following
flags:

- `-w`: write changes back to source files, instead of `stdout`.
- `-l`: list files which have been changed. If `-w` is not active, it will only
   output the list of files with changes, and not write the formatted changes
   anywhere.
- `-cols int`: change the number of columns to wrap lines at (default: 80.)

`vmdfmt` uses the
[blackfriday.v2](https://github.com/russross/blackfriday/tree/v2) markdown
library to parse a large set of input markdown formats, but emits the parsed AST
in single output format.

## VMD mdformatter package

Also included is the `mdformatter` package, which exposes a very simple
interface to render markdown. It may be imported from:

```
github.com/bobertlo/vmd/pkg/mdformatter
```

To render a `[]byte` slice of markdown, returning a formatted `[]byte` slice,
for example:

```
md := mdformatter.New(80) // New takes the column number to wrap at
out, err := md.RenderBytes(input)
```

## Versioned Markdown Specification

After parsing a document, the formatter will emit each of the following top
level blocks from the AST. Each block will have a single blank line between
them, with no black newline at the end of the file.

### Paragraphs

Paragraphs are formatted as a collection of words (and punctuation) with inline
formatting, emitted as tokens, which are line wrapped at a predefined column
number (default: 80). There is a single space between each token.

#### Inline Formatting

```
Inline formatting blocks include *italic*, **bold**, and `code` formatting.
```

Inline formatting blocks include *italic*, **bold**, and `code` formatting.

These formatting blocks will be parsed inline as a single output string with
their siblings. Then the formatter wraps the lines at spaces if they are in an
appropriate parent container.

#### Links

Links whose text are the same as their content can be represented in this way:

```
<http://www.example.com/>
```

<http://www.example.com/>

And links with description text which is not the url itself, will be written as:

```
[descriptive text](http://www.example.com/)
```

[descriptive text](http://www.example.com/)

Links are treated as a single token, because of their formatting. If their
descriptive text contains spaces they may be linewrapped, otherwise they are an
exception to the linewrapping rule.

#### Images

Images are supported with the following format:

```
![descriptive text](http://www.example.org/logo.jpg)
```

The discriptive text may not contain any formatting, only plain text.

### Headings

Headings will only be rendered in `ATX Heading` format, So they will be rendered
as:

```
# <heading level 1 text>

## <heading level 2 text>
```

> Note: Heading bodies may not contain inline formatting, only text.

### Block Quotes

Block quotes are treated almost identically to paragraphs, except that each line
will begin with a '>' and a single space before the text begins. If a block
quote is embedded in another block quote, an additional '>' and another single
space will be added in addition to the first.

Block quotes may only contain content which is valid in paragraphs, in addition
to nested block quotes.

> Note: Block Quotes next to eachother with empty lines are parsed as a single
> block quote. This is not a style issue, but inherited from the
> *blackfriday.v2* parser.

### Code Blocks

Code blocks will only be emitted in the `backtick` format, shown as follows:

``````
```
Code block example:

The text inside is passed on verbatim.
```
``````

```
Code block example:

The text inside is passed on verbatim.
```

To quote codeblock notation inside of a codeblock, you may increase the number
of backticks. Any backticks of a lesser magnitude inside of the code block will
be ignored. The reference implementation simply reuses the fence length from the
original document.

`````````
``````
```
code block inside a code block
```
``````
`````````

``````
```
code block inside a code block
```
``````

> A formatting implementation may accept multiple formats of block quotes, but
> any tilde fences or indented fences containing backtick fences may have
> undefined behaviour.

### Tables

Tables have the following format:

```
| Table       | Heads  | should be      | formatted    | like | this   |
|-------------|--------|----------------|--------------|------|--------|
| the content | itself | should         | be formatted | like | this   |
| with        | single | spaces between | words and    | just | one    |
| extra space | around | each of the    | longest cell | per  | column |
```

| Table       | Heads  | should be      | formatted    | like | this   |
|-------------|--------|----------------|--------------|------|--------|
| the content | itself | should         | be formatted | like | this   |
| with        | single | spaces between | words and    | just | one    |
| extra space | around | each of the    | longest cell | per  | column |

> Note: tables are an exception to the line wrapping rule. Their formatting is
> line based so they cannot be wrapped

### Lists

Unordered list items start with a `-` and a space, and continue with three
columns of indentation if a line wraps or a sublist is started.

```
- potato
- potato
   - tomato
   - tomato
```

- potato
- potato
   - tomato
   - tomato

Ordered lists will start with their index (i.e`1.`) followed by a space, and
continue with three columns of indentation if a line wraps, or a sublist is
started.

```
1. one
2. two
3. three
   1. uno
   2. dos
```

1. one
2. two
3. three
   1. uno
   2. dos

Sublists may be of a different type than the parent list, but list types may not
be mixed.

```
1. one
   - uno 
   - un
2. two
   - dos
   - deux
3. three
   - tres
   - trois
```

1. one
   - uno
   - un
2. two
   - dos
   - deux
3. three
   - tres
   - trois

> Note: list inputs must have at least 3 columns of indentation when wrapping
> lines or creating sublists. This is not a style issue, but another limitation
> inherited from the *blackfriday* markdown parser.
