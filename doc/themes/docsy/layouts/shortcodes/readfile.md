{{$file := .Get "file"}}
{{- if eq (.Get "markdown") "true" -}}
{{- $file | readFile | markdownify -}}
{{- else if  (.Get "highlight") -}}
{{- highlight ($file | readFile) (.Get "highlight") "" -}}
{{- else -}}
{{ $file | readFile | safeHTML }}
{{- end -}}
