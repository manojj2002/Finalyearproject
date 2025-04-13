// src/store/reducers/scanReducer.ts
import { SET_TRIVY_RESULTS, SET_FALCO_ALERTS } from "../Actions/scanActions"

const initialState = {
  trivyResults: [],
  falcoAlerts: [],
};

const scanReducer = (state = initialState, action: any) => {
  switch (action.type) {
    case SET_TRIVY_RESULTS:
      return { ...state, trivyResults: action.payload };
    case SET_FALCO_ALERTS:
      return { ...state, falcoAlerts: action.payload };
    default:
      return state;
  }
};

export default scanReducer;
