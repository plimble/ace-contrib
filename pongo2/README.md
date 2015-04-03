Pongo2
========
more information about [pongo2](https://github.com/flosch/pongo2)

###Example

######main.go

```
package main

import (
	"github.com/plimble/ace"
	"github.com/plimble/ace-contrib/pongo2"
)

func main() {
	a := ace.New()

	render := pongo2.Pongo2(&pongo2.TemplateOptions{
		Directory:     "./public",
		Extensions:    []string{"html"},
		IsDevelopment: true, // if true enable cache template
	})
	a.HtmlTemplate(render)
	a.GET("/", homePage)

	a.Run(":5000")
}

func homePage(c *ace.C) {
	c.HTML("index.html", map[string]interface{}{"name": "john"})
}

```

######public/index.html

```
<html>
<head>
  <meta charset="UTF-8">
  <title>Ace Pongo2</title>
</head>
<body>
  <h1>Hello {{name}}</h1>
</body>
</html>
```