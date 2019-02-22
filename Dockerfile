FROM golang:1.11 as api-builder
WORKDIR /unoconv-api
COPY . /unoconv-api
RUN go build


FROM ubuntu:xenial

LABEL maintainer="kaufmann.r@gmail.com"

COPY --from=api-builder /unoconv-api/unoconv-api /opt/unoconv-api/unoconv-api

#Install unoconv
RUN \
	apt-get update && \
	DEBIAN_FRONTEND=noninteractive \
	    apt-get upgrade -y && \
		apt-get install -y \
		        locales \
			unoconv \
			supervisor && \
        apt-get remove -y && \
	    apt-get autoremove -y && \
        apt-get clean && \
			rm -rf /var/lib/apt/lists/

# Set the locale
RUN locale-gen de_DE.UTF-8  
ENV LANG de_DE.UTF-8  
ENV LANGUAGE de_DE:de  
ENV LC_ALL de_DE.UTF-8  

COPY supervisord.conf /etc/supervisor/conf.d/supervisord.conf

# Expose 3000
EXPOSE 3000

HEALTHCHECK --interval=30s --timeout=10s \
    CMD curl http://localhost:3000/unoconv/health

# Startup
ENTRYPOINT ["/usr/bin/supervisord"]
