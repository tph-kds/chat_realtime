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