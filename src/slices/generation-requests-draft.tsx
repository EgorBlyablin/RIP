import { createAsyncThunk, createSlice } from '@reduxjs/toolkit';
import type { DsDraftGenerationRequestsBriefInfo, DsGenerationRequest, DsTurbineGenerationRequest } from '../api/Api';
import { api } from '../api';

interface generationRequestData extends DsGenerationRequest { }

interface GenerationRequestDraftState extends DsDraftGenerationRequestsBriefInfo {
    generationRequestId: number;
    turbinesCount: number;

    isDraft: boolean;
    turbines: DsTurbineGenerationRequest[];
    generationRequestData: generationRequestData;
    error: string | null;
    loading: boolean;
}

const initialState: GenerationRequestDraftState = {
    generationRequestId: NaN,
    turbinesCount: NaN,

    isDraft: false,
    turbines: [],
    generationRequestData: {},
    error: null,
    loading: false,
};

const generationRequestDraftSlice = createSlice({
    name: 'generationRequestDraft',
    initialState,
    reducers: {
        setGenerationRequestId: (state, action) => {
            state.generationRequestId = action.payload;
        },
        setTurbinesCount: (state, action) => {
            state.turbinesCount = action.payload;
        },
        setError: (state, action) => {
            state.error = action.payload;
        },
    },
    extraReducers: (builder) => {
        builder
            .addCase(getGenerationRequest.pending, (state) => {
                state.loading = true;
                state.error = null;
            })
            .addCase(getGenerationRequest.fulfilled, (state, action) => {
                const { turbine_generation_requests, id, period_days } = action.payload;

                if (turbine_generation_requests && id) {
                    state.generationRequestId = id;
                    state.turbinesCount = turbine_generation_requests.length;

                    state.generationRequestData = {
                        period_days: period_days,
                    };
                    state.turbines = turbine_generation_requests;
                }

                state.isDraft = action.payload.status === "draft";
                state.generationRequestData = action.payload;
                state.error = null;
                state.loading = false;
            })
            .addCase(getGenerationRequest.rejected, (state) => {
                state.error = "Ошибка при загрузке данных";
                state.loading = false;
            })

            .addCase(deleteGenerationRequest.fulfilled, (state) => {
                state.generationRequestId = NaN;
                state.turbinesCount = 0;
                state.turbines = [];
                state.generationRequestData = {
                    period_days: 0
                }
                state.error = null;
            })
            .addCase(deleteGenerationRequest.rejected, (state) => {
                state.error = 'Ошибка при удалении заявки расчета';
            })

            .addCase(setGenerationRequestData.fulfilled, (state, action) => {
                state.generationRequestData.period_days = action.payload.period_days;
                state.error = null;
            })
            .addCase(setGenerationRequestData.rejected, (state) => {
                state.error = 'Ошибка при обновлении заявки расчета';
            })

            .addCase(sendGenerationRequest.fulfilled, (state) => {
                state.generationRequestId = NaN
                state.turbinesCount = NaN

                state.turbines = []
                state.generationRequestData = {period_days: undefined}
                state.loading = false

                state.error = null;
            })
            .addCase(sendGenerationRequest.rejected, (state) => {
                state.error = 'Ошибка при отправке заявки расчета';
            })
    }
});

export const getGenerationRequest = createAsyncThunk(
    'generationRequestDraft/getGenerationRequest',
    async (generationRequestId: number) => {
        const response = await api.generationRequestsList2(generationRequestId);
        return response.data;
    }
)

export const setGenerationRequestData = createAsyncThunk(
    'generationRequestDraft/setGenerationRequestData',
    async (periodDays: number) => {
        const response = await api.generationRequestsDraftUpdate({
            period_days: periodDays
        });
        return response.data;
    }
);

export const sendGenerationRequest = createAsyncThunk(
    'generationRequestDraft/sendGenerationRequest',
    async () => {
        const response = await api.generationRequestsDraftSubmitUpdate();
        return response.data;
    }
);

export const deleteGenerationRequest = createAsyncThunk(
    'generationRequestDraft/deleteGenerationRequest',
    async () => {
        const response = await api.generationRequestsDraftDelete();
        return response.data;
    }
);

export const deleteTurbineFromGenerationRequest = createAsyncThunk(
    'generationRequestDraft/deleteTurbineFromGenerationRequest',
    async (turbineId: number, { dispatch, getState }) => {
        const { generationRequestDraft }: any = getState();

        const response = await api.generationRequestsDraftDelete2(turbineId);

        dispatch(getGenerationRequest(generationRequestDraft.generationRequestId))

        return response.data;
    }
);

export const updateTurbineInGenerationRequest = createAsyncThunk(
    'generationRequestDraft/updateTurbineInGenerationRequest',
    async ({
        turbineId,
        avgWind,
        alpha,
    }: {
        turbineId: number;
        avgWind: number | undefined;
        alpha: number | undefined;
    }, { dispatch, getState }) => {
        const { generationRequestDraft }: any = getState();

        const response = await api.generationRequestsDraftUpdate2(turbineId, {
            avg_velocity: avgWind,
            alpha: alpha
        });

        dispatch(getGenerationRequest(generationRequestDraft.generationRequestId))

        return response.data;
    }
);

export const addTurbineToGenerationRequest = createAsyncThunk(
    'generationRequestDraft/addTurbineToGenerationRequest',
    async (turbineId: number, { dispatch }) => {
        const response = await api.generationRequestsDraftCreate(turbineId);

        const draftInfo = await api.generationRequestsDraftList();

        dispatch(setGenerationRequestId(draftInfo.data.generationRequestId));
        dispatch(setTurbinesCount(draftInfo.data.turbinesCount));

        return response.data;
    }
)

export const { setGenerationRequestId, setTurbinesCount, setError } = generationRequestDraftSlice.actions;
export const generationRequestDraftReducer = generationRequestDraftSlice.reducer;