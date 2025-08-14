import axios from 'axios';
import { API_BASE_URL, PROD_API_BASE_URL } from '../configs/contants';

export const axiosInstance = axios.create(
    {
        baseURL: import.meta.env.MODE === 'development' ? API_BASE_URL : PROD_API_BASE_URL,
        withCredentials: true,
        headers: {
            'Content-Type': 'application/json',
        },
    }
)

axiosInstance.interceptors.request.use( (config) => {
    const token = localStorage.getItem('token');
    if (token) {
        config.headers.Authorization = `Bearer ${token}`;
    }
    return config;
},
(error) => {
    return Promise.reject(error);
});