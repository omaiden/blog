FROM registry.moonrhythm.io/builder

ENV CGO_ENABLED=0

WORKDIR /workspace

ADD .tool-versions .
RUN asdf install

ADD go.mod go.sum ./
RUN go mod download
ADD . .
RUN go build -o .build/api -ldflags "-w -s" .

FROM gcr.io/distroless/static

WORKDIR /app

COPY --from=0 --link /workspace/.build/* ./
ENTRYPOINT ["/app/api"]
