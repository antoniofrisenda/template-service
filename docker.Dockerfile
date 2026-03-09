FROM docker:latest

RUN apk add --no-cache curl openjdk17-jre

RUN mkdir -p /home/jenkins/agent

WORKDIR /home/jenkins

CMD curl -sO http://jenkins:8080/jnlpJars/agent.jar && \
    java -jar agent.jar \
    -url http://jenkins:8080/ \
    -secret f9c920548f870c737e07634460719649ee69543504b75abac24b2de7b64f2caa \
    -name "docker-agent" \
    -webSocket \
    -workDir "/home/jenkins/agent"