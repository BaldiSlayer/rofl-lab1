kubectl create secret -n prod generic llm-secrets --from-literal=api-key="${MISTRAL_API}"
kubectl create secret -n prod generic db-secrets --from-literal=username="${DB_USERNAME}" --from-literal=password="${DB_PASSWORD}" --from-literal=db="${DB_NAME}"
kubectl create secret -n prod generic backend-secrets --from-literal=github-token="${GITHUB_TOKEN}" --from-literal=tg-api-key="${TG_API_KEY}"
