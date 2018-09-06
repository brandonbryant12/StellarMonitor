FROM golang:1.9
RUN apt-get update && apt-get install -y unzip --no-install-recommends && \
    apt-get autoremove -y && apt-get clean -y && \
    wget -O dep.zip https://github.com/golang/dep/releases/download/v0.3.0/dep-linux-amd64.zip && \
    echo '96c191251164b1404332793fb7d1e5d8de2641706b128bf8d65772363758f364  dep.zip' | sha256sum -c - && \
    unzip -d /usr/bin dep.zip && rm dep.zip

RUN mkdir -p /go/src/github.com/app
WORKDIR /go/src/github.com/app
COPY . . 
RUN dep ensure -vendor-only
RUN go build BlockchainListener.go
RUN go build PaymentListener.go Ledger.go Payment.go
