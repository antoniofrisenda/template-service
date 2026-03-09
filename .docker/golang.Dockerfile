FROM golang:latest

RUN apt-get update && \
    apt-get install -y curl default-jre && \
    rm -rf /var/lib/apt/lists/*

RUN mkdir -p /home/jenkins/agent

WORKDIR /home/jenkins

CMD curl -sO http://jenkins:8080/jnlpJars/agent.jar && \
    java -jar agent.jar \
    -url http://jenkins:8080/ \
    -secret 160628eb86fca6fb9c33b5570733de32cfa2730a3f5cf85b29f8a5e2c4451f9c \
    -name "agent-golang" \
    -webSocket \
    -workDir "/home/jenkins/agent"