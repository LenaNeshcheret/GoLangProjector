FROM golang

WORKDIR /bin

COPY . .
RUN go get -d -v ./...
RUN go build -o api .

CMD [ "/bin/api" ]
