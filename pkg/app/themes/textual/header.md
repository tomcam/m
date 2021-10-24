{{- /*  IMPORTANT: No need to change any of
        this manually. Just fill in the 
        appropriate parts of the site.toml file
        (probably found in .mb/site/site.toml).

        Automatically name first item in header    
        based on company name, then author name.
        If those fail, use the the branding name 
        of the theme or just the theme name if
        neither of those was specified.
        
        ONE MORE NOTE: None of this is required
        by the theme. You can just replace it with
        whatever text or Markdown you please.
*/ -}}
{{- if .Site.Company.Name -}}
{{- $name := .Site.Company.Name -}}
* [{{ $name -}}](/)
{{- else if .Site.Author.FullName -}}
{{- $name := .Site.Author.FullName -}}
* [{{ $name -}}](/)
{{- else if .Page.Theme.PageType.Branding -}}
{{- $name := .Page.Theme.PageType.Branding -}}
* [{{ $name -}}](/)
{{- else }}
* [{{.FrontMatter.Theme}} {{.FrontMatter.PageType}}](/)
{{- end }} 
* [Events](/)
* [Podcast](/)
* [Subscribe](/)

