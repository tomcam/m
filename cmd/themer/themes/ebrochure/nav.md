[Contact](/)  [About](/) 
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

{{- if .Site.Social.GitHub }}[![GitHub logo](foo.svg)]({{ .Site.Social.GitHub }}){{ end -}} 

{{- if .Site.Social.Twitter }}[![Twitter logo](twitter-blue-30x30.svg)]({{ .Site.Social.Twitter }}){{ end -}}

{{- if .Site.Social.Facebook }}[![Facebook logo](facebook-blueish-30x30.svg)]({{ .Site.Social.Facebook }}){{ end -}} 

{{- if .Site.Social.Instagram }}[![Instagram logo](instagram-blue-30x30.svg)]({{ .Site.Social.Instagram }}){{ end -}} 



