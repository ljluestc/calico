# Copyright (c) 2018-2025 Tigera, Inc. All rights reserved.
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

FROM scratch AS source

ARG TARGETARCH

COPY LICENSE /licenses/LICENSE

COPY bin/calico-typha-${TARGETARCH} /usr/bin/calico-typha
COPY typha.cfg /etc/calico/typha.cfg

FROM ${CALICO_BASE}

ARG GIT_VERSION=unknown

LABEL org.opencontainers.image.description="Calico Typha is a fan-out datastore proxy"
LABEL org.opencontainers.image.authors="maintainers@tigera.io"
LABEL org.opencontainers.image.source="https://github.com/projectcalico/calico"
LABEL org.opencontainers.image.title="Calico Typha"
LABEL org.opencontainers.image.vendor="Project Calico"
LABEL org.opencontainers.image.version="${GIT_VERSION}"
LABEL org.opencontainers.image.licenses="Apache-2.0"

COPY --from=source / /

USER 999

CMD ["/usr/bin/calico-typha"]
