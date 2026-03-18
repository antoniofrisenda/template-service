FROM jenkins/inbound-agent:latest

USER root

RUN apt-get update && \
    apt-get install -y --no-install-recommends curl ca-certificates unzip && \
    curl -L -o /tmp/sonar-scanner.zip \
      https://binaries.sonarsource.com/Distribution/sonar-scanner-cli/sonar-scanner-cli-8.0.1.6346-linux-x64.zip && \
    unzip /tmp/sonar-scanner.zip -d /opt && \
    mv /opt/sonar-scanner-8.0.1.6346-linux-x64 /opt/sonar-scanner && \
    rm /tmp/sonar-scanner.zip && \
    rm -rf /var/lib/apt/lists/*

ENV SONAR_SCANNER_HOME=/opt/sonar-scanner \
    PATH="/opt/sonar-scanner/bin:${PATH}"

RUN mkdir -p /var/jenkins && chown -R jenkins:jenkins /var/jenkins
WORKDIR /var/jenkins
USER jenkins
