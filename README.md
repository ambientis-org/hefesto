# Hefesto

## Backend

### Instalación
Para ejecutar este servicio de manera sencilla se recomienda
instalar Docker y docker-commpose, para más información, revisar

 - [Instalación de Docker](https://docs.docker.com/install/)

 - [Instalación de docker-compose](https://docs.docker.com/compose/install/)

Una vez instaladas ambas herraminentas, ejecutar en una terminal

```
$ git clone https://github.com/ambientis-org/hefesto.git
$ cd hefesto
$ docker-compose build && docker-compose up
```

Para comprobar si Hefesto funciona de forma correcta, ir a algún
navegador de internet y navegar a `http://localhost:8080/healthcheck`

Si Hefesto ha sido levantado exitósamente, el mensaje `Hefesto funciona correctamente`
podrá leerse en pantalla.