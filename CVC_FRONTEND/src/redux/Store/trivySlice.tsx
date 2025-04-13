// src/store/trivySlice.ts
import { createSlice, PayloadAction } from "@reduxjs/toolkit";

interface TrivyScanState {
  results: any[];
}

const initialState: TrivyScanState = {
  results: [],
};

const trivySlice = createSlice({
  name: "trivy",
  initialState,
  reducers: {
    setTrivyResults(state, action: PayloadAction<any[]>) {
      state.results = action.payload;
    },
  },
});

export const { setTrivyResults } = trivySlice.actions;

export default trivySlice.reducer;
