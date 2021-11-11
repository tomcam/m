{{ if .Site.Company.Name }}
{{- $name := .Site.Company.Name -}}
{{ else if .Site.Author.FullName }}
{{- $name := .Site.Author.FullName -}} 
* [{{ $name }}](/)
{{ end }} 

