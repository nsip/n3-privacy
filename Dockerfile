FROM alpine
RUN mkdir /n3-privacy
COPY . / /n3-privacy/
WORKDIR /n3-privacy/
CMD ["./server"]

# FROM scratch
# ADD ./Server/build/linux64/server /
# ADD ./Server/build/linux64/config.toml /
# CMD ["./server"]

### docker build --tag=n3-privacy .

### ! run this docker image
### docker run --name privacy --net host n3-privacy:latest

### docker tag IMAGE_ID dockerhub-user/n3-privacy:latest
### docker login
### docker push dockerhub-user/n3-privacy
