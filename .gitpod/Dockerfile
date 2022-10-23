FROM gitpod/workspace-full
USER root
# Install Kubectl
RUN curl -LO "https://dl.k8s.io/release/$(curl -L -s https://dl.k8s.io/release/stable.txt)/bin/linux/amd64/kubectl" \
  && install -o root -g root -m 0755 kubectl /usr/local/bin/kubectl \
  && kubectl completion bash > /etc/bash_completion.d/kubectl \
  && kubectl version --client
# Install Helm
RUN curl https://raw.githubusercontent.com/helm/helm/main/scripts/get-helm-3 | bash \
  && helm completion bash > /etc/bash_completion.d/helm \
  && helm version
# Install faas-cli
RUN curl -sSL https://cli.openfaas.com | sh
USER gitpod
