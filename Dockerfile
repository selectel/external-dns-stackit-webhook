FROM gcr.io/distroless/static-debian11:nonroot

COPY external-dns-webhook /external-dns-webhook

ENTRYPOINT ["/external-dns-webhook"]
