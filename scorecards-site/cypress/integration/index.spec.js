const URL = Cypress.env("PREVIEW_URL")
describe("Scorecard site", () => {
    it("can view homepage", () => {
      cy.visit(URL, { failOnStatusCode: false })

      cy.get('title')
        .should('not.be.empty');

      cy.get('meta[name="keywords"]')
        .should("have.attr", "content").should('not.be.empty');
      
      cy.get('meta[name="description"]')
        .should("have.attr", "content").should('not.be.empty');

        cy.get('video')
            .should('have.prop', 'paused', false)
            .and('have.prop', 'ended', false)

      cy.contains("Build better security habits, one test at a time")
      cy.get('h1').contains("Build better security habits, one test at a time").should('not.be.empty')
      cy.get('a').first().click()
      cy.contains('Run the checks').should('not.be.empty')

      cy.get('.theme-code-group__nav-tab').first().should('not.be.empty').contains('Homebrew')

      cy.get('.theme-code-group__nav-tab').first().click()

    })
})