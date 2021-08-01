FROM golang:1.16 AS builder

WORKDIR /app
COPY . .
RUN go build ./cmd/pdflatexserver

FROM debian:10.10

RUN apt-get update -q \
    && apt-get install -qy build-essential wget libfontconfig1 \
    && rm -rf /var/lib/apt/lists/*

# Install TexLive with scheme-basic
# https://github.com/blang/latex-docker/blob/master/Dockerfile.basic
RUN wget http://mirror.ctan.org/systems/texlive/tlnet/install-tl-unx.tar.gz; \
	mkdir /install-tl-unx; \
	tar -xvf install-tl-unx.tar.gz -C /install-tl-unx --strip-components=1; \
    echo "selected_scheme scheme-basic" >> /install-tl-unx/texlive.profile; \
	/install-tl-unx/install-tl -profile /install-tl-unx/texlive.profile; \
    rm -r /install-tl-unx; \
	rm install-tl-unx.tar.gz

COPY --from=builder /app/pdflatexserver /usr/bin/pdflatexserver

ENTRYPOINT ["/usr/bin/pdflatexserver"]
