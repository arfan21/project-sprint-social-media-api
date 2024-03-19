package s3

import "go.opentelemetry.io/otel"

var tracer = otel.Tracer("github.com/arfan21/project-sprint-social-media-api/pkg/s3")
