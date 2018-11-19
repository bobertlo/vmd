# Versioned Markdown (VMD)

Versioned Markdown is a 
[markdown](https://en.wikipedia.org/wiki/Markdown) specification tailored to
group collaboration on documents. A key component of this project is 
a formatting tool, similar to `go fmt`, which can parse markdown into a
tree, and re-render it in a standard and stable output format.

> The reference implementation is a go library that utilizes the
> [blackfriday.v2](https://github.com/russross/blackfriday/tree/v2) markdown
> library to parse a robust segment of markdown formats into a standard
> output.

## Versioned Markdown Specification

This is a basic specification of the Versioned Markdown ouput format. On
the root level of the document, each of these elements must have exactly
one empty line between it and the next element.

### Paragraphs

Paragraphs are formatted as a collection of words (and punctuation) with 
inline formatting, emitted as tokens, which are line wrapped at a predefined
column number (default: 80). There is a single space between each token.

> Links are treated as a single token, because of their formatting.

### Headings

Headings will only be rendered in `ATX Heading` format, So they will be
rendered as:

```
# <heading level 1 text>

## <heading level 2 text>

etc.
```

> Heading bodies may only contain plain text. They may not contain
links, emphasis, or any other formatting.

### Block quotes

Block quotes are treated almost identically to paragraphs, except that each
line will begin with a '>' and a single space before the text begins.

> Block quotes may only contain content which is valid in paragraphs, in 
> addition to nested block quotes.

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

Tables have the following stable output format:

```
| Table | Heads | should just be | formatted strictly | like | this |
| ----- | ----- | -------------- | ------------------ | ---- | ---- |
| the content itself | should | just | be formatted | like this | x |
| so that a change | to | any | cell | on | a single line |
| will | be | stable | on the | rest | of the table | 
```

| Table | Heads | should just be | formatted strictly | like | this |
| ----- | ----- | -------------- | ------------------ | ---- | ---- |
| the content itself | should | just | be formatted | like this | x |
| so that a change | to | any | cell | on | a single line |
| will | be | stable | on the | rest | of the table | 

> Note: tables are an exception to the line wrapping rule. Their formatting
> may require them to be longer than an arbitrary limit.

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

> Note: links are another exemption to the line wrapping rule. Their
> formatting requires them to exist on a single line.

> Note: footnote link notation is not supported.

### Lists

Unordered list items start with a `-` and a space, and continue with two
columns of indentation if a line wraps or a sublist is to be started.

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

Ordered lists will start with `1.` followed by a space, and continue with 
three columns of indentation if a line wraps, or a sublist is to be started.

```
1. one
1. two
1. three
   1. uno
   1. dos
```

1. one
1. two
1. three
   1. uno
   1. dos

The one last feature of lists is that they can be added to eachother, but
an ordered list may only contain ordered items, while sub-lists may be 
of another, also homogenous, class.

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

Note that each level of the list applies its precise indentation in addition
to the indentation from the parent list.
