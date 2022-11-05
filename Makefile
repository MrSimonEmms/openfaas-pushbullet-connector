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
