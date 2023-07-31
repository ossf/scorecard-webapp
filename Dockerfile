# Copyright 2021 OpenSSF Scorecard Authors
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#      http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

FROM golang:1.20.6@sha256:010a0ffe47398a3646993df44906c065c526eabf309d01fb0cbc9a5696024a60 AS builder
ENV APP_ROOT=/opt/app-root
ENV GOPATH=$APP_ROOT

WORKDIR $APP_ROOT/src/
ADD . $APP_ROOT/src/
RUN go mod download

ARG TARGETOS
ARG TARGETARCH
RUN CGO_ENABLED=0 make scorecard-webapp

# Multi-Stage production build
FROM golang:1.20.6@sha256:010a0ffe47398a3646993df44906c065c526eabf309d01fb0cbc9a5696024a60 as deploy
# Retrieve the binary from the previous stage
COPY --from=builder /opt/app-root/src/scorecard-webapp /usr/local/bin/scorecard-webapp

# Set the binary as the entrypoint of the container
ENTRYPOINT scorecard-webapp --host="0.0.0.0" --port=8080
EXPOSE 8080
