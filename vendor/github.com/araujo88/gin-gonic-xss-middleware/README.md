# XSS Middleware 

XssMw is an middleware written in [Golang](https://golang.org/) for the [Gin](https://github.com/gin-gonic/gin) web framework. Although, it should be useable with any Go web framework which utilizes Golang's "net/http" native library in a similiar way to Gin. The idea behind XssMw is to "auto remove XSS" from user submitted input. It's applied on http GET, POST, PATCH and PUT requests only.

The XSS filtering is performed by HTML sanitizer [Bluemonday](https://github.com/microcosm-cc/bluemonday).

The default is to the strictest policy - StrictPolicy()

## Installation

`go get -u github.com/araujo88/gin-gonic-xss-middleware` <br>
`go mod tidy`

## How To Use it?

Using the defaults,
It will skip filtering for a field named 'password' but will run the filter on everything else.
Uses the Bluemonday strictest policy - StrictPolicy()

```go
package main

import (
    xss "github.com/araujo88/gin-gonic-xss-middleware"
    "github.com/gin-gonic/gin"
)

func main() {
    r := gin.Default()

    // include as standard middleware
    var xssMdlwr xss.XssMw
    r.Use(xssMdlwr.RemoveXss())

    r.GET("/ping", func(c *gin.Context) {
        c.JSON(200, gin.H{
            "message": "pong",
        })
    })
    r.Run() // listen and serve on 0.0.0.0:8080
}

```

Using some config options, here It will skip filtering for a fields named 'password', "create_date" and "token" but will run the filter  on everything else.

Uses Bluemonday UGCPolicy


```go
package main

import (
    xss "github.com/araujo88/gin-gonic-xss-middleware"
    "github.com/gin-gonic/gin"
)

func main() {
    r := gin.Default()

    xssMdlwr := &xss.XssMw{
            FieldsToSkip: []string{"password", "create_date", "token"},
            BmPolicy:     "UGCPolicy",
    }
    r.Use(xssMdlwr.RemoveXss())

    r.GET("/ping", func(c *gin.Context) {
        c.JSON(200, gin.H{
            "message": "pong",
        })
    })
    r.Run() // listen and serve on 0.0.0.0:8080
}

```

## Data

Currently, it removes (deletes) all HTML and malicious detected input from user input on  the submitted request to the server. It handles three Request types:

* JSON requests - Content-Type application/json
* Form Encoded - Content-Type application/x-www-form-urlencoded
* Multipart Form Data - Content-Type multipart/form-data

A future plan is have a feature to store all user submitted data intact and have the option to filter it out on the http response, so you can choose your preference. In other words: data would be stored in the database as it was submitted and removed in Responses back to the user. Pros: data integrity. Cons: XSS exploits still present.

## NOTE: This is beta level code with minimal actual real world usage

## Contributing 

 - You are welcome to contribute to this project. Please open a issue first for discussing before opening a pull request.
 - Please update/add tests as appropriate.
 - Send pull request against the develop branch.
 - Please use the same formatting as the Go authors. Run code through gofmt before submitting. 

## Acknowledgements

Thanks to dvwright (https://github.com/dvwright/xss-mw) from which this repository was forked.

https://github.com/goware/jsonp <br>
https://github.com/appleboy/gin-jwt/tree/v2.1.1 <br>
https://github.com/microcosm-cc/bluemonday <br>

Learning source: https://static.googleusercontent.com/intl/hu/about/appsecurity/learning/xss/

> A note on manually escaping input
> Writing your own code for escaping input and then properly and consistently applying it is extremely difficult. 
> We do not recommend that you manually escape user-supplied data. Instead, we strongly recommend that you 
> use a templating system or web development framework that provides context-aware auto-escaping. 

