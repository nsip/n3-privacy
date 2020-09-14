# FROM alpine
# RUN mkdir /n3-privacy
# COPY . / /n3-privacy/
# WORKDIR /n3-privacy/
# CMD ["./server"]

### docker build --tag=n3-privacy .

### ! run this docker image
### docker run --name privacy --net host n3-privacy:latest

### docker tag IMAGE_ID dockerhub-user/n3-privacy:latest
### docker login
### docker push dockerhub-user/n3-privacy

###########################
# INSTRUCTIONS
############################
# BUILD
#	docker build --rm -t nsip/n3-privacy:latest -t nsip/n3-privacy:v0.1.0 .
# TEST: docker run -it -v $PWD/test/data:/data -v $PWD/test/config.json:/config.json nsip/n3-privacy:develop .
# RUN: docker run -d nsip/n3-privacy:develop
#
# PUSH
#	Public:
#		docker push nsip/n3-privacy:v0.1.0
#		docker push nsip/n3-privacy:latest
#
#	Private:
#		docker tag nsip/n3-privacy:v0.1.0 the.hub.nsip.edu.au:3500/nsip/n3-privacy:v0.1.0
#		docker tag nsip/n3-privacy:latest the.hub.nsip.edu.au:3500/nsip/n3-privacy:latest
#		docker push the.hub.nsip.edu.au:3500/nsip/n3-privacy:v0.1.0
#		docker push the.hub.nsip.edu.au:3500/nsip/n3-privacy:latest
#
###########################
# DOCUMENTATION
############################


###########################
# STEP 0 Get them certificates
############################
# (note, step 2 is using alpine now) 
# FROM alpine:latest as certs

############################
# STEP 1 build executable binary (go.mod version)
############################
FROM golang:1.15.2-alpine3.12 as builder
RUN apk add --no-cache ca-certificates
RUN apk update && apk add --no-cache git bash
RUN mkdir -p /n3-privacy
COPY . / /n3-privacy/
WORKDIR /n3-privacy/
RUN ["/bin/bash", "-c", "./build_d.sh"]
RUN ["/bin/bash", "-c", "./release_d.sh"]

############################
# STEP 2 build a small image
############################
FROM alpine
COPY --from=builder /n3-privacy/app/ /
# NOTE - make sure it is the last build that still copies the files
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
WORKDIR /
CMD ["./server"]