import { createAsyncThunk, createSlice } from '@reduxjs/toolkit';
import type { DsGenerationRequest } from '../api/Api';
import { api } from '../api';

interface GenerationRequestsState {
    formedAtBegin?: string;
    formedAtEnd?: string;
    status?: "sent" | "completed" | "rejected";

    generationRequests: DsGenerationRequest[];
    error: string | null;
    loading: boolean;
}

const initialState: GenerationRequestsState = {
    formedAtBegin: new Date(Date.now() - 24 * 60 * 60 * 1000).toISOString().split("T")[0],

    generationRequests: [{}],
    error: null,
    loading: false,
};

const generationRequestsSlice = createSlice({
    name: 'generationRequests',
    initialState,
    reducers: {
        setFormedAtBegin(state, action) {
            state.formedAtBegin = action.payload;
        },
        setFormedAtEnd(state, action) {
            state.formedAtEnd = action.payload;
        },
        setStatus(state, action) {
            state.status = action.payload;
        }
    },
    extraReducers: (builder) => {
        builder
            .addCase(getGenerationRequests.pending, (state) => {
                state.loading = true;
                state.error = null;
            })
            .addCase(getGenerationRequests.fulfilled, (state, action) => {
                state.generationRequests = action.payload.sort((a, b) => (b.id || 0) - (a.id || 0));
                state.error = null;
                state.loading = false;
            })
            .addCase(getGenerationRequests.rejected, (state) => {
                state.error = "Ошибка при загрузке данных";
                state.loading = false;
            })
    }
});

export const getGenerationRequests = createAsyncThunk(
    'generationRequests/getGenerationRequests',
    async (_, { getState }) => {
        const { generationRequests }: any = getState()

        const response = await api.generationRequestsList(generationRequests.formedAtBegin ? new Date(generationRequests.formedAtBegin).toISOString() : undefined, generationRequests.formedAtEnd ? new Date(generationRequests.formedAtEnd).toISOString() : undefined, generationRequests.status);
        return response.data;
    }
)

export const completeGenerationRequest = createAsyncThunk(
    'generationRequests/completeGenerationRequest',
    async (generationRequestId: number) => {
        const response = await api.generationRequestsCloseUpdate(generationRequestId, { status: "completed" });
        return response.data;
    }
)

export const rejectGenerationRequest = createAsyncThunk(
    'generationRequests/rejectGenerationRequest',
    async (generationRequestId: number, { dispatch }) => {
        const response = await api.generationRequestsCloseUpdate(generationRequestId, { status: "rejected" });

        dispatch(getGenerationRequests())

        return response.data;
    }
)

export const { setFormedAtBegin, setFormedAtEnd, setStatus } = generationRequestsSlice.actions;
export const generationRequestsReducer = generationRequestsSlice.reducer;