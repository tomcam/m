{{- if .Site.Company.Name -}}
{{- $name := .Site.Company.Name -}}
* [{{ $name -}}](/)
{{- else if .Site.Author.FullName -}}
{{- $name := .Site.Author.FullName -}}
* [{{ $name -}}](/)
{{- else if .Page.Theme.Branding -}}
{{- $name := .Page.Theme.Branding -}}
* [{{ $name -}}](/)
{{- end }} 
* [Services](services.html)
* [Rates](rates.html)
* [Contact](contact.html)
* [About](about.html)

