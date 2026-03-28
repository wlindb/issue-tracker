import { defineConfig } from 'orval'

export default defineConfig({
  issueTracker: {
    input: {
      target: '../backend/api/openapi.yaml',
    },
    output: {
      target: './src/api/generated',
      override: {
        mutator: {
          path: './src/api/mutator.ts',
          name: 'customFetch',
        },
      },
    },
  },
})
