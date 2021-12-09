# `/web`

Web application specific components: static web assets, server side templates and SPAs.

https://github.com/golang-standards/project-layout/tree/master/web

### `/web/i18n`

Please name your translation files according to the locale (e.g. `de.toml`, `en.toml` or `en-uk.toml` and `en-us.toml`). We assume that any translation file hold all keys (no key mixing between locales)!

### `/web/templates`

This directory should e.g. hold email related templates (used by `/internal/mailer`).