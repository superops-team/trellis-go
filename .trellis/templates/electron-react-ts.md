# Electron + React + TypeScript Spec Template

## Overview

Desktop application using Electron, React frontend, and TypeScript.

## Project Structure

```
project/
├── src/
│   ├── main/                  # Electron main process
│   │   ├── index.ts
│   │   ├── ipc/               # IPC handlers
│   │   └── menu.ts
│   ├── preload/               # Preload scripts
│   │   └── index.ts
│   └── renderer/              # React app
│       ├── App.tsx
│       ├── components/
│       ├── pages/
│       └── styles/
├── electron-builder.yml
├── vite.config.ts
├── tsconfig.json
├── tsconfig.node.json
└── package.json
```

## Naming Conventions

- **Main process**: `kebab-case.ts`
- **Renderer components**: PascalCase
- **IPC channels**: `namespace:action` (e.g., `file:save`)
- **Events**: past tense (`file-saved`, `window-closed`)

## Testing Strategy

| Layer | Tool | Scope |
|-------|------|-------|
| Unit | Vitest | Main process, utilities |
| Component | Vitest + Testing Library | React components |
| Integration | Vitest + Electron Testing | IPC flows |
| E2E | Playwright + Electron | Full app flows |

## Deployment Checklist

- [ ] Code signing configured
- [ ] Auto-update setup (electron-updater)
- [ ] Platform-specific builds tested
- [ ] App icon set
- [ ] Crash reporting (Sentry)
- [ ] Installer tested on clean OS

## Security Considerations

- Context isolation enabled
- Node integration disabled in renderer
- IPC input validation
- CSP headers in renderer
- Auto-update signature verification
- No remote content loading
