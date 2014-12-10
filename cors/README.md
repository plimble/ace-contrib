# ace-cors

CORS middleware for [Gin].

## Installation

``` bash
$ go get github.com/plimble/ace/cors
```

## Usage

``` go
import (
    "github.com/gin-gonic/gin"
    "github.com/plimble/ace/cors"
)

func main(){
    g := gin.New()
    g.Use(cors.Middleware(cors.Options{}))
}
```

[Gin]: http://gin-gonic.github.io/gin/
