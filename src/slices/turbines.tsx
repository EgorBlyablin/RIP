import { createAsyncThunk, createSlice } from "@reduxjs/toolkit"

import type { DsTurbine } from "../api/Api"
import { turbinesMock } from "../api/mock"
import { api } from "../api"


interface TurbinesState {
  searchFilter: string;
  turbines: DsTurbine[];
  loading: boolean;
}

const initialState: TurbinesState = {
  searchFilter: "",
  turbines: [],
  loading: false
}

export const getTurbinesList = createAsyncThunk(
  "turbines/getList",
  async (_, { getState }) => {
    const { turbines }: any = getState();

    const response = await api.turbinesList({ turbineTitle: turbines.searchFilter })

    return response.data;
  }
)

const turbinesSlice = createSlice({
  name: "turbines",
  initialState: initialState,
  reducers: {
    setTurbinesFilter: (state, action) => {
      state.searchFilter = action.payload.trim();
    }
  },
  extraReducers: (builder) => {
    builder
      .addCase(getTurbinesList.pending, (state) => {
        state.loading = true;
      })
      .addCase(getTurbinesList.fulfilled, (state, action) => {
        state.loading = false;
        state.turbines = action.payload;
      })
      .addCase(getTurbinesList.rejected, (state) => {
        state.loading = false;
        state.turbines = turbinesMock.filter((turbine) =>
          turbine.title?.toLocaleLowerCase().includes(state.searchFilter.toLocaleLowerCase())
        );
      });
  },

})

export const { setTurbinesFilter } = turbinesSlice.actions;
export const turbinesReducer = turbinesSlice.reducer;
