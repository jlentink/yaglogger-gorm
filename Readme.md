# Yaglogger wrapper for Gorm

This wrapper allows for easy usage of YagLogger with Gorm.
Ensuring all can be used via the same logger.

This is a different package than the normal YagLogger package to ensure
the dependencies of the original logger package is small as possible.

example usage:
```go
import (
    "gorm.io/driver/postgres"
    "gorm.io/gorm"
    "github.com/jlentink/yaglogger"
    yagloggerGorm "github.com/jlentink/yaglogger-gorm"
)

db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{Logger: .NewLogger()})
```