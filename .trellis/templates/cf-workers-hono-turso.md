# Cloudflare Workers + Hono + Turso Spec Template

## Overview

Edge-deployed API using Cloudflare Workers, Hono framework, and Turso (SQLite) database.

## Project Structure

```
project/
├── src/
│   ├── index.ts               # Entry point
│   ├── routes/                # Hono route handlers
│   ├── middleware/            # Hono middleware
│   ├── db/                    # Turso/Drizzle schema
│   │   ├── schema.ts
│   │   └── queries.ts
│   └── utils/                 # Shared utilities
├── migrations/                # Turso migrations
├── wrangler.toml
├── tsconfig.json
├── drizzle.config.ts
└── package.json
```

## Naming Conventions

- **Files**: `kebab-case.ts`
- **Routes**: RESTful, plural nouns
- **Middleware**: camelCase
- **Database**: snake_case

## Testing Strategy

| Layer | Tool | Scope |
|-------|------|-------|
| Unit | Vitest | Pure functions |
| Integration | Vitest + Miniflare | Worker endpoints |
| E2E | Playwright | Full edge flow |

## Deployment Checklist

- [ ] `wrangler.toml` configured
- [ ] Environment secrets set
- [ ] Turso database created
- [ ] Migrations applied
- [ ] Custom domain configured
- [ ] D1/Turso replication verified

## Security Considerations

- Validate all inputs with Zod
- Use Workers KV for rate limiting
- Implement CORS properly
- Secure service bindings
- Environment variable validation at startup
