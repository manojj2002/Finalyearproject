// src/store/actions/scanActions.ts
export const SET_TRIVY_RESULTS = 'SET_TRIVY_RESULTS';
export const SET_FALCO_ALERTS = 'SET_FALCO_ALERTS';

export const setTrivyResults = (results: any[]) => ({
  type: SET_TRIVY_RESULTS,
  payload: results,
});

export const setFalcoAlerts = (alerts: any[]) => ({
  type: SET_FALCO_ALERTS,
  payload: alerts,
});
