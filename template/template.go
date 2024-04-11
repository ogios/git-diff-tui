package template

import (
	"encoding/json"
	"os"
	"strings"
)

const A = `
<!DOCTYPE html>
<html lang="en">
  <head>
    <meta charset="UTF-8" />
    <meta name="viewport" content="width=device-width, initial-scale=1" />
    <title></title>
    <!-- <link href="css/style.css" rel="stylesheet" /> -->
    <link href="https://unpkg.com/gridjs/dist/theme/mermaid.min.css" rel="stylesheet" />
    <style>
    td {
    white-space: pre-wrap;
    }
    </style>
  </head>
  <body>
    <div id="wrapper"></div>

    <script src="https://unpkg.com/gridjs/dist/gridjs.umd.js"></script>
    <script>
      const DATA = data-replace
      new gridjs.Grid({
        // columns: ['name', 'hash', 'comment', 'time', 'tag', 'author'],
        columns: ['name', 'comment'],
        data: DATA,
        // autoWidth: false,
      }).render(document.getElementById('wrapper'))
    </script>
  </body>
</html>
`

/*
生成html文件

传any, 里面转json然后替换模板
*/
func GenTemplate(data any) {
	s, err := json.Marshal(data)
	if err != nil {
		panic(err)
	}
	f, err := os.Create("index.html")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	_, err = f.WriteString(strings.Replace(A, "data-replace", string(s), 1))
	if err != nil {
		panic(err)
	}
}
