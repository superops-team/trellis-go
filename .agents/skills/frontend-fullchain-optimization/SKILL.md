---
name: frontend-fullchain-optimization
description: |
  Web Vitals-driven frontend performance optimization skill.
  Use when you need to audit, diagnose, or improve frontend performance
  (LCP, FID, CLS, INP, TTFB).
---

# Frontend Full-Stack Optimization

## Trigger Check

Use this skill when:
- Page load time exceeds 3 seconds
- Lighthouse score is below 80
- Core Web Vitals fail (LCP > 2.5s, FID > 100ms, CLS > 0.1)
- User reports sluggish interactions
- Bundle size exceeds 500KB (gzipped)

## Steps

### 1. Audit

1. Run Lighthouse (mobile + desktop)
2. Check Core Web Vitals via `web-vitals` library
3. Analyze bundle with `vite-bundle-visualizer` or `webpack-bundle-analyzer`
4. Check network waterfall (blocking resources, slow API calls)
5. Profile runtime performance (long tasks, layout thrashing)

### 2. Diagnose

For each issue found, identify:
- Root cause (render-blocking script? large image? slow API?)
- Impact on user experience
- Fix priority (P0 = blocks rendering, P1 = slows interaction, P2 = cosmetic)

### 3. Optimize

| Area | Technique | Impact |
|------|-----------|--------|
| Images | WebP/AVIF, lazy loading, responsive sizes | LCP |
| Fonts | `font-display: swap`, subsetting | CLS |
| JS | Code splitting, tree shaking, defer | TTI |
| CSS | Critical CSS inlining, unused CSS removal | FCP |
| API | Caching, prefetching, streaming | TTFB |
| Rendering | Virtual scrolling, `content-visibility` | CLS |

### 4. Verify

1. Re-run Lighthouse — target score ≥ 90
2. Check Web Vitals in production
3. Load test with simulated slow network (3G throttling)
4. Verify no regressions in functionality

## Output

- Performance audit report with before/after metrics
- Optimized bundle configuration
- Updated components with performance fixes
