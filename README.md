# Versioned Markdown (VMD)

Versioned Markdown is a 
[markdown](https://en.wikipedia.org/wiki/Markdown) specification tailored
group collaboration on documents. A key component of this project will be
a formatting tool, similar to `go fmt`, which can read markdown input with the
[blackfriday.v2](https://github.com/russross/blackfriday/tree/v2) markdown
library, and re-render the markdown in a stable, consistent format.

## Elements of Versioned Markdown

Before we give the the example format we must define the elements of the 
Versioned Markdown format. They are as follows:

- **Paragraphs:** which can be any series of string tokens, including basic 
  strings of text (words and punctuation) as well as intermixed `code` text,
  *emphasis* text, **bold** text, and [links](https://example.com/),
  sperated by very forgiving whitespaces, which cannot extend over
  a blank newline. These lines are treated as a list of string tokens, with 
  all whitspace rendered as a single space.
- **Headings:** which may only contain plain text title, 
- **Quote Blocks:** which are treated the same as paragpraphs, but enclosed in 
  a special formatting, to indent them inward in the final document.
- **Code Blocks:** which simply present their content verbatin.
- **Lists:** which are an ordered container for paragraphs, and may also
  contain sub-lists.
- **Tables:** which present paragraphs in a fomatted grid by the renderer.

Each of these elements *must* have one blank line between it and another
element in the output format (excluding links and inline formatting inside
a paragraph) and aside from paragraphs, none of these elements may be nested
inside of one another.

## Markdown style guide

This is a basic specification of the Versioned Markdown ouput format.

### Paragraphs

Paragraphs are formatted as a collection of words and inline formatting blocks,
rendered with a single space between each token, and wrapped at an arbitrary
column limit (the default is 80 columns).

### Headings are simple

Headings will only be rendered in `ATX Heading` format, So they will be
rendered as:

```
# <heading level 1 text>
```

or

```
## <heading level 2 text>
```

### Block quotes

Block quotes are treated almost identically to paragraphs, except that each
line will begin with a '>' and a single space before the text begins.

> Note: lists are not supported in block quotes.

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

To include backtick fences inside a codeblock, the opening fence must have
more backticks than any backtick sequence contained inside. The reference
implementation will use the same length of fence as is parsed, for
simplicity.


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

> If pretty list formatting is enabled, ordered lists will be numbered by 
> their index, while the default settings will be more stable on edits, and 
> let the renderer handle indexing of items.

The one last feature of lists is that they can be added to eachother, but
an ordered list may only contain ordered items, while sub-lists may be 
of another, homogenous type.

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

Note that each level of the list applies its indentation on top of the 
previous item's indentation.
