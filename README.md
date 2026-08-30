OneBill
=======

**One number for what you spend on cloud and AI.**

You run AWS, Azure, Hetzner, plus OpenAI and Anthropic. Each one shows its own
bill in its own dashboard. Nobody can answer the simple question: *what did we
spend last month, in total, across everything?*

OneBill answers it. Add read-only API credentials per provider, a daily job
pulls the cost data, and one page shows the sum.

Self-hosted, so your credentials never leave your infrastructure.

Status
------

**Early. Not usable yet.**

This is the application skeleton — users, sessions, authentication, settings,
a REST API, migrations, and Docker packaging — carved out of
[Miniflux 2](https://github.com/miniflux/v2). The cost providers are not
written yet. See [CLAUDE.md](CLAUDE.md) for how the code is laid out and what
is left to build.

What it will do
---------------

- Total spend this month, last month, and a breakdown by provider.
- Providers: Hetzner, AWS, and Anthropic to start. One file each, so adding
  another is a small contribution.
- Daily sync. Credentials encrypted at rest.
- Original and converted amounts, with the exchange rate and its date
  recorded, so the total can be checked.

What it will not do: forecasting, anomaly detection, rightsizing, unit
economics, chargeback. Other tools do those. This one is a number.

Why
---

Every commercial cost tool is SaaS-only, priced for companies with a FinOps
team, and none of them know the European hosters exist. OneBill is
self-hosted, free, open source, and speaks Hetzner.

Requirements
------------

- Go >= 1.26
- PostgreSQL

Development
-----------

```bash
docker run --rm --name onebill-db -p 5432:5432 \
  -e POSTGRES_DB=onebill -e POSTGRES_USER=postgres -e POSTGRES_PASSWORD=postgres postgres

make run     # http://127.0.0.1:8080, admin / test123
make test
make lint
```

`make onebill` builds the binary. See the [man page](onebill.1) for command
line flags and configuration options.

Credits
-------

OneBill is a fork of [Miniflux 2](https://github.com/miniflux/v2) by
Frédéric Guillot and contributors, and reuses its HTTP server, storage layer,
migration runner, template engine, configuration parser, and CLI. See
[NOTICE](NOTICE).

License
-------

Apache 2.0. See [LICENSE](LICENSE).
