FROM scratch
MAINTAINER Red Hat, Inc. <container-tools@redhat.com>

ADD bin/atomicgo /atomicgo
CMD ["/atomicgo"]
