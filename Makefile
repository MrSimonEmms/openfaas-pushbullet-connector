HELM_CHART_NAME ?= pushbullet-connector
OPENFAAS_NAMESPACE ?= openfaas

helm_install:
	@kubectl create secret generic \
		-n ${OPENFAAS_NAMESPACE} \
		pushbullet \
		--from-literal=token=${PUSHBULLET_TOKEN} \
		--dry-run=client -o yaml | \
		kubectl replace --force -f -

	@helm upgrade \
		--atomic \
		--cleanup-on-fail \
		--create-namespace \
		--install \
		--namespace="${OPENFAAS_NAMESPACE}" \
		--reset-values \
		--wait \
		${HELM_CHART_NAME} \
		./chart/openfaas-pushbullet-connector
.PHONY: helm_install

install_openfaas:
	@echo "Installing OpenFaaS"

	@kubectl apply -f https://raw.githubusercontent.com/openfaas/faas-netes/master/namespaces.yml

	@helm upgrade \
		--atomic \
		--cleanup-on-fail \
		--install \
		--namespace openfaas \
		--repo https://openfaas.github.io/faas-netes \
		--reset-values \
		--set functionNamespace=openfaas-fn \
		--set generateBasicAuth=true \
		--wait \
		openfaas \
		openfaas
.PHONY: install_openfaas

run:
	go run . run
.PHONY: run
