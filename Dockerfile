FROM golang:1.19-alpine AS builder

RUN /sbin/apk update && \
	/sbin/apk --no-cache add ca-certificates git tzdata && \
	/usr/sbin/update-ca-certificates

RUN adduser -D -g '' greenlight

WORKDIR /home/greenlight

COPY . /home/greenlight

ARG VERSION
ARG BUILDTIME

RUN CGO_ENABLED=0 go build -a -tags netgo,osusergo \
    -ldflags "-extldflags '-static' -s -w" \
    -ldflags "-s -X main.buildTime=$BUILDTIME -X main.version=$VERSION" -o greenlight ./cmd/api

FROM busybox:musl

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /etc/group /etc/group
COPY --from=builder /home/greenlight/images /home/greenlight/images
COPY --from=builder /home/greenlight/migrations /home/greenlight/migrations
COPY --from=builder /home/greenlight/internal/mailer/templates /home/greenlight/internal/mailer/templates
COPY --from=builder /home/greenlight/vendor /home/greenlight/vendor
COPY --from=builder /home/greenlight/greenlight /home/greenlight/greenlight

RUN chown -R greenlight:greenlight /home/greenlight

USER greenlight
WORKDIR /home/greenlight

ENTRYPOINT ["/home/greenlight/greenlight"]