// src/slices/falcoSlice.ts
import { createSlice } from "@reduxjs/toolkit";

interface FalcoState {
  isRunning: boolean;
}

const initialState: FalcoState = {
  isRunning: localStorage.getItem("falcoRunning") === "true",
};

const falcoSlice = createSlice({
  name: "falco",
  initialState,
  reducers: {
    startFalco(state) {
      state.isRunning = true;
      localStorage.setItem("falcoRunning", "true");
    },
    stopFalco(state) {
      state.isRunning = false;
      localStorage.setItem("falcoRunning", "false");
    },
    toggleFalco(state) {
      state.isRunning = !state.isRunning;
      localStorage.setItem("falcoRunning", state.isRunning.toString());
    },
  },
});

export const { startFalco, stopFalco, toggleFalco } = falcoSlice.actions;
export default falcoSlice.reducer;
