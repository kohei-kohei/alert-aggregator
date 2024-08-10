FROM golang:1.22.6-bookworm AS builder

ENV CGO_ENABLED=0

RUN --mount=type=bind,rw,target=. go build -o /bin/main


FROM gcr.io/distroless/static-debian12 AS executor

USER nonroot

COPY --from=builder /bin/main /main

CMD [ "/main" ]
