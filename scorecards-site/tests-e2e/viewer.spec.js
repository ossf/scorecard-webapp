// @ts-check
import path from 'path'
import fs from 'fs'
import { test, expect } from '@chromatic-com/playwright'

const mockScorecardDataPath = path.resolve(
  __dirname,
  'mock-data',
  'scorecard-results.json'
)
const mockScorecardData = fs.readFileSync(mockScorecardDataPath, 'utf-8')

test('viewer welcome', async ({ page }) => {
  await page.goto('http://localhost:3000/viewer')

  // Prevent Visual Regressions
  if (process.env.ENABLE_SNAPSHOTS) {
    expect(await page.screenshot({ fullPage: true })).toMatchSnapshot(
      'viewer-welcome.png'
    )
  }

  // Find input field and type the repo URL.
  const input = page.locator('input[placeholder="github.com/ossf/scorecard"]')
  await input.fill('github.com/ossf/scorecard')

  // Click the Submit button.
  await page.getByRole('button', { name: 'Submit' }).click()

  // Wait for the results page to load.
  await expect(page).toHaveURL(
    'http://localhost:3000/viewer/?uri=github.com%2Fossf%2Fscorecard'
  )
})

test('viewer details', async ({ page }) => {
  // Mock Scorecard results endpoint
  await page.route(
    'https://api.scorecard.dev/projects/github.com/ossf/scorecard',
    (route) => {
      route.fulfill({
        status: 200,
        contentType: 'application/json',
        body: mockScorecardData,
      })
    }
  )

  await page.goto('http://localhost:3000/viewer/?uri=github.com/ossf/scorecard')

  // Wait for the results to load.
  await expect(
    page.getByText(
      "Determines if the project's GitHub Action workflows avoid dangerous patterns."
    )
  ).toBeVisible()

  // Wait for 5 seconds
  await page.waitForTimeout(5 * 1000)

  // Prevent Visual Regressions
  if (process.env.ENABLE_SNAPSHOTS) {
    expect(await page.screenshot({ fullPage: true })).toMatchSnapshot(
      'viewer-details.png'
    )
  }
})
