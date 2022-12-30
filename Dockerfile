FROM golang:1.18


WORKDIR /app

COPY . .



RUN  go build  -o main .

RUN mkdir temp 

 
EXPOSE 5000

CMD [ "./main" ]
