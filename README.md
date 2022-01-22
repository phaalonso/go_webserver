# Go Webserver
Um pequeno Web server que estou escrevendo em [Go Lang](https://go.dev/) enquanto me familizarizo com a linguagem. Este web server será inicialmente baseado na [wiki](https://go.dev/doc/articles/wiki/) sobre escrevendo aplicações web, no entanto existem planos para nos desviar deste.

- [X] Estrutura de dados para carregar e salvar métodos
- [X] Usando o pacote `net/http` para construir aplicações Web
- [X] Usando o pacote `html/template` para processar templates HTML
- [X] Usando cache nos templates
- [X] Usando regexp para validar entrada de dados
  - Corrigi problemas de segurança que permite o usuário criar uma página em qualquer caminho que desejar através de url's como: `http://localhost:8080/save/\.\./\.\./notvalidpath`
- [ ] Usando closures