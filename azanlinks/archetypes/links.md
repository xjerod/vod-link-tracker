---
title: "{{ .Name }}"
date: {{ .Date.Format "2006-01-02T15:04:05-07:00" }}
draft: true
---

# Links
{{ range .Links }}
- [{{ . }}](https://{{ . }})
{{ end }}
