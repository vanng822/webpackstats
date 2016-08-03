# webpackstats
  For rendering the bundle url generated by webpack based on metadata generated by https://github.com/owais/webpack-bundle-tracker

# Usage

  ```go
  templates := template.New("")
  templates = templates.Funcs(webpackstats.WebpackStats("/path/to/stats.json"))
  templates = template.Must(templates.ParseGlob("templates/*"))
  ```

  ```html
  <script type="text/javascript" src="{{webpackUrl "bundle"}}"></script>
  ```
