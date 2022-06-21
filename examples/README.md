# README

Create file called `.env` use the example `.env.example` and update the values to suit your environment.

Once you have the `.env` file you can run any of the examples like

e.g.

```shell
drone exec --secret-file=.env
```

(OR)

```shell
drone exec --secret-file=.env .drone-unauthenticated.yml
```
