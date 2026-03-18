FROM golang:latest

RUN apt-get update && \
    apt-get install -y curl default-jre docker.io ca-certificates unzip && \
    curl -L -o /tmp/sonar-scanner.zip \
      https://binaries.sonarsource.com/Distribution/sonar-scanner-cli/sonar-scanner-cli-8.0.1.6346-linux-x64.zip && \
    unzip /tmp/sonar-scanner.zip -d /opt && \
    mv /opt/sonar-scanner-8.0.1.6346-linux-x64 /opt/sonar-scanner && \
    rm /tmp/sonar-scanner.zip && \
    rm -rf /var/lib/apt/lists/*

ENV SONAR_SCANNER_HOME=/opt/sonar-scanner \
    PATH="/opt/sonar-scanner/bin:${PATH}"

RUN mkdir -p /var/jenkins

WORKDIR /var/jenkins

CMD curl -sO http://jenkins:8080/jnlpJars/agent.jar && \
    java -jar agent.jar \
    -url http://jenkins:8080/ \
    -secret cedb3609bb8426abb2e2a666d174ce8731d00287180a39aeeb54e3faccc2c158 \
    -name "agent-golang" \
    -webSocket \
    -workDir "/var/jenkins"
