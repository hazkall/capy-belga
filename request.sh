#!/bin/bash

 # Criação de 1100 usuários
for i in {1..1100}
do
  curl -s -X POST http://localhost:8080/contrate/discount-club/user \
    -H "Content-Type: application/json" \
    -d "{\"name\": \"Usuario$i\", \"email\": \"usuario$i@email.com\"}"
done

echo "1100 usuários criados."

curl -s -X POST http://localhost:8080/contrate/discount-club \
  -H "Content-Type: application/json" \
  -d '{"name": "Clube Basic Online Web", "description": "Clube básico online website", "aquisition_channel": "online", "aquisition_location": "website", "plan_type": "basic"}'
curl -s -X POST http://localhost:8080/contrate/discount-club \
  -H "Content-Type: application/json" \
  -d '{"name": "Clube Basic Offline Store", "description": "Clube básico offline store", "aquisition_channel": "offline", "aquisition_location": "store", "plan_type": "basic"}'
curl -s -X POST http://localhost:8080/contrate/discount-club \
  -H "Content-Type: application/json" \
  -d '{"name": "Clube Premium Online Store", "description": "Clube premium online store", "aquisition_channel": "online", "aquisition_location": "store", "plan_type": "premium"}'
curl -s -X POST http://localhost:8080/contrate/discount-club \
  -H "Content-Type: application/json" \
  -d '{"name": "Clube Premium Offline Web", "description": "Clube premium offline website", "aquisition_channel": "offline", "aquisition_location": "website", "plan_type": "premium"}'
curl -s -X POST http://localhost:8080/contrate/discount-club \
  -H "Content-Type: application/json" \
  -d '{"name": "Clube Basic Offline Web", "description": "Clube básico offline website", "aquisition_channel": "offline", "aquisition_location": "website", "plan_type": "basic"}'

echo "5 clubes criados."

sleep 15 # Aguarda 2 segundos para garantir que os clubes foram persistidos no banco


CLUBES=("Clube Basic Online Web" "Clube Basic Offline Store" "Clube Premium Online Store" "Clube Premium Offline Web" "Clube Basic Offline Web")

for i in {1..1100}
do
  idx=$((RANDOM % 5))
  club="${CLUBES[$idx]}"
  if [ -z "$club" ]; then
    echo "[ERRO] Nome do clube vazio para usuario$i! Pulando..."
    continue
  fi
  echo "Inscrevendo usuario$i no clube: $club"
  curl -s -X POST http://localhost:8080/contrate/discount-club/signup \
    -H "Content-Type: application/json" \
    -d "{\"email\": \"usuario$i@email.com\", \"club\": \"$club\"}"
done

echo "1100 inscrições criadas."


curl -X GET http://localhost:8080/user/state \
-H "Content-Type: application/json" \
-d '{"email": "usuario1@email.com", "name": "Usuario1"}'


curl -X POST http://localhost:8080/user/cancel/club \
-H "Content-Type: application/json" \
-d '{"email": "usuario1@email.com"}'

curl -X POST http://localhost:8080/user/cancel/club \
-H "Content-Type: application/json" \
-d '{"email": "usuario47@email.com"}'


curl -X POST http://localhost:8080/user/cancel/club \
-H "Content-Type: application/json" \
-d '{"email": "usuario23@email.com"}'

curl -X GET http://localhost:8080/user/plan/status \
  -H "Content-Type: application/json" \
  -d '{"email": "usuario15@email.com"}'
