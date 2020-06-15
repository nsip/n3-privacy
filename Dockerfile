FROM alpine
RUN mkdir n3-privacy
COPY ./n3-privacy /n3-privacy
WORKDIR /n3-privacy/Server/build/linux64
CMD ["./server"]

# FROM scratch
# ADD ./Server/build/linux64/server /
# ADD ./Server/build/linux64/config.toml /
# CMD ["./server"]

### docker build --tag=n3-privacy .

### docker tag IMAGE_ID cdutwhu/n3-privacy:latest
### docker login
### docker push cdutwhu/n3-privacy

### find ./ -type d -not -path "*/\.*"