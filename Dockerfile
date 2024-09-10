FROM gsoci.azurecr.io/giantswarm/alpine:3.20.2

RUN apk update && apk --no-cache add ca-certificates && update-ca-certificates

RUN mkdir -p /usr/loca/bin/aws-fulfillment/content

ADD ./aws-fulfillment /usr/local/bin/aws-fulfillment/aws-fulfillment
ADD ./content /usr/local/bin/aws-fulfillment/content/

EXPOSE 8000
USER 9000:9000

WORKDIR /usr/local/bin/aws-fulfillment
ENTRYPOINT ["/usr/local/bin/aws-fulfillment/aws-fulfillment"]
