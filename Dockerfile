from golang as builder
WORKDIR /go/src/github.com/mulesoft-consulting/terraform-provider-anypoint

COPY . .
RUN go get .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o terraform-provider-anypoint .