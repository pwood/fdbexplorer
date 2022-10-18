FROM --platform=linux/amd64 golang:1.19

RUN curl -L https://github.com/apple/foundationdb/releases/download/7.1.22/foundationdb-clients_7.1.22-1_amd64.deb -o /tmp/fdbclient.deb && dpkg -i /tmp/fdbclient.deb && rm /tmp/fdbclient.deb
