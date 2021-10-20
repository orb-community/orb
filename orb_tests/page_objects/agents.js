var loginActions = {
    AgentsPage: function () {
        return this
        .waitForElementVisible('@path', "Agents path is visible")
        .assert.containsText('@agentPath', 'Fleet Management', "Agents is inherited from Fleet Management")
        .assert.containsText('@view', 'Agent List', "Agent view is named 'Agent List'")
        .assert.containsText('@header', "All Agents", "Agents Header is 'All Agents'")
        .waitForElementVisible('@table', "Agent table is visible")
        .waitForElementVisible("@new", "New Agent button is visible")
        .waitForElementVisible("@filter", "Filter type is visible")
        .waitForElementVisible("@search", "Search by filter is visible")

    }
  }
  
  module.exports = {
    url: '/pages/fleet/agents',
    commands: [loginActions],
    elements: {
      path: 'xng-breadcrumb.orb-breadcrumb',
      agentPath: '.xng-breadcrumb-link',
      view: '.xng-breadcrumb-trail',
      header: 'ngx-agent-list-component.ng-star-inserted > div:nth-child(1) > header:nth-child(1) > h4:nth-child(2)',
      table:'.datatable-body',
      new:'.status-primary',
      filter:'.select-button',
      search:'input.size-medium',
      agentsListed: '.datatable-row-wrapper',
      info: '.sink-info-accent',
      emptyRow: '.empty-row',
      countMessage: '.justify-content-between > div:nth-child(1)',
      count:'.page-count'


    }
  }
