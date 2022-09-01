# Overview

This project is wrapped apach bench mark. This makes apache bench easy to run on multiple cores.

# Prerequisites

- ab

# Example

```
./ab-wrap-go -goroutineNum 5 -n 10 -c 10 -url https://www.example.com/
```
