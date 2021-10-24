{{- /*  Automatically name first item in header    
        based on company name, author name name
        if no company was specified, or just 
        the name of the theme if neither of those
        was specified.
        
*/ -}}
{{- if .Site.Company.Name -}}
{{- $name := .Site.Company.Name -}}
* [{{ $name -}}](/)
{{- else if .Page.Theme.PageType.Branding -}}
{{- $name := .Page.Theme.PageType.Branding -}}
* [{{ $name -}}](/)
{{- else if .Site.Author.FullName -}}
{{- $name := .Site.Author.FullName -}}
* [{{ $name -}}](/)
{{- else }}
* [{{.FrontMatter.Theme}} {{.FrontMatter.PageType}}](/)
{{- end }} 
* [Events](/)
* [Podcast](/)
* [Subscribe](/)

