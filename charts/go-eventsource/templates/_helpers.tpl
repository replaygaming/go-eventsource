{{/* vim: set filetype=mustache: */}}
{{/*
Expand the name of the chart.
*/}}
{{- define "go-eventsource.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" -}}
{{- end -}}

{{/*
Create a default fully qualified app name.
We truncate at 63 chars because some Kubernetes name fields are limited to this (by the DNS naming spec).
If release name contains chart name it will be used as a full name.
*/}}
{{- define "go-eventsource.fullname" -}}
{{- if .Values.fullnameOverride -}}
{{- .Values.fullnameOverride | trunc 63 | trimSuffix "-" -}}
{{- else -}}
{{- $name := default .Chart.Name .Values.nameOverride -}}
{{- if contains $name .Release.Name -}}
{{- .Release.Name | trunc 63 | trimSuffix "-" -}}
{{- else -}}
{{- printf "%s-%s" .Release.Name $name | trunc 63 | trimSuffix "-" -}}
{{- end -}}
{{- end -}}
{{- end -}}

{{/*
Create chart name and version as used by the chart label.
*/}}
{{- define "go-eventsource.chart" -}}
{{- printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" -}}
{{- end -}}

{{- define "go-eventsource.affinity" -}}
{{- if .Values.affinity.enabled -}}
affinity:
  nodeAffinity:
    {{- if eq .Values.affinity.type "soft" }}
    preferredDuringSchedulingIgnoredDuringExecution:
    - weight: 100
      preference:
        matchExpressions:
        - key: {{ .Values.affinity.labelKey }}
          operator: In
          values:
          - {{ .Values.affinity.labelValue }}
    {{- else if eq .Values.affinity.type "hard" }}
    requiredDuringSchedulingIgnoredDuringExecution:
      nodeSelectorTerms:
      - matchExpressions:
        - key: {{ .Values.affinity.labelKey }}
          operator: In
          values:
          - {{ .Values.affinity.labelValue }}
    {{- end }}
{{- end -}}
{{- end -}}

{{- define "go-eventsource.image" -}}
{{ list .Values.image.name .Values.image.tag | join ":" }}
{{- end -}}
