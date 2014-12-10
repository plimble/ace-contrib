# gin-sessions 
Gin middleware/handler for easy session management.

## Usage

~~~ go
package main

import (
  "github.com/gin-gonic/gin"
  "github.com/plimble/ace/sessions"
)

func main() {
  n := gin.New()

  store := sessions.NewCookieStore([]byte("secret123"))  
  n.Use(sessions.Sessions("my_session", store))

  n.GET("/show", func(c *gin.Context) {
    session := sessions.GetSession(c.Request)
    session.Set("hello", "world")
  })

  n.Run(":3000")
}

~~~
