// src/slices/userSlice.ts
import { createSlice, PayloadAction } from "@reduxjs/toolkit";

interface UserState {
  token: string | null;
  username: string | null
}

const initialState: UserState = {
  token: localStorage.getItem("token"),
  username: localStorage.getItem("username")
};

const userSlice = createSlice({
  name: "user",
  initialState,
  reducers: {
    setUser(state, action: PayloadAction<UserState>) {
      const { token, username, } = action.payload;
      state.token = token;
      state.username = username;
      localStorage.setItem("token", token || "");
      localStorage.setItem("username", username || "");
    },
    logout(state) {
      state.token = null;
      state.username = null;
      localStorage.removeItem("token");
      localStorage.removeItem("username");
    },
  },
});

export const { setUser, logout } = userSlice.actions;
export default userSlice.reducer;
