FROM golang:latest

RUN apt-get update && \
    apt-get install -y curl default-jre docker.io && \
    rm -rf /var/lib/apt/lists/*

RUN mkdir -p /var/jenkins

WORKDIR /var/jenkins

CMD curl -sO http://jenkins:8080/jnlpJars/agent.jar && \
    java -jar agent.jar \
    -url http://jenkins:8080/ \
    -secret cedb3609bb8426abb2e2a666d174ce8731d00287180a39aeeb54e3faccc2c158 \
    -name "agent-golang" \
    -webSocket \
    -workDir "/var/jenkins"