# Go Lexandra
Lexer for Golang

# How it work?
A Lexer its just an process that tokenize (create tokens) a string of chars. For do that it has to be created each individual Tokenizer, i.e: the Number Tokenizer, the Word Tokenizer, etc. 

**Those Tokenizer are called: `Lex`**.

So, a Lexer is just a set of `Lex`s. 

```
         Lexer
       /   |   \  
      /    |    \
     /  LexWord  \
    |             |
LexNumb        LexString
```

When a string of char (as *bytes*) is reading, it will asking to each **Lex**: *do you allow this char?*. In case of *yes*, the **Lex** will run; otherwise it will continue with the next one.

# Quick start

Let's create a simple number lexer.

First, create `main.go` and instance a new `Lexer`.

```go
package main

import golex "github.com/agustin-del-pino/go-lexandra"

func main(){
    lexer := golex.NewLexer[int]()
}
```

Then, create the `Lex` instance.

```go
var lexNumb := &golex.Lex {
    AllowedBytes: golex.BytesRange(0x30, 0x39),
}
```

The `AllowedBytes` attribute indicates the bytes must be match for run this Lex. In this case is using a `ByteRange` meaning the bytes greater or equal than `0x30` and less or equal thant `0x39` are allowed.

*From `0x30` to `0x39` are the ASCII Numbers (0 to 9)*.

Let's add the `Run` function of the Lex.

```go
var lexNumb := &golex.Lex[int] {
    AllowedBytes: golex.BytesRange(0x30, 0x39),

    // The 'Run' func takes the Cursor and the Lex itself. 
    // It must return the created token.
    Run: func(c *golex.Cursor, l *golex.Lex[int]) golex.Token[int] {
        // Creates the token. As Token Type is using the: 0. 
        // And adds as initial value the current read char.
        t := golex.NewToken(0, c.Char)

        // Advance to the next char.
        c.Advance()

        // Read the char until reach the end of the bytes 
        // or the current char is not allowed.
        for c.HasChar() && l.AllowedBytes(c.Char) {
            
            // Adds the current char to the Value of the token.
            t.Value = append(t.Value, c.Char)
            
            // Advance to the next char.
            c.Advance()
        }

        // Returns the token.
        return t
    }
}
```
*Read the comments*

What's happening here?. Well, the `Cursor` is the representation of the *"eye"* that reads the string of char. With this *"eye"*, it can be read the current char, asking whether there's no more char to read and also, and very import, advance to the next char. In other words, it's the action that tells the `Cursor` it has to read the next char.

So far so good. Let's finish this!.

```go
package main

import golex "github.com/agustin-del-pino/go-lexandra"

var lexNumb := &golex.Lex[int] {
    AllowedBytes: golex.BytesRange(0x30, 0x39),

    // The 'Run' func takes the Cursor and the Lex itself. 
    // It must return the created token.
    Run: func(c *golex.Cursor, l *golex.Lex[int]) golex.Token[int] {
        // Creates the token. As Token Type is using the: 0. 
        // And adds as initial value the current read char.
        t := golex.NewToken(0, c.Char)

        // Advance to the next char.
        c.Advance()

        // Read the char until reach the end of the bytes 
        // or the current char is not allowed.
        for c.HasChar() && l.AllowedBytes(c.Char) {
            
            // Adds the current char to the Value of the token.
            t.Value = append(t.Value, c.Char)
            
            // Advance to the next char.
            c.Advance()
        }

        // Returns the token.
        return t
    }
}

func main(){
    lexer := golex.NewLexer[int]()

    // Register the Lex.
    lexer.Numbs(lexNumbs)

    // Starts the tokenize process 
    // and then returns the Token Slice.
    tokens := lexer.Tokenize([]bytes("0123456789"))
}
```

Good this is done. But before to end, it necessary to review some few things.

### The Generic Token Type
As you can see before. It was declared as Generic to the Lexer and it propagates to the Lex. 

This Generic indicates the type of the Token Type. In the example it was used the `int` type, but it can be used any type. It's recommend use an Alias Type like `type TokenType int` and defined an `Enum` with the different Token Types.

### The Registration of a Lex
By default the lexer as a pre-defined Lex registers, those are and in order of execution:
- Numbs
- Strings
- Words
- Delimiters
  
### Extending the Lexer
In case you want to register a non pre-defined Lex, you have to extend the Lexer. For do that just use the following template:

```go

// The user-defined Lex.
var lexComment = &golex.Lex[TokenType]{
	AllowedBytes: golex.ByteSingle(0x23),
	Run: func(c *golex.Cursor, l *golex.Lex[TokenType]) golex.Token[TokenType] {
		t := golex.NewToken(TComment)
		c.Advance()
		for c.HasChar() && c.Char != 0x0A {
			t.Value = append(t.Value, c.Char)
			c.Advance()
		}
		return t
	},
}

// The extension function.
func lexExtend(l *golex.Lexer[TokenType], c *golex.Cursor, t *golex.Tokens[TokenType]) func() bool {
	m := l.Get("comments", true)

	return func() bool {
		if m.AllowedBytes(c.Char) {
			*t = append(*t, m.Run(c, m))
			return true
		}

		return false
	}
}


func main(){
    lexer := golex.NewLexer[int]()

    // Register the user-defined Lex.
    lexer.Register("comments", lexComments)
    
    // Adds the extension fo the Lexer.
    lexer.Extension(lexExtend)

    tokens := lexer.Tokenize([]bytes("# I'm a comment."))
}
```

### Ignoring chars
In case you want to ignore some specific chars, just declared them as params of the Lexer Constructor.

```go
// It's ignoring the New Lines and the White Spaces.
var lexer := golex.NewLexer[int](0x0A, 0x20)
```