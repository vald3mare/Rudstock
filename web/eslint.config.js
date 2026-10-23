// Стандарты проекта по Bulletproof React:
// https://github.com/alan2207/bulletproof-react/blob/master/docs/project-standards.md
import js from '@eslint/js'
import prettier from 'eslint-config-prettier'
import checkFile from 'eslint-plugin-check-file'
import { importX } from 'eslint-plugin-import-x'
import { createTypeScriptImportResolver } from 'eslint-import-resolver-typescript'
import jsxA11y from 'eslint-plugin-jsx-a11y'
import react from 'eslint-plugin-react'
import reactHooks from 'eslint-plugin-react-hooks'
import globals from 'globals'
import tseslint from 'typescript-eslint'

export default tseslint.config(
  { ignores: ['dist', 'node_modules'] },

  {
    files: ['**/*.{ts,tsx}'],
    extends: [
      js.configs.recommended,
      // строгие правила с учётом типов: ловят забытый await, лишние проверки, небезопасный any
      tseslint.configs.strictTypeChecked,
      tseslint.configs.stylisticTypeChecked,
      react.configs.flat.recommended,
      react.configs.flat['jsx-runtime'],
      reactHooks.configs.flat.recommended,
      // доступность: alt у картинок, роли, подписи
      jsxA11y.flatConfigs.recommended,
      importX.flatConfigs.recommended,
      importX.flatConfigs.typescript,
    ],
    languageOptions: {
      ecmaVersion: 2023,
      globals: globals.browser,
      parserOptions: {
        projectService: true,
        tsconfigRootDir: import.meta.dirname,
      },
    },
    settings: {
      react: { version: 'detect' },
      // понимает алиас @/ из tsconfig
      'import-x/resolver-next': [createTypeScriptImportResolver()],
    },
    plugins: {
      'check-file': checkFile,
    },
    rules: {
      // типы props проверяет TypeScript, PropTypes не нужны
      'react/prop-types': 'off',

      // в проекте описываем типы через type, как в Bulletproof, а не interface
      '@typescript-eslint/consistent-type-definitions': ['error', 'type'],

      // имена файлов и папок в kebab-case: product-card.tsx, catalog-card.ts
      'check-file/filename-naming-convention': [
        'error',
        { 'src/**/*.{ts,tsx}': 'KEBAB_CASE' },
        { ignoreMiddleExtensions: true }, // product-card.test.tsx
      ],
      'check-file/folder-naming-convention': ['error', { 'src/**/!(__tests__)': 'KEBAB_CASE' }],

      // однонаправленная архитектура: shared -> features -> app
      'import-x/no-restricted-paths': [
        'error',
        {
          zones: [
            // фичи не импортируют друг друга; общее выносится в shared или собирается в app
            {
              target: './src/features/catalog',
              from: './src/features',
              except: ['./catalog'],
            },
            {
              target: './src/features/user',
              from: './src/features',
              except: ['./user'],
            },
            // фичи не знают про app
            {
              target: './src/features',
              from: './src/app',
            },
            // общий код не знает ни про фичи, ни про app
            {
              target: [
                './src/components',
                './src/config',
                './src/hooks',
                './src/lib',
                './src/types',
                './src/utils',
              ],
              from: ['./src/features', './src/app'],
            },
          ],
        },
      ],
    },
  },

  // конфиги в корне (eslint.config.js и т.п.) без проверок, которым нужны типы
  {
    files: ['**/*.js'],
    extends: [js.configs.recommended, tseslint.configs.disableTypeChecked],
    languageOptions: { globals: globals.node },
  },

  // отключает правила оформления, которые конфликтуют с Prettier: оформлением занимается он
  prettier,
)
