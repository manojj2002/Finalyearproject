// src/store.ts
import { configureStore } from "@reduxjs/toolkit";
import userReducer from "./Store/userSlice";  
import falcoReducer from "./Store/falcoSlice";
import trivyReducer from "./Store/trivySlice"
export const store = configureStore({
  reducer: {
    user: userReducer,
    falco: falcoReducer,
    trivy: trivyReducer
  },
});

export type RootState = ReturnType<typeof store.getState>;
export type AppDispatch = typeof store.dispatch;
