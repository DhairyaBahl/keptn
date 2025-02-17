import { interceptEnvironmentScreen } from '../support/intercept';
import EnvironmentPage from '../support/pageobjects/EnvironmentPage';

describe('Environment Screen', () => {
  const environmentPage = new EnvironmentPage();
  const project = 'sockshop';
  const stage = 'dev';

  beforeEach(() => {
    interceptEnvironmentScreen();
    environmentPage.visit('sockshop');
  });

  it('should show evaluation history loading indicators', () => {
    const service = 'carts';
    cy.intercept(environmentPage.getEvaluationHistoryURL(project, stage, service, 6), {
      delay: 10_000,
    });
    environmentPage.selectStage(stage).assertEvaluationHistoryLoadingCount(service, 5);
  });

  it('should not show evaluation history loading indicators', () => {
    environmentPage.selectStage(stage).assertEvaluationHistoryLoadingCount('carts', 0);
  });

  it('should not show evaluation history', () => {
    environmentPage.selectStage(stage).assertEvaluationHistoryCount('carts', 0);
  });

  it('should not show evaluation', () => {
    environmentPage.selectStage(stage).assertEvaluationInDetails('carts-db', '-');
  });

  it('should show evaluations in history if sequence does not have an evaluation task', () => {
    const service = 'carts-db';
    cy.intercept(environmentPage.getEvaluationHistoryURL(project, stage, service, 5), {
      fixture: 'get.environment.evaluation.history.carts-db.mock',
    });
    environmentPage
      .selectStage(stage)
      .assertEvaluationHistoryLoadingCount(service, 0)
      .assertEvaluationHistoryCount(service, 5)
      .assertEvaluationInDetails(service, '-');
  });

  it('should show 2 evaluations in history and should not show current evaluation in history', () => {
    const service = 'carts';
    cy.intercept(environmentPage.getEvaluationHistoryURL(project, stage, service, 6), {
      fixture: 'get.environment.evaluation.history.limited.mock', // 3 events, including the current one
    });
    environmentPage
      .selectStage(stage)
      .assertEvaluationHistoryCount(service, 2)
      .assertEvaluationInDetails(service, 0, 'success');
  });

  it('should show 5 evaluation in history', () => {
    const service = 'carts';
    cy.intercept(environmentPage.getEvaluationHistoryURL(project, stage, service, 6), {
      fixture: 'get.environment.evaluation.history.mock',
    });
    environmentPage
      .selectStage(stage)
      .assertEvaluationHistoryCount(service, 5)
      .assertEvaluationInDetails(service, 0, 'success');
  });
});
