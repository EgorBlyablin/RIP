import { Api } from './Api';

export const api = new Api({
    baseURL: 'https://localhost:3000',
    securityWorker: () => {
        const headers = {
            headers: {
                Authorization: localStorage.getItem('access_token') ? `Bearer ${localStorage.getItem('access_token')}` : ''
            }
        }
        return headers;
    },
}).api;