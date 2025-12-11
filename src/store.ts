import { configureStore } from "@reduxjs/toolkit"
import { type TypedUseSelectorHook, useDispatch as useUntypedDispatch, useSelector as useUntypedSelector } from 'react-redux';

import { turbinesReducer } from "./slices/turbines"
import { userReducer } from "./slices/user";
import { generationRequestDraftReducer } from "./slices/generation-requests-draft";
import { generationRequestsReducer } from "./slices/generation-requests";


const store = configureStore({
    reducer: {
        turbines: turbinesReducer,
        user: userReducer,
        generationRequestDraft: generationRequestDraftReducer,
        generationRequests: generationRequestsReducer
    },
});

export default store

export type RootState = ReturnType<typeof store.getState>;
export type AppDispatch = typeof store.dispatch;

export const useDispatch: () => AppDispatch = useUntypedDispatch;
export const useSelector: TypedUseSelectorHook<RootState> = useUntypedSelector;