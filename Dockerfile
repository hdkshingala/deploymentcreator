FROM golang:1.18.0
WORKDIR /deploymentcreator
ADD . .
RUN go mod download && CGO_ENABLED=0 go build

FROM scratch
WORKDIR /deploymentcreator
COPY --from=0 deploymentcreator .
ENTRYPOINT [ "./deploymentcreator" ]
