# Versioned Markdown (VMD)

Versioned Markdown is a 
[markdown](https://en.wikipedia.org/wiki/Markdown) specification tailored to
group collaboration on documents. A key component of this project is 
a formatting tool, similar to `go fmt`, which can parse markdown into a
tree, and re-render it in a standard and stable output format.

> The reference implementation is a go library that utilizes the
> [blackfriday.v2](https://github.com/russross/blackfriday/tree/v2) markdown
> library to parse a robust segment of markdown formats into a standard
> output. See `vmdfmt`.

## Versioned Markdown Specification

This is a basic specification of the Versioned Markdown ouput format. On
the root level of the document, each of these elements must have exactly
one empty line between it and the next element. An output file will end
with an empty line.

### Paragraphs

Paragraphs are formatted as a collection of words (and punctuation) with 
inline formatting, emitted as tokens, which are line wrapped at a predefined
column number (default: 80). There is a single space between each token.

> Links are treated as a single token, because of their formatting.

#### Inline Formatting

```
Inline formatting blocks include *italic*, **bold**, and `code` formatting.
```

Inline formatting blocks include *italic*, **bold**, and `code` formatting.

### Headings

Headings will only be rendered in `ATX Heading` format, So they will be
rendered as:

```
# <heading level 1 text>

## <heading level 2 text>

etc.
```

> Note: Heading bodies may only contain plain text.

### Block Quotes

Block quotes are treated almost identically to paragraphs, except that each
line will begin with a '>' and a single space before the text begins. If a
block quote is embedded in another block quote, an additional '>' and
another single space will be added in addition to the first.

Block quotes may only contain content which is valid in paragraphs, in 
addition to nested block quotes.

> Note: Block Quotes next to eachother with empty lines are parsed as a single
> block quote. This is not a style issue, but inherited from the *blackfriday*
> parser.

### Code Blocks

The only output format accepted is the `backtick` format, shown as follows:

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

To quote codeblock notation inside of a codeblock, you may increase the 
number of backticks. Any backticks of a lesser magnitude inside of the
code block will be ignored. The reference implementation simply reuses
the fence length from the original document.

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

> A formatting implementation may accept multiple formats of block quotes,
> but any tilde fences or indented fences containing backtick fences may
> have undefined behaviour.

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

> Note: tables are an exception to the line wrapping rule. Their formatting
> is line based so they cannot be wrapped

### Links

Links are simple and have two representations.

Links whose text are the same as their content can be represented in this way:

```
<http://www.example.com/>
```

<http://www.example.com/>

And links with description text which is not the url itself, will be written
as:

```
[descriptive text](http://www.example.com/)
```

[descriptive text](http://www.example.com/)

No other format of links is supported.

> Note: links are another exemption to the line wrapping rule.

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

Sublists may be of a different type than the parent list, but list types
may not be mixed.

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

> Note: list inputs must have at least 3 columns of indentation. this is not
> a style issue, but rather a limitation of the *blackfriday* markdown parser.
