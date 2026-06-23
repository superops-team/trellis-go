# Next.js + oRPC + PostgreSQL Spec Template

## Overview

Full-stack web application using Next.js App Router, oRPC for type-safe API, and PostgreSQL with Drizzle ORM.

## Project Structure

```
project/
├── src/
│   ├── app/                    # Next.js App Router pages
│   │   ├── layout.tsx
│   │   ├── page.tsx
│   │   ├── api/                # Route handlers (thin layer)
│   │   └── (auth)/             # Auth-required routes
│   ├── components/             # React components
│   │   ├── ui/                 # Primitive UI components
│   │   └── features/           # Feature-specific components
│   ├── lib/
│   │   ├── orpc/               # oRPC client & server setup
│   │   ├── db/                 # Drizzle schema & queries
│   │   │   ├── schema/         # Table definitions
│   │   │   └── queries/        # Reusable query functions
│   │   └── auth/               # Authentication utilities
│   └── types/                  # Shared TypeScript types
├── migrations/                 # Drizzle migrations
├── drizzle.config.ts
├── next.config.ts
├── tsconfig.json
└── package.json
```

## Naming Conventions

- **Files**: `kebab-case.ts` (components: `PascalCase.tsx`)
- **Components**: PascalCase
- **Functions**: camelCase
- **Database tables**: snake_case
- **Database columns**: snake_case
- **API routes**: kebab-case
- **Environment variables**: `UPPER_SNAKE_CASE`

## Testing Strategy

| Layer | Tool | Scope |
|-------|------|-------|
| Unit | Vitest | Pure functions, utilities |
| Component | Vitest + Testing Library | React components |
| API | Vitest + supertest | Route handlers |
| E2E | Playwright | Full user flows |

## Deployment Checklist

- [ ] Environment variables configured
- [ ] Database migrations run
- [ ] SSL certificate configured
- [ ] CDN cache headers set
- [ ] Error monitoring (Sentry)
- [ ] Health check endpoint
- [ ] Rate limiting configured
- [ ] CORS policy set

## Security Considerations

- Input validation on all API endpoints
- SQL injection prevention (Drizzle parameterized queries)
- CSRF protection (Next.js built-in)
- Rate limiting on auth endpoints
- Session management with HTTP-only cookies
- Content Security Policy headers
