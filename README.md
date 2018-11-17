# Versioned Markdown (VMD)

Versioned Markdown is a markdown specification aimed at defining a markdown
standard tailored for group collaboration on documentation in the context of
version control systems (i.e. git.) A key component of this project will be
a formatting tool, similar to `go fmt`, which can read almost any markdown
robustly and output a stable and consistent format.

## Elements of markdown

Before we give the the example format we must define the elements of the 
markdown format. They are as follows:

- **Paragraphs:** which can be any series of string tokens, including basic 
  strings of text (words and punctuation) as well as intermixed `code` text,
  *emphasis* text, **strong** text, and [links](https://example.com/),
  sperated by very forgiving whitespaces, which cannot extend over
  a blank newline. These lines are treated as an array of tokens, with all
  whitspace rendered as a single space.
- **Headings:** which may only contain a simple title string, and a heading
  level.
- **Quote Blocks:** which are treated the same as paragpraphs, but enclosed in 
  a special formatting, to indent them inward in the final document.
- **Code Blocks:** which simply present their content verbatin.
- **Lists:** which format sub-paragraphs in an ordered fashion, and may contain
  sublists.
- **Tables:** which present paragraphs in a fomatted grid by the renderer.

Each element of markdown **must** have one blank line between it and another
element in the output format.

> It is important to note that each of these elements, except for paragraphs,
> may only appear on the top level of the document, and may not be embedded 
> inside eachother.

## Markdown style guide

This is a basic specification of the Versioned Markdown ouput format.

### Headings are simple

The parser emits heading nodes with heading level and text content. Headings
will only be rendered in `ATX Heading` format, So they will be rendered as:

```
# <heading level 1 text>

## <heading level 2 text>
```

### Paragraphs

Paragraphs are a list of elementary text nodes. Each word, link, empasis, etc.
may be emitted as a string and tokenized. These tokens are then emmited, and
wrapped at the column limit, forming stable paragraphs of text. Paragraphs 
are valid on their own at the top level of the document, as well as in the 
bodies of list items, and table cells.

### Block quotes

Block quotes are parsed as a `BlockQuote` node and will be handled in the 
same way as paragraphs, with the exception that a `>` followed by a space will
be written at the beginning of every line.

> Note: lists are not supported in block quotes.

### Code Blocks

Code blocks are very simple. The only output format accepted is the `backtick`
format, shown as follows:


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

To include the backtick escape sequence inside a codeblock, the codeblock
must have more backticks than any backtick sequence contained inside.


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

Indented or other types of code blocks may be accepted by the parser, but
only the backtick fenced code blocks will be emitted.

### Tables

Tables have the following stable output format:

```
| Table | Heads | should just be | formatted strictly | like | this |
| ----- | ----- | -------------- | ------------------ | ---- | ---- |
| the content itself | should | just | be formatted | like this | x |
| so that a change | to | any | cell | x | x |
| will | be | stable | on the | rest | of the table | 
```

| Table | Heads | should just be | formatted strictly | like | this |
| ----- | ----- | -------------- | ------------------ | ---- | ---- |
| the content itself | should | just | be formatted | like this | x |
| so that a change | to | any | cell | x | x |
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

> Note: links are another example of an exception to the line wrapping rule.
> Links may require a text token longer than an arbitrary limit, but will 
> be treated as one token in a paragraph (i.e. they will start on a new
> line if too long, or blended in with other tokens if they fit.)

> Note: footnote link notation is not supported

### Lists

Markdown lists have two forms: ordered and unordered. They share very common
syntax and may be nested together.

Unordered list items start with a `-` and a space, and continue with two
columns of indentation if a line wraps or a sublist is to be started.

```
- potato - potato
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
   2. dos
```

1. one
1. two
1. three
   1. uno
   1. dos

> If pretty list formatting is enabled, ordered lists will be numbered by 
> their index, while the default settings will be more stable on edits, and 
> let the renderer handle indexing of items.

The one last feature of lists is that they can be added to eachother, but
an ordered list may only contain ordered items, while child lists may be 
of another, homogenous type.

```
1. one
   - uno 
     - spanish
   - un
     - french
2. two
   - dos
   - deux
3. three
   - tres
   - trois
```

1. one
   - uno 
     - spanish
   - un
     - french
2. two
   - dos
   - deux
3. three
   - tres
   - trois

Note that each level of the list applies its indentation on top of the 
previous item's indentation.
