import { useEffect } from 'react'
import { BrowserRouter, Route, Routes } from 'react-router'

import 'bootstrap/dist/css/bootstrap.min.css'
import { Layout } from './components/layout'
import { GenerationRequestPage } from './pages/generation-request'
import LoginPage from './pages/login'
import { MainPage } from './pages/main'
import { TurbineDetailsPage } from './pages/turbine-details'
import { TurbinesListPage } from './pages/turbines-list'
import { ROUTES } from './routes'
import { checkAuthThunk } from './slices/user'
import { useDispatch } from './store'
import { ProfilePage } from './pages/profile'
import { RegistrationPage } from './pages/registration'
import { GenerationRequestsPage } from './pages/generation-requests'


export const App = () => {
    const dispatch = useDispatch();

    useEffect(() => {
        dispatch(checkAuthThunk());
    }, []);

    return (
        <BrowserRouter>
            <Routes>
                <Route path={ROUTES.HOME} element={<MainPage />} />
                <Route path={ROUTES.LOGIN} element={<LoginPage />} />
                <Route path={ROUTES.REGISTRATION} element={<RegistrationPage />} />
                <Route path={ROUTES.TURBINES_LIST} element={<Layout />}>
                    <Route index element={<TurbinesListPage />} />
                    <Route path=":turbineId" element={<TurbineDetailsPage />} />
                </Route>
                <Route path={ROUTES.GENERATION_REQUESTS} element={<Layout />}>
                    <Route index element={<GenerationRequestsPage />} />
                    <Route path=":generationRequestId" element={<GenerationRequestPage />} />
                </Route>
                <Route path={ROUTES.PROFILE} element={<Layout />}>
                    <Route index element={<ProfilePage />} />
                </Route>
            </Routes>
        </BrowserRouter>
    )
}
