FROM ubuntu:xenial

MAINTAINER Rene Kaufmann <kaufmann.r@gmail.com>

ENV GOPATH /go
ENV PATH $GOPATH/bin:/usr/local/go/bin:$PATH
ENV GO15VENDOREXPERIMENT 1

ADD . /go/src/github.com/HeavyHorst/unoconv-api

#Install unoconv
RUN \
	apt-get update && \
	DEBIAN_FRONTEND=noninteractive \
	    apt-get upgrade -y && \
		apt-get install -y \
		        locales \
			unoconv \
			gcc \
			supervisor \
			golang-go && \
		go install github.com/HeavyHorst/unoconv-api && \
        apt-get remove -y golang-go gcc && \
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

# Startup
ENTRYPOINT ["/usr/bin/supervisord"]
