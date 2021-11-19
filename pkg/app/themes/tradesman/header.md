{{- /*  IMPORTANT: No need to change any of
        this manually. Just fill in the 
        appropriate parts of the site configuration file
        (probably found in .mb/site/site.yaml).
        ![Eastside Emerald Logo](eastside-emerald-64x64.png)

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
{{- else if .Page.Theme.Branding -}}
{{- $name := .Page.Theme.Branding -}}
* [{{ $name -}}](/)
{{- end }} 
* [Services](services.html)
* [Rates](rates.html)
* [Contact](contact.html)
* [About](about.html)

