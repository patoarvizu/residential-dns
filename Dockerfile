FROM golang:1.13 as builder
ARG TARGETARCH
ARG TARGETVARIANT

WORKDIR /workspace

COPY go.mod go.mod
COPY go.sum go.sum

RUN go mod download

COPY main.go main.go

RUN CGO_ENABLED=0 GOOS=linux GOARM=$(if [ "$TARGETVARIANT" = "v7" ]; then echo "7"; fi) GOARCH=$TARGETARCH GO111MODULE=on go build -a -o residential-dns main.go

FROM gcr.io/distroless/static:nonroot-amd64

ARG GIT_COMMIT="unspecified"
LABEL GIT_COMMIT=$GIT_COMMIT

ARG GIT_TAG=""
LABEL GIT_TAG=$GIT_TAG

ARG COMMIT_TIMESTAMP="unspecified"
LABEL COMMIT_TIMESTAMP=$COMMIT_TIMESTAMP

ARG AUTHOR_EMAIL="unspecified"
LABEL AUTHOR_EMAIL=$AUTHOR_EMAIL

ARG SIGNATURE_KEY="undefined"
LABEL SIGNATURE_KEY=$SIGNATURE_KEY

WORKDIR /
COPY --from=builder /workspace/residential-dns .
USER nonroot:nonroot

ENTRYPOINT ["/residential-dns"]
