# Versioned Markdown (VMD)

Versioned Markdown is a markdown specification aimed at defining a markdown
standard tailored for group collaboration on documentation in the context of
version control systems (i.e. git.) A key component of this project will be
a formatting tool, similar to `go fmt`, which can read almost any markdown
robustly and output a stable and consistent format.

## Markdown style guide

This is a basic mock-up of the markdown ouput style.

### Headings are simple

The parser emits heading nodes with heading level and text content. Headings
will only be rendered in `ATX Heading` format, So they will be rendered as:

```
# <heading level 1 text>

## <heading level 2 text>

etc.
```

### Paragraphs

Paragraph nodes are a list of `text`, `quote`, `link`, `emphasis`, etc.,
nodes. These are the basic container for string type tokens in the parser.
Each type of node contained in a paragraph may be rendered and tokenized,
to be written out with one space between tokens and wrapped at an arbitrary
number of characters (in most cases, 80.)

### Code Blocks

Code blocks are also simple. They are all represented the same, and the node
contains metadata noting whether or not they are fenced.

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

Indented or other types of code blocks may be accepted by the parser, but
only the backtick fenced code blocks will be emitted.

### Tables

Tables will be output in the following manner. Note that, while pretty
formatting of the tables may be visually appealing, this format is both
stable in the context of version control diffs, and also easier to work with
when editing.

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

### Block quotes

Block quotes are parsed as a `BlockQuote` node and will be handled in the 
same way as paragraphs, with the exception that a `>` followed by a space will
be written at the beginning of every line.

> Note: lists are not supported in block quotes.

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
