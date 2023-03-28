/*
Package graceful is a Go 1.8+ package enabling graceful shutdown of http.Handler servers.

Usage:

  package main

  import (
      "log"

      "github.com/ETZhangSX/graceful"
      "github.com/gin-gonic/gin"
  )

  func main() {
      r := gin.Default()
      if err := graceful.ListenAndServe(":8080", r.Handler()); err != nil {
          log.Fatal(err)
      }
  }
*/
package graceful
