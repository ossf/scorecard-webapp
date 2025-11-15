// @ts-check
import path from 'path'
import fs from 'fs'
import { test, expect } from '@playwright/test'

const mockRepoDataPath = path.resolve(
  __dirname,
  'mock-data',
  'github-repo.json'
)
const mockRepoData = fs.readFileSync(mockRepoDataPath, 'utf-8')

const mockCommitsDataPath = path.resolve(
  __dirname,
  'mock-data',
  'github-commits.json'
)
const mockCommitsData = fs.readFileSync(mockCommitsDataPath, 'utf-8')

test('Home Page', async ({ page }) => {
  // Mock GitHub repo metadata used by Header.getTotalCommits()
  await page.route('https://api.github.com/repos/ossf/scorecard', (route) => {
    route.fulfill({
      status: 200,
      contentType: 'application/json',
      body: mockRepoData,
    })
  })

  // Mock commits endpoint used by Header.fetchData()
  await page.route(
    'https://api.github.com/repos/ossf/scorecard/commits*',
    (route) => {
      route.fulfill({
        status: 200,
        contentType: 'application/json',
        body: mockCommitsData,
      })
    }
  )

  await page.goto('http://localhost:3000/')

  // Expect a title "to contain" a substring.
  await expect(page).toHaveTitle(/OpenSSF Scorecard/)

  // Check the video player is present and in autoplay mode
  const videoElement = page.locator('video:visible')
  await expect(videoElement).toBeVisible()
  const autoplay = await videoElement.getAttribute('autoplay')
  expect(autoplay).not.toBeNull()

  // Hide all video players for consistent screenshots
  await page.evaluate(() => {
    document
      .querySelectorAll('video')
      .forEach((v) => (v.style.display = 'none'))
  })

  // Prevent Visual Regressions
  expect(await page.screenshot({ fullPage: true })).toMatchSnapshot('home.png')
})
