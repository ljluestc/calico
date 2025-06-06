# Copyright (c) 2015-2024 Tigera, Inc. All rights reserved.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

ARG CALICO_BASE

FROM registry.access.redhat.com/ubi9/ubi:latest AS ubi

RUN dnf upgrade -y

COPY docker-image/nginx.repo /etc/yum.repos.d/nginx.repo

RUN dnf --enablerepo=nginx-stable install -y \
    nginx

FROM scratch AS source

# binary and dependencies
COPY --from=ubi /usr/sbin/nginx /usr/sbin/nginx

COPY --from=ubi /bin/sh /bin/sh
COPY --from=ubi /usr/bin/cat /usr/bin/cat
COPY --from=ubi /usr/bin/cd /usr/bin/cd
COPY --from=ubi /usr/bin/mkdir /usr/bin/mkdir
COPY --from=ubi /usr/bin/coreutils /usr/bin/coreutils
COPY --from=ubi /usr/bin/sed /usr/bin/sed

COPY --from=ubi /lib64/libacl.so.1 /lib64/libacl.so.1
COPY --from=ubi /lib64/libattr.so.1 /lib64/libattr.so.1
COPY --from=ubi /lib64/libcap.so.2 /lib64/libcap.so.2
COPY --from=ubi /lib64/libcrypt.so.2 /lib64/libcrypt.so.2
COPY --from=ubi /lib64/libcrypto.so.3 /lib64/libcrypto.so.3
COPY --from=ubi /lib64/libpcre2-8.so.0 /lib64/libpcre2-8.so.0
COPY --from=ubi /lib64/libselinux.so.1 /lib64/libselinux.so.1
COPY --from=ubi /lib64/libssl.so.3 /lib64/libssl.so.3
COPY --from=ubi /lib64/libtinfo.so.6 /lib64/libtinfo.so.6
COPY --from=ubi /lib64/libz.so.1 /lib64/libz.so.1

# web server configs and folder permissions
COPY --from=ubi --chown=10001:10001 /etc/nginx /etc/nginx/
# nginx-start.sh writes clientsideConfig.js into this folder
COPY --from=ubi --chown=10001:10001 /usr/share/nginx/html /usr/share/nginx/html/
COPY --from=ubi --chown=10001:10001 /var/cache/nginx /var/cache/nginx/
COPY --from=ubi --chown=10001:10001 /var/log/nginx /var/log/nginx/

COPY docker-image/nginx-start.sh /usr/bin/nginx-start.sh

COPY docker-image/default.conf /etc/nginx/conf.d/default.conf
COPY docker-image/nginx.conf /etc/nginx/nginx.conf

COPY dist /usr/share/nginx/html/

COPY LICENSE /licenses/LICENSE

FROM ${CALICO_BASE}

ARG GIT_VERSION

LABEL org.opencontainers.image.description="Whisker provides a UI to view flows."
LABEL org.opencontainers.image.authors="maintainers@tigera.io"
LABEL org.opencontainers.image.source="https://github.com/projectcalico/calico"
LABEL org.opencontainers.image.title="Whisker"
LABEL org.opencontainers.image.vendor="Project Calico"
LABEL org.opencontainers.image.version="${GIT_VERSION}"
LABEL org.opencontainers.image.licenses="Apache-2.0"

COPY --from=source / /

USER 10001:10001

CMD ["/usr/bin/nginx-start.sh"]
