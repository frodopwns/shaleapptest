FROM quay.io/deis/base:v0.3.5

# Add user and group
RUN adduser --system \
	--shell /bin/bash \
	--disabled-password \
	--home /opt/shaleapptest \
	--group \
	shaleapptest

COPY . /

# Fix some permission since we'll be running as a non-root user
RUN mkdir -p /opt/shaleapptest/repository/cache \
    && chown -R shaleapptest:shaleapptest /opt/shaleapptest \
    && chmod +x /opt/shaleapptest/sbin/* 

USER shaleapptest
WORKDIR /opt/shaleapptest/sbin/
CMD ["./shaleapptest"]

ARG VERSION
ARG BUILD_DATE
ENV VERSION $VERSION
ENV BUILD_DATE $BUILD_DATE
ENV PATH $PATH:/opt/shaleapptest/sbin