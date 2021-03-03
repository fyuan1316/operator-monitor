# Build the manager binary
FROM harbor-b.alauda.cn/asm/builder:0.4-alpine3.12.1 AS builder

COPY ./bin/ /opt/

RUN ARCH="" && dpkgArch="$(arch)" \
  && case "${dpkgArch}" in \
  x86_64) ARCH='amd64' && upx /opt/${ARCH}/manager ;; \
  aarch64) ARCH='arm64' && upx /opt/${ARCH}/manager  ;; \
  *) echo "unsupported architecture"; exit 1 ;; \
  esac \
  && cp /opt/${ARCH}/manager /manager


FROM harbor-b.alauda.cn/asm/runner:0.1-alpine3.12.1
RUN apk --no-cache --update add ca-certificates
WORKDIR /
COPY --from=builder /manager /manager

ENTRYPOINT ["/manager"]
