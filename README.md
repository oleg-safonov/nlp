# nlp

Библиотека реализует лемматизацию текста на русском языке. Пользоваться максимально просто:

```go
package main

import (
    "fmt"
    "log"

    "github.com/oleg-safonov/nlp"
    nlprudata "github.com/oleg-safonov/nlp-ru-data"
)

func main() {
    base, err := nlprudata.Load()
    if err != nil {
        log.Fatal(err)
    }

    lem, err := nlp.NewLemmatizer(base)
    if err != nil {
        log.Fatal(err)
    }

    fmt.Println(lem.LemmatizeText("Съешь ещё этих мягких французских булок да выпей чаю"))
}
```

Результат:

```
[съесть еще этот мягкий французский булка да выпить чай]
```