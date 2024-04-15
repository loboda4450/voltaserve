module.exports = {
  root: true,
  env: { browser: true, es2020: true },
  files: ['./src/**/*.{ts,tsx}'],
  extends: [
    'eslint:recommended',
    'plugin:@typescript-eslint/recommended',
    'plugin:react-hooks/recommended',
    'prettier',
  ],
  ignorePatterns: ['dist', '.eslintrc.cjs'],
  parser: '@typescript-eslint/parser',
  plugins: ['react-refresh'],
  rules: {
    'react-refresh/only-export-components': [
      'warn',
      { allowConstantExport: true },
    ],
    '@typescript-eslint/no-explicit-any': 'off',
    'react-hooks/rules-of-hooks': 'off',
    'react-refresh/only-export-components': 'off',
  },
}