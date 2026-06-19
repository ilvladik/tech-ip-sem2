Set-Location "$PSScriptRoot\deploy\rabbit"
echo "Starting RabbitMQ AMQP :5672, UI http://localhost:15672 (guest/guest)"
docker compose up -d
docker compose ps
