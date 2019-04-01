from golang as builder
WORKDIR /go/src/github.com/tech-nico/terraform-provider-anypoint
ENV GO111MODULE=on

COPY . .
RUN go get .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o ./build/terraform-provider-anypoint .
ENTRYPOINT ["tail", "-f", "/dev/null"]
