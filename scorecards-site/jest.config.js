module.exports = {
    moduleNameMapper: {
        "^@/(.*svg)(\\?inline)$": "<rootDir>/$1",
      '^@/(.*)$': '<rootDir>/$1',
      '^~/(.*)$': '<rootDir>/$1',
      '^vue$': 'vue/dist/vue.common.js',
    },
    modulePathIgnorePatterns: ["<rootDir>/dist/", "<rootDir>/cypress/*"],
    moduleFileExtensions: ['js', 'vue', 'json'],
    transform: {
      '^.+\\.js$': 'babel-jest',
      '.*\\.(vue)$': 'vue-jest',
      '.+\\.(css|styl|less|sass|scss|svg|png|jpg|ttf|woff|woff2)(\\?inline)?$': 'jest-transform-stub',
    },
    collectCoverage: true,
    collectCoverageFrom: [
      '<rootDir>/components/**/*.vue',
      '<rootDir>/modules/**/*.vue',
    ],
    testEnvironment: 'jsdom',
}
  