package docvu

import (
  "fmt"
)

func ExampleWordCount() {
  fmt.Println(WordCount("hello"))
  fmt.Println(WordCount("hello goodbye"))
  fmt.Println(WordCount("hello, goodbye"))
  fmt.Println(WordCount("hello + goodbye"))
  fmt.Println(WordCount("hello & goodbye"))
  fmt.Println(WordCount("hello - goodbye"))
  fmt.Println(WordCount("hello-goodbye"))
  fmt.Println(WordCount("hello,goodbye"))
  fmt.Println(WordCount("hello--goodbye"))
  fmt.Println(WordCount("hello.goodbye"))
  fmt.Println(WordCount("hello—goodbye")) // em-dash

  // Output:
  // 1
  // 2
  // 2
  // 3
  // 3
  // 2
  // 1
  // 2
  // 2
  // 2
  // 2
}