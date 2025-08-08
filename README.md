
# CapyBelga

## Sobre o Projeto

CapyBelga é uma aplicação de exemplo desenvolvida para demonstrar o uso de OpenTelemetry em aplicações modernas. O projeto simula um sistema de clube de descontos, onde usuários podem se cadastrar, assinar clubes e cancelar inscrições, tudo monitorado com métricas, traces e logs via OpenTelemetry.

O objetivo é mostrar na prática como instrumentar uma aplicação Go para observabilidade, incluindo:
- Tracing distribuído
- Métricas customizadas (gauge, counter, histogram)
- Exportação para Prometheus, Jaeger, Grafana, etc.

## Como Usar

1. **Clone o repositório:**
   ```bash
   git clone https://github.com/hazkall/capy-belga.git
   cd capy-belga
   ```

2. **Suba os serviços com Docker Compose:**
   ```bash
   docker compose up postgres rabbitmq otel-collector prometheus jaeger grafana --build -d

   docker compose up capybelga --build
   ```

3. **Acesse os endpoints:**
   - Cadastro de usuário: `POST /contrate/discount-club/user`
   - Cadastro de clube: `POST /contrate/discount-club`
   - Inscrição em clube: `POST /contrate/discount-club/signup`
   - Cancelamento: `POST /user/cancel/club`
   - Estado do usuário: `GET /user/state`
   - Status do plano: `GET /user/plan/status`

4. **Observabilidade:**
   - As métricas e traces são exportados automaticamente para os backends configurados (Prometheus, Jaeger, etc). Consulte a documentação dos serviços para visualizar os dados.

## Imagem

<p align="center">
  <img src="./assets/capyhunter.png" width="350" alt="CapyBelga">
</p>

## Licença

MIT
