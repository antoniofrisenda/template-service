FROM docker:latest

RUN apk add --no-cache curl openjdk17-jre

RUN mkdir -p /home/jenkins/agent

WORKDIR /home/jenkins

CMD curl -sO http://jenkins:8080/jnlpJars/agent.jar && \
    java -jar agent.jar \
    -url http://jenkins:8080/ \
    -secret e9c394a719ae6b048db0c747769abdad9527e2c5dfb2b809b0a8ab595103073b \
    -name "agent-docker" \
    -webSocket \
    -workDir "/home/jenkins/agent"