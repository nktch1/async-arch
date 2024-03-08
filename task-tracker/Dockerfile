FROM golang:1.22

ENV GOBIN=/go/bin

COPY . /task-tracker-sources

RUN make install -C /task-tracker-sources

ENTRYPOINT ${GOBIN}/task-tracker
