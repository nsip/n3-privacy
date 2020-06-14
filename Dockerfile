FROM alpine

RUN mkdir -p /Server
COPY ./Server/build/linux64/server /Server
COPY ./Server/build/linux64/config.toml /Server

RUN mkdir -p /Enforcer
COPY ./Enforcer/build/linux64/enforcer /Enforcer

RUN mkdir -p /Client
COPY ./Client/build/linux64/client /Client
COPY ./Client/build/linux64/cfg-clt-privacy.toml /Client

WORKDIR /Server
CMD ["./server"]

# FROM scratch
# ADD ./Server/build/linux64/server /
# ADD ./Server/build/linux64/config.toml /
# CMD ["./server"]

### docker build --tag=n3-privacy .