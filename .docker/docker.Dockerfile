FROM docker:latest

RUN apk add --no-cache curl openjdk17-jre

RUN mkdir -p /var/jenkins

WORKDIR /var/jenkins

CMD curl -sO http://jenkins:8080/jnlpJars/agent.jar && \
    java -jar agent.jar \
    -url http://jenkins:8080/ \
    -secret d92293fdd7f92e693caa25da5190823fe9e0e5dba552797a4ca83900ad9aa3c1 \
    -name "agent-docker" \
    -webSocket \
    -workDir "/var/jenkins"