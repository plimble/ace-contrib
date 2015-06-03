# ace-cors

CORS middleware for [ACE](https://github.com/plimble/ace)

## Installation

``` bash
$ go get github.com/plimble/ace-contrib/cors
```

## Usage

``` go
import (
    "github.com/plimble/ace"
    "github.com/plimble/ace-contrib/cors"
)

func main(){
    a := ace.New()
    a.Use(cors.Cors(cors.Options{}))
}
```

[ACE]: https://github.com/plimble/ace
