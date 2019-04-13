from golang:1.12.4-stretch  as builder
WORKDIR /app/terraform-provider-anypoint
ENV GO111MODULE=on
ENV TERRAFORM_VER 0.11.7

COPY . .
#RUN go get .
#RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o ./build/terraform-provider-anypoint .
RUN ( \
    apt-get update && \
    apt-get install --yes --no-install-recommends \
    curl unzip python python-pip curl unzip groff && \
    apt-get clean && \
    rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/* \
    )

RUN ( \
    curl -fsSL -o /tmp/terraform.zip https://releases.hashicorp.com/terraform/${TERRAFORM_VER}/terraform_${TERRAFORM_VER}_linux_amd64.zip && \
    unzip /tmp/terraform.zip -d /usr/local/bin/ && \
    rm -f /tmp/terraform.zip \
    )

RUN make -f ./build/GNUmakefile test
RUN make -f ./build/GNUmakefile build

#Now build the final image containing just the executable
FROM scratch
COPY --from=BUILDER /usr/local/bin/terraform /usr/local/bin/terraform
COPY --from=builder /app/terraform-provider-anypoint/build/terraform-provider-anypoint /usr/local/bin/terraform-provider-anypoint
ENTRYPOINT ["/usr/local/bin/terraform"]