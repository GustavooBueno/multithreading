# Multithreading - resultado mais rápido entre duas APIs distintas.

As duas requisições são feitas simultaneamente para as seguintes APIs:

https://brasilapi.com.br/api/cep/v1/01153000 + cep

http://viacep.com.br/ws/" + cep + "/json/

- Acata a API que entrega a resposta mais rápida e descarta a resposta mais lenta.

- O resultado da request é exibido no command line com os dados do endereço, bem como qual API a enviou.

- Limita o tempo de resposta em 1 segundo. Caso contrário, o erro de timeout é exibido.

# Como rodar?
- basta clonar o projeto
- rodar o comando: go run main.go
