FROM docker.io/library/golang:trixie AS build
ADD . /go/src/app
WORKDIR /go/src/app
RUN make onebill

FROM gcr.io/distroless/base-debian13:nonroot

LABEL org.opencontainers.image.title=OneBill
LABEL org.opencontainers.image.description="OneBill is a minimalist and opinionated feed reader"
LABEL org.opencontainers.image.vendor="Frédéric Guillot"
LABEL org.opencontainers.image.licenses=Apache-2.0
LABEL org.opencontainers.image.url=https://github.com/onebill/onebill
LABEL org.opencontainers.image.source=https://github.com/onebill/onebill

EXPOSE 8080
ENV LISTEN_ADDR=0.0.0.0:8080
COPY --from=build /go/src/app/onebill /usr/bin/onebill
CMD ["/usr/bin/onebill"]
