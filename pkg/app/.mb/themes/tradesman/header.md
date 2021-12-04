{{- if .Site.Company.Name -}}
{{- $name := .Site.Company.Name -}}
{{- if .Site.Company.HeaderLogo -}}
* ![Logo]({{ .Site.Company.HeaderLogo }}) [{{ $name -}}](/)
{{- else -}}
* [{{ $name -}}](/)
{{- end }} 
{{- end }} 
* [Services](services.html)
* [Rates](rates.html)
* [Contact](contact.html)
* [About](about.html)

