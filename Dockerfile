FROM golang
ADD . /go/src/github.com/tobyjsullivan/moneypenny
RUN  go install github.com/tobyjsullivan/moneypenny
CMD /go/bin/moneypenny
