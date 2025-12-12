import { createSlice, createAsyncThunk } from '@reduxjs/toolkit';

import { api } from '../api';
import { getTurbinesList, setTurbinesFilter } from './turbines';
import { setGenerationRequestId, setTurbinesCount } from './generation-requests-draft';

interface UserState {
  username: string;
  isAuthenticated: boolean;
  isEngineer: boolean;
  error?: string | null;
}

const initialState: UserState = {
  username: '',
  isAuthenticated: false,
  isEngineer: false,
  error: null,
};


// Асинхронное действие для деавторизации
export const logoutThunk = createAsyncThunk(
  'user/logout',
  async (_, { rejectWithValue, dispatch }) => {
    try {
      const response = await api.usersLogoutCreate();

      dispatch(setTurbinesFilter(''));
      dispatch(getTurbinesList());

      dispatch(setGenerationRequestId(NaN));
      dispatch(setTurbinesCount(0));

      localStorage.removeItem('access_token');

      return response.data;
    } catch (error) {
      return rejectWithValue('Ошибка при выходе из системы');
    }
  }
);

const userSlice = createSlice({
  name: 'user',
  initialState,
  reducers: {
    setIsEngineer(state, action) {
      state.isEngineer = action.payload;
    }
  },
  extraReducers: (builder) => {
    builder
      .addCase(loginThunk.pending, (state) => {
        state.error = null;
      })
      .addCase(loginThunk.fulfilled, (state, action) => {
        const { username } = action.payload;
        state.username = username;
        state.isAuthenticated = true;
        state.error = null;
      })
      .addCase(loginThunk.rejected, (state, action) => {
        state.error = action.payload as string;
        state.isAuthenticated = false;
      })

      .addCase(logoutThunk.fulfilled, (state) => {
        state.username = '';
        state.isAuthenticated = false;
        state.error = null;
      })
      .addCase(logoutThunk.rejected, (state, action) => {
        state.error = action.payload as string;
      });
  },
});

// Асинхронное действие для авторизации
export const loginThunk = createAsyncThunk(
  'user/login',
  async (credentials: { username: string; password: string }, { rejectWithValue, dispatch }) => {
    try {
      const response = await api.usersLoginCreate({
        login: credentials.username,
        password: credentials.password
      });

      if (response.data.access_token) {
        localStorage.setItem('access_token', response.data.access_token);
      }

      const draftInfo = await api.generationRequestsDraftList();

      dispatch(setGenerationRequestId(draftInfo.data.generationRequestId));
      dispatch(setTurbinesCount(draftInfo.data.turbinesCount));

      const userDetails = await api.usersList();
      dispatch(setIsEngineer(userDetails.data.is_moderator))

      return { ...response.data, username: credentials.username };
    } catch (error) {
      return rejectWithValue('Ошибка авторизации'); // Возвращаем ошибку в случае неудачи
    }
  }
);

export const { setIsEngineer } = userSlice.actions;
export const userReducer = userSlice.reducer;