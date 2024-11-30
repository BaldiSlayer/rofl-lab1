docker network create --gateway 172.16.1.1 --subnet 172.16.1.0/24 rofl-lab1
minikube start --driver=docker --network rofl-lab1 --ports=8443:30000 --listen-address=0.0.0.0
minikube kubectl -- create namespace argocd
minikube kubectl -- apply -n argocd -f https://raw.githubusercontent.com/argoproj/argo-cd/stable/manifests/install.yaml
watch -n 1 "minikube kubectl -- get pods -n argocd"

minikube kubectl -- port-forward --address 0.0.0.0 svc/argocd-server -n argocd 8966:443

minikube kubectl -- -n argocd get secret argocd-initial-admin-secret -o jsonpath="{.data.password}" | base64 -d

minikube kubectl -- create namespace prod

minikube kubectl -- create secret -n prod generic llm-secrets --from-literal=api-key="${MISTRAL_API_KEY}"
minikube kubectl -- create secret -n prod generic db-secrets --from-literal=username="${POSTGRES_USER}" --from-literal=password="${POSTGRES_PASSWORD}" --from-literal=db="${POSTGRES_DB}"
minikube kubectl -- create secret -n prod generic backend-secrets --from-literal=github-token="${GHTOKEN}" --from-literal=tg-api-key="${TGTOKEN}"

minikube service postgres-external -n prod --url
export $(cat .env | xargs) && docker run -it --rm --network rofl-lab1 -v $(pwd)/postgresql/migrations/:/migrations/migrations urbica/pgmigrate -d /migrations -t latest migrate -t 4 -c "port=30001 host=172.16.1.2 dbname=$POSTGRES_DB user=$POSTGRES_USER password=$POSTGRES_PASSWORD"
