
{{ if .FrontMatter.PageType }}
### This is the {{ .FrontMatter.PageType }} pagetype 
{{ else }}
### This is the {{ .FrontMatter.Theme }} theme
{{ end }}



