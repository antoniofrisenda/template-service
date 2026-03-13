FROM jenkins/inbound-agent:latest

USER root

RUN apt-get update && \
    apt-get install -y curl ca-certificates docker.io && \
    curl -LO "https://dl.k8s.io/release/$(curl -L -s https://dl.k8s.io/release/stable.txt)/bin/linux/amd64/kubectl" && \
    install -o root -g root -m 0755 kubectl /usr/local/bin/kubectl && \
    rm kubectl && \
    rm -rf /var/lib/apt/lists/*

RUN mkdir -p /var/jenkins && chown jenkins:jenkins /var/jenkins
USER jenkins