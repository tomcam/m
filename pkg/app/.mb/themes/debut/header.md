{{ if .Site.Company.Name }}
{{- $name := .Site.Company.Name -}}
* [{{- $name -}}](/)
{{- else if .Site.Author.FullName -}}
{{- $name := .Site.Author.FullName -}}
* [{{ $name -}}](/)
{{- else }}
* [{{.Page.FrontMatter.Theme}}](#)
{{- end }} 
* [Gallery](/)
* [Docs](/)
* [FAQ](/)
* [About](/)

