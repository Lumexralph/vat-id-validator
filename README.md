VAT ID Validator Microservice
==
A VAT Number validator microservice, to validate on German VAT ID.


Application Set-Up
---

You need to clone the repository to your computer by doing the following:

```bash
git clone https://github.com/Lumexralph/vat-id-validator.git

cd vat-id-validator
```

You can run the application in 3 ways:
1. Using the `Go tool`.
2. Using the application binary built with `Makefile`.
3. Using a Docker container image built using the `Makefile`.

Using the `Go tool`
---
1. Ensure you are in the root directory of the repository.
2. Run the following command:
```bash
go run cmd/main.go
```

Using the application binary built with `Makefile`
---
1. Ensure you are in the root directory of the repository.
2. Run `make build`, it will create the binary `vat-id-validator` in the root directory of the repository.
3. Run the following command:

```bash
vat-id-validator
```

Using a Docker container image built using the `Makefile`.
---
1. Ensure you are in the root directory of the repository.
2. Run `make docker-build`, this will create the container image `vat-id-validator:latest`
3. Run the following:

```bash
docker run vat-id-validator:latest 
```

You can run tests, fmt, lint and build using just `make` but this will require you have
`golangci-lint` installed.

PS: Download instructions will be outputted during the `make` process.

